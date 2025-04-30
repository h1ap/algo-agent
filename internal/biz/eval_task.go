package biz

import (
	"algo-agent/internal/cons/mq"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	cc "algo-agent/internal/cons/container"
	"algo-agent/internal/cons/file"
	taskArgs "algo-agent/internal/cons/task"
	ct "algo-agent/internal/cons/train"
	mc "algo-agent/internal/model/container"
	"algo-agent/internal/model/eval"
	"algo-agent/internal/mq/event"
	"algo-agent/internal/utils"

	"github.com/go-kratos/kratos/v2/log"
)

// EvalTaskUsecase 评估任务用例
type EvalTaskUsecase struct {
	etm EvalTaskManager

	mq  MqService
	d   DockerService
	oss OSSService

	log *log.Helper

	filePath string
	tsn      string // 训练服务名称

	// 常量
	evalScriptName string
}

// 更新状态并发送消息
func (etu *EvalTaskUsecase) editTaskAndSendMq(ctx context.Context, task *eval.EvalTaskInfo) error {
	etu.etm.UpdateTask(ctx, task)
	return etu.sendStatusChangeMessage(ctx, task)
}

// 发送状态变更消息
func (etu *EvalTaskUsecase) sendStatusChangeMessage(ctx context.Context, taskInfo *eval.EvalTaskInfo) error {
	reply := &event.EvalReceiveMessage{
		TaskId: taskInfo.TaskId,
		Status: taskInfo.TaskStatus,
		Remark: taskInfo.Remark,
	}

	return etu.mq.SendToService(ctx, etu.tsn, &event.ReqMessage{
		Type:    mq.TASK_EVALUATE.Code(),
		Payload: reply,
	})
}

// StartEvaluation 开始评估任务
func (etu *EvalTaskUsecase) StartEvaluation(ctx context.Context, taskInfo *eval.EvalTaskInfo) error {
	eventInfo := taskInfo.EvalSendMessage
	if eventInfo == nil || eventInfo.TaskId == "" {
		etu.log.Error("开始评估失败：任务ID为空")
		return errors.New("任务ID不允许为空")
	}
	etu.log.Infof("开始评估任务，任务ID: %s", eventInfo.TaskId)

	// 设置任务状态为"评估中"
	taskInfo.TaskStatus = ct.TrainJobStatusStarting.Code
	if err := etu.etm.AddTask(ctx, taskInfo); err != nil {
		return fmt.Errorf("添加任务失败: %v", err)
	}
	// 发送任务开始的状态到系统消息队列
	if err := etu.sendStatusChangeMessage(ctx, taskInfo); err != nil {
		return fmt.Errorf("发送状态变更消息失败: %v", err)
	}

	// 创建必要的目录
	evalBasePath := filepath.Join(etu.filePath, file.EVAL)
	imagePath := filepath.Join(etu.filePath, file.IMAGE)
	utils.EnsureDirectoryExists(evalBasePath)
	utils.EnsureDirectoryExists(imagePath)

	// 检查Docker镜像是否存在
	etu.log.Infof("检查评估镜像是否存在: %s", eventInfo.AlgorithmTrainImageName)
	imageExists, err := etu.d.FindImageByName(ctx, eventInfo.AlgorithmTrainImageName)
	if err != nil || !imageExists {
		etu.log.Info("本地未找到评估镜像，准备下载")
		taskInfo.TaskStatus = ct.TrainJobStatusDownloadTrainImage.Code
		if err := etu.editTaskAndSendMq(ctx, taskInfo); err != nil {
			return fmt.Errorf("更新任务状态失败: %v", err)
		}

		// 下载镜像文件
		imageFileName := strings.Replace(eventInfo.AlgorithmTrainImageName, ":", "-", -1)
		tarPath := filepath.Join(imagePath, imageFileName+".tar")
		err = etu.oss.DownloadSingleFile(
			ctx,
			eventInfo.AlgorithmTrainImageMinioBucket,
			eventInfo.AlgorithmTrainImageMinioURL,
			imagePath,
			imageFileName+".tar",
		)
		if err != nil {
			return fmt.Errorf("下载镜像文件失败: %v", err)
		}

		// 导入镜像
		err = etu.d.ImportAndTagImage(ctx, tarPath, eventInfo.AlgorithmTrainImageName)
		if err != nil {
			return fmt.Errorf("导入镜像失败: %v", err)
		}
	} else {
		etu.log.Info("找到已存在的评估镜像")
	}

	// 设置任务路径
	taskPath := filepath.Join(evalBasePath, eventInfo.TaskId)

	// 下载算法脚本
	etu.log.Info("开始下载算法脚本")
	taskInfo.TaskStatus = ct.TrainJobStatusDownloadAlgoScripts.Code
	if err := etu.editTaskAndSendMq(ctx, taskInfo); err != nil {
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	algorithmScriptZipFile := "script_" + eventInfo.TaskId + ".zip"
	// 下载算法脚本
	err = etu.oss.DownloadSingleFile(
		ctx,
		eventInfo.AlgorithmScriptMinioBucket,
		eventInfo.AlgorithmScriptZipMinioURL,
		taskPath,
		algorithmScriptZipFile,
	)
	if err != nil {
		return fmt.Errorf("下载算法脚本失败: %v", err)
	}

	// 解压脚本文件
	etu.log.Info("解压算法脚本")
	scriptZipPath := filepath.Join(taskPath, algorithmScriptZipFile)
	scriptDestPath := filepath.Join(taskPath, file.SCRIPT)
	if err := utils.Unzip(scriptZipPath, scriptDestPath); err != nil {
		return fmt.Errorf("解压算法脚本失败: %v", err)
	}

	// 下载评估数据
	etu.log.Info("开始下载评估数据")
	taskInfo.TaskStatus = ct.TrainJobStatusDownloadTrainData.Code
	if err := etu.editTaskAndSendMq(ctx, taskInfo); err != nil {
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	datasetPath := filepath.Join(taskPath, file.DATASET)
	utils.EnsureDirectoryExists(datasetPath)

	datasetDataPath := filepath.Join(datasetPath, file.DATASET_DATA)
	utils.EnsureDirectoryExists(datasetDataPath)

	if len(eventInfo.DatasetFolders) == 0 {
		return errors.New("数据集文件夹列表不能为空")
	}

	for _, folder := range eventInfo.DatasetFolders {
		etu.log.Infof("从文件夹下载数据集: %s", folder)
		files, err := etu.oss.ListFiles(ctx, eventInfo.DatasetBucket, folder)
		if err != nil {
			return fmt.Errorf("列出数据集文件失败: %v", err)
		}

		for _, file := range files {
			err = etu.oss.DownloadSingleFile(
				ctx,
				eventInfo.DatasetBucket,
				file,
				datasetDataPath,
				filepath.Base(file),
			)
			if err != nil {
				return fmt.Errorf("下载数据集文件失败: %v", err)
			}
		}
	}

	if eventInfo.DataZip {
		etu.log.Info("解压数据集文件")
		files, err := filepath.Glob(filepath.Join(datasetDataPath, "*.zip"))
		if err == nil && len(files) > 0 {
			for _, zipFile := range files {
				if err := utils.Unzip(zipFile, datasetDataPath); err != nil {
					etu.log.Warnf("解压数据集文件失败: %v", err)
				}
			}
		}
	}

	// 下载标注文件
	etu.log.Info("开始下载标注文件")
	datasetAnnotationPath := filepath.Join(datasetPath, file.DATASET_ANNOTATION)
	utils.EnsureDirectoryExists(datasetAnnotationPath)

	if len(eventInfo.AnnotationFolders) == 0 {
		return errors.New("标注文件夹列表不能为空")
	}

	for _, folder := range eventInfo.AnnotationFolders {
		etu.log.Infof("从文件夹下载标注: %s", folder)
		files, err := etu.oss.ListFiles(ctx, eventInfo.DatasetBucket, folder)
		if err != nil {
			return fmt.Errorf("列出标注文件失败: %v", err)
		}

		for _, file := range files {
			err = etu.oss.DownloadSingleFile(
				ctx,
				eventInfo.DatasetBucket,
				file,
				datasetAnnotationPath,
				filepath.Base(file),
			)
			if err != nil {
				return fmt.Errorf("下载标注文件失败: %v", err)
			}
		}
	}

	if eventInfo.AnnotationZip {
		etu.log.Info("解压标注文件")
		files, err := filepath.Glob(filepath.Join(datasetAnnotationPath, "*.zip"))
		if err == nil && len(files) > 0 {
			for _, zipFile := range files {
				if err := utils.Unzip(zipFile, datasetAnnotationPath); err != nil {
					etu.log.Warnf("解压标注文件失败: %v", err)
				}
			}
		}
	}

	// 下载需要评估的模型或检查点
	if eventInfo.ModelOrCheckpointBucket == "" || eventInfo.ModelOrCheckpointFileURL == "" {
		return errors.New("模型或检查点的存储桶或URL不能为空")
	}

	if eventInfo.EvalType == "1" { // 模型评估类型
		etu.log.Infof("下载模型, taskId: %s", eventInfo.TaskId)
		taskInfo.TaskStatus = ct.TrainJobStatusDownloadPreWeights.Code
		if err := etu.editTaskAndSendMq(ctx, taskInfo); err != nil {
			return fmt.Errorf("更新任务状态失败: %v", err)
		}

		modelDir := filepath.Join(taskPath, file.MODEL)
		utils.EnsureDirectoryExists(modelDir)

		err = etu.oss.DownloadSingleFile(
			ctx,
			eventInfo.ModelOrCheckpointBucket,
			eventInfo.ModelOrCheckpointFileURL,
			modelDir,
			filepath.Base(eventInfo.ModelOrCheckpointFileURL),
		)
		if err != nil {
			return fmt.Errorf("下载模型文件失败: %v", err)
		}
	} else if eventInfo.EvalType == "2" { // 检查点评估类型
		etu.log.Infof("下载检查点, taskId: %s", eventInfo.TaskId)

		checkpointDir := filepath.Join(taskPath, file.CHECKPOINT)
		utils.EnsureDirectoryExists(checkpointDir)

		err = etu.oss.DownloadSingleFile(
			ctx,
			eventInfo.ModelOrCheckpointBucket,
			eventInfo.ModelOrCheckpointFileURL,
			checkpointDir,
			filepath.Base(eventInfo.ModelOrCheckpointFileURL),
		)
		if err != nil {
			return fmt.Errorf("下载检查点文件失败: %v", err)
		}
	} else {
		return fmt.Errorf("无效的评估类型: %s", eventInfo.EvalType)
	}

	// 设置容器运行参数
	hostPath, _ := filepath.Abs(taskPath)
	containerPath := file.UNIX_SEPARATOR + file.WORKSPACE
	scriptPath := containerPath + file.UNIX_SEPARATOR + file.SCRIPT + file.UNIX_SEPARATOR + etu.evalScriptName

	etu.log.Info("设置容器参数")
	args := []string{
		taskArgs.ArgTaskID, eventInfo.TaskId,
	}

	// 添加评估任务参数
	args = append(args, eventInfo.Args...)
	args = append(args,
		taskArgs.ArgAnnotationDir, containerPath+file.UNIX_SEPARATOR+file.DATASET+file.UNIX_SEPARATOR+file.DATASET_ANNOTATION+file.UNIX_SEPARATOR,
		taskArgs.ArgDataDir, containerPath+file.UNIX_SEPARATOR+file.DATASET+file.UNIX_SEPARATOR+file.DATASET_DATA+file.UNIX_SEPARATOR,
	)

	// 添加评估类型和路径
	args = append(args, taskArgs.ArgEvalType, eventInfo.EvalType)

	fileName := filepath.Base(eventInfo.ModelOrCheckpointFileURL)
	if eventInfo.EvalType == "1" { // 模型评估类型
		args = append(args,
			taskArgs.ArgModelPath, containerPath+file.UNIX_SEPARATOR+file.MODEL+file.UNIX_SEPARATOR+fileName,
		)
	} else { // 检查点评估类型
		args = append(args,
			taskArgs.ArgCheckpointPath, containerPath+file.UNIX_SEPARATOR+file.CHECKPOINT+file.UNIX_SEPARATOR+fileName,
		)
	}

	// 添加数据集标签
	if eventInfo.DatasetLabel != "" {
		utils.AddLabels(eventInfo.DatasetLabel, &args)
	}

	// 启动容器
	etu.log.Infof("使用镜像启动容器: %s", eventInfo.AlgorithmTrainImageName)
	containerInfo, err := etu.d.RunAndStartContainer(
		ctx,
		eventInfo.AlgorithmTrainImageName,
		hostPath,
		containerPath,
		scriptPath,
		args,
	)
	if err != nil {
		return fmt.Errorf("启动容器失败: %v", err)
	}

	// 启动日志流
	etu.d.StartLogStream(ctx, containerInfo.ContainerId, func(logText string) {
		etu.sendDockerLogData(ctx, eventInfo.TaskId, logText)
	})

	// 更新任务状态为"运行中"
	etu.log.Infof("容器成功启动，更新任务状态为'运行中'，容器名称: %s", containerInfo.ContainerName)
	taskInfo.TaskStatus = ct.TrainJobStatusRunning.Code
	taskInfo.TrainingContainerName = containerInfo.ContainerName
	return etu.editTaskAndSendMq(ctx, taskInfo)
}

// JustStop 只停止容器，不删除容器和相关文件
func (etu *EvalTaskUsecase) JustStop(ctx context.Context, taskId string, remove bool) {
	etu.log.Infof("开始停止评估任务，任务ID: %s", taskId)
	task := etu.etm.FindTaskById(ctx, taskId)
	if task == nil {
		etu.log.Warnf("停止失败，未找到评估任务。taskId: %s", taskId)
		return
	}
	etu.log.Infof("找到评估任务信息: taskId=%s, containerName=%s", taskId, task.TrainingContainerName)

	// 停止日志
	containerName := task.TrainingContainerName
	if containerName != "" {
		etu.d.StopLogStream(ctx, containerName, func(logText string) {
			etu.sendDockerLogData(ctx, taskId, logText)
		})

		// 停止容器
		etu.log.Infof("停止容器: %s, remove=%v", containerName, remove)
		err := etu.d.StopContainerByName(ctx, containerName, remove)
		if err != nil {
			etu.log.Errorf("停止容器失败: %v", err)
		} else {
			etu.log.Infof("容器已成功停止: %s", containerName)
		}
	}
}

// CleanupTaskDirectoryAndRecord 清理任务相关的目录和记录
func (etu *EvalTaskUsecase) CleanupTaskDirectoryAndRecord(ctx context.Context, taskId string) {
	// 删除目录
	storePath := filepath.Join(etu.filePath, file.EVAL)
	taskPath := filepath.Join(storePath, taskId)
	etu.log.Infof("删除评估任务目录: %s", taskPath)
	err := utils.RemoveDirectory(taskPath)
	if err != nil {
		etu.log.Errorf("删除评估任务目录失败。taskId: %s, 错误: %v", taskId, err)
	} else {
		etu.log.Infof("删除评估任务目录成功。taskId: %s", taskId)
	}

	// 删除记录
	etu.log.Infof("从管理器中移除评估任务, taskId: %s", taskId)
	etu.etm.RemoveTask(ctx, taskId)
	etu.log.Infof("评估任务已成功从管理器中移除, taskId: %s", taskId)
}

// DestroyAndDelete 销毁容器，并删除相关文件
func (etu *EvalTaskUsecase) DestroyAndDelete(ctx context.Context, taskId string) {
	etu.log.Infof("开始销毁和删除评估任务，任务ID: %s", taskId)

	// 停止容器
	etu.JustStop(ctx, taskId, true)
	etu.log.Infof("容器已成功停止并移除, taskId: %s", taskId)

	// 清理任务相关的目录和记录
	etu.CleanupTaskDirectoryAndRecord(ctx, taskId)
	etu.log.Infof("评估任务已成功销毁和删除, taskId: %s", taskId)
}

// BatchInfoHandle 处理评估任务批次信息
func (etu *EvalTaskUsecase) BatchInfoHandle(ctx context.Context, batchInfo *eval.EvalBatchInfo) {
	etu.log.Infof("处理评估任务批次信息, taskId: %s, 详情数量: %d", batchInfo.TaskId, len(batchInfo.Details))

	taskInfo := etu.etm.FindTaskById(ctx, batchInfo.TaskId)
	if taskInfo == nil {
		etu.log.Errorf("处理批次信息失败，未找到任务。taskId: %s", batchInfo.TaskId)
		return
	}

	details := batchInfo.Details
	if len(details) == 0 {
		etu.log.Warnf("评估任务批次详情为空, taskId: %s", batchInfo.TaskId)
		return
	}

	// 构建详情列表
	var detailList []event.EvalDetail
	for _, detail := range details {
		detailList = append(detailList, event.EvalDetail{
			DataUUID: detail.DataUuid,
			EvalData: detail.EvalData,
		})
	}

	// 发送批次信息
	reply := &event.EvalReceiveMessage{
		TaskId:     batchInfo.TaskId,
		Status:     ct.TrainJobStatusRunning.Code,
		DetailList: detailList,
	}

	if err := etu.mq.SendToService(ctx, etu.tsn, &event.ReqMessage{
		Type:    mq.TASK_EVALUATE.Code(),
		Payload: reply,
	}); err != nil {
		etu.log.Errorf("发送批次信息失败: %v", err)
	}
}

// FinishHandle 处理评估任务完成
func (etu *EvalTaskUsecase) FinishHandle(ctx context.Context, result *eval.EvalTaskResult) {
	etu.log.Infof("处理评估任务完成, taskId: %s", result.TaskId)

	taskInfo := etu.etm.FindTaskById(ctx, result.TaskId)
	if taskInfo == nil {
		etu.log.Errorf("处理评估完成失败，未找到任务。taskId: %s", result.TaskId)
		return
	}

	// 更新任务状态
	taskInfo.TaskStatus = ct.TrainJobStatusSucceed.Code
	etu.etm.UpdateTask(ctx, taskInfo)

	// 发送评估完成信息
	resultJSON, err := utils.ToJSON(result)
	if err != nil {
		etu.log.Errorf("转换评估结果为JSON失败: %v", err)
		return
	}

	reply := &event.EvalReceiveMessage{
		TaskId: result.TaskId,
		Status: ct.TrainJobStatusSucceed.Code,
		Result: resultJSON,
	}

	if err := etu.mq.SendToService(ctx, etu.tsn, &event.ReqMessage{
		Type:    mq.TASK_EVALUATE.Code(),
		Payload: reply,
	}); err != nil {
		etu.log.Errorf("发送评估完成信息失败: %v", err)
		return
	}

	// 停止任务
	etu.JustStop(ctx, result.TaskId, false)
}

// HandleEvalTask 处理评估任务消息
func (etu *EvalTaskUsecase) HandleEvalTask(ctx context.Context, task *event.EvalSendMessage) {
	taskId := task.TaskId
	taskInfo := eval.NewEvalTaskInfoWithMessage(task)

	if task.Op == ct.TrainJobOpStart.Code {
		etu.log.Infof("收到启动评估任务请求，任务ID: %s", taskId)
		err := etu.StartEvaluation(ctx, taskInfo)
		if err != nil {
			etu.log.Errorf("启动评估任务失败! taskId: %s, 错误: %v", taskId, err)
			taskInfo.TaskStatus = ct.TrainJobStatusFail.Code
			taskInfo.Remark = "启动错误: " + err.Error()
			etu.editTaskAndSendMq(ctx, taskInfo)
			etu.JustStop(ctx, taskId, false)
		}
	} else if task.Op == ct.TrainJobOpStop.Code {
		etu.log.Infof("收到停止评估任务请求，任务ID: %s", taskId)
		existingTask := etu.etm.FindTaskById(ctx, taskId)
		if existingTask == nil {
			etu.log.Warnf("停止失败，未找到任务。taskId: %s", taskId)
			emptyTask := &eval.EvalTaskInfo{
				TaskId:     taskId,
				TaskStatus: ct.TrainJobStatusFail.Code,
				Remark:     "未找到任务!",
			}
			etu.sendStatusChangeMessage(ctx, emptyTask)
			return
		}

		etu.JustStop(ctx, taskId, false)
		existingTask.TaskStatus = ct.TrainJobStatusStopped.Code
		etu.sendStatusChangeMessage(ctx, existingTask)
	} else {
		etu.log.Errorf("未知的任务操作! op: %d, taskId: %s", task.Op, taskId)
	}
}

// CheckTask 检查所有任务的状态
func (etu *EvalTaskUsecase) CheckTask(ctx context.Context) {
	etu.log.Debug("开始检查评估任务状态...")
	tasks := etu.etm.GetTaskList(ctx)
	if len(tasks) == 0 {
		etu.log.Debug("没有需要检查的评估任务")
		return
	}

	for _, task := range tasks {
		// 使用匿名函数避免循环中的错误影响到所有任务的检查
		func(task *eval.EvalTaskInfo) {
			defer func() {
				if r := recover(); r != nil {
					etu.log.Errorf("检查评估任务时发生异常: %v", r)
				}
			}()

			etu.checkSingleTask(ctx, task)
		}(task)
	}
	etu.log.Info("评估任务状态检查完成")
}

// checkSingleTask 检查单个任务的状态
func (etu *EvalTaskUsecase) checkSingleTask(ctx context.Context, task *eval.EvalTaskInfo) {
	taskId := task.TaskId
	status := task.TaskStatus

	// 只处理特定状态的任务
	if !etu.isStatusNeedCheck(status) {
		return
	}

	containerName := task.TrainingContainerName
	if containerName == "" {
		etu.log.Errorf("任务 %s 没有关联的容器名称，清理任务目录和记录", taskId)
		etu.CleanupTaskDirectoryAndRecord(ctx, taskId)
		return
	}

	containerInfo, err := etu.d.FindContainerByName(ctx, containerName)
	if err != nil {
		etu.log.Errorf("查找容器信息失败: %v", err)
		return
	}

	if status == ct.TrainJobStatusRunning.Code {
		etu.handleRunningTask(ctx, task, containerInfo)
	} else {
		etu.handleFinishedTask(ctx, task, containerInfo)
	}
}

// isStatusNeedCheck 判断任务状态是否需要检查
func (etu *EvalTaskUsecase) isStatusNeedCheck(status int32) bool {
	return status == ct.TrainJobStatusRunning.Code ||
		status == ct.TrainJobStatusStopped.Code ||
		status == ct.TrainJobStatusFail.Code ||
		status == ct.TrainJobStatusSucceed.Code
}

// handleRunningTask 处理运行中的任务状态
func (etu *EvalTaskUsecase) handleRunningTask(
	ctx context.Context,
	task *eval.EvalTaskInfo,
	containerInfo *mc.ContainerInfo,
) {
	taskId := task.TaskId

	if containerInfo == nil {
		etu.log.Errorf("任务 %s 标记为运行中但容器不存在，清理任务目录和记录", taskId)
		task.TaskStatus = ct.TrainJobStatusFail.Code
		task.Remark = "容器不存在!"
		etu.editTaskAndSendMq(ctx, task)
		etu.CleanupTaskDirectoryAndRecord(ctx, taskId)
		return
	}

	// 获取容器状态
	containerState, err := etu.d.GetContainerState(ctx, containerInfo.ContainerId)
	if err != nil {
		etu.log.Errorf("获取容器状态失败: %v", err)
		return
	}

	if containerState.State != string(cc.RUNNING) {
		etu.log.Errorf("任务 %s 标记为运行中但容器已停止", taskId)
		logs, err := etu.d.GetContainerLastLogs(ctx, containerInfo.ContainerId, 10)
		if err != nil {
			etu.log.Errorf("获取容器日志失败: %v", err)
		}

		task.TaskStatus = ct.TrainJobStatusFail.Code
		task.Remark = "容器意外停止: " + logs
		etu.editTaskAndSendMq(ctx, task)
	}
}

// handleFinishedTask 处理已完成的任务状态
func (etu *EvalTaskUsecase) handleFinishedTask(ctx context.Context, task *eval.EvalTaskInfo, containerInfo *mc.ContainerInfo) {
	taskId := task.TaskId
	status := task.TaskStatus

	if containerInfo == nil {
		etu.log.Infof("任务 %s (状态:%d) 没有关联的容器", taskId, status)
		etu.CleanupTaskDirectoryAndRecord(ctx, taskId)
		return
	}

	// 获取容器状态
	containerState, err := etu.d.GetContainerState(ctx, containerInfo.ContainerId)
	if err != nil {
		etu.log.Errorf("获取容器状态失败: %v", err)
		return
	}

	// 如果容器还在运行，先停止容器
	if containerState.State == string(cc.RUNNING) {
		etu.log.Infof("任务 %s 的容器仍在运行，先停止容器", taskId)
		etu.d.StopContainerByName(ctx, task.TrainingContainerName, false)
		return
	}

	etu.checkAndCleanupStoppedContainer(ctx, task, containerInfo)
}

// checkAndCleanupStoppedContainer 检查并清理停止的容器
func (etu *EvalTaskUsecase) checkAndCleanupStoppedContainer(
	ctx context.Context,
	task *eval.EvalTaskInfo,
	containerInfo *mc.ContainerInfo,
) {
	taskId := task.TaskId
	status := task.TaskStatus
	stopTime, err := etu.d.GetContainerStopTime(ctx, containerInfo.ContainerId)
	if err != nil {
		etu.log.Errorf("获取容器停止时间失败: %v", err)
		return
	}

	currentTime := time.Now().Unix()
	timeDifferenceHours := float64(currentTime-stopTime) / 3600.0

	etu.log.Infof("任务 %s (状态:%d) 的容器已停止 %.2f 小时",
		taskId,
		status,
		timeDifferenceHours)

	// 如果停止超过6小时，清理任务和容器
	if timeDifferenceHours >= 6.0 {
		etu.log.Infof("清理任务 %s (状态:%d) 的容器，停止时间超过6小时",
			taskId,
			status)
		etu.DestroyAndDelete(ctx, taskId)
	}
}

// 发送日志数据
func (etu *EvalTaskUsecase) sendDockerLogData(ctx context.Context, taskId string, logText string) {
	taskIdInt, err := strconv.ParseInt(taskId, 10, 64)
	if err != nil {
		etu.log.Errorf("转换任务ID失败: %v", err)
		return
	}

	logMsg := &event.DockerLogRespMessage{
		TaskId:   taskIdInt,
		Log:      logText,
		TaskType: 1, // 评估类型
	}

	if err := etu.mq.SendToService(ctx, etu.tsn, &event.ReqMessage{
		Type:    mq.DOCKER_LOG.Code(),
		Payload: logMsg,
	}); err != nil {
		etu.log.Errorf("发送日志消息失败: %v", err)
	}
}
