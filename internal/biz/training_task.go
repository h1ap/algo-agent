package biz

import (
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
	"algo-agent/internal/model/train"
	"algo-agent/internal/mq/event"
	"algo-agent/internal/utils"

	"github.com/go-kratos/kratos/v2/log"
)

// TrainingTaskUsecase 训练任务用例
type TrainingTaskUsecase struct {
	ttm TrainingTaskManager

	mq  MqService
	d   DockerService
	oss OSSService

	log *log.Helper

	filePath string
	tsn      string // 训练服务名称

	// 路径常量
	checkpointPathPrefix string
	modelPathPrefix      string
	trainScriptName      string
}

// 更新状态并发送消息
func (ttu *TrainingTaskUsecase) editTaskAndSendMq(ctx context.Context, task *train.TrainingTaskInfo) error {
	ttu.ttm.UpdateTask(ctx, task)
	return ttu.sendStatusChangeMessage(ctx, task)
}

// 发送状态变更消息
func (ttu *TrainingTaskUsecase) sendStatusChangeMessage(ctx context.Context, taskInfo *train.TrainingTaskInfo) error {
	reply := &event.TrainTaskRespMessage{
		TaskId:     taskInfo.TaskId,
		TaskStatus: taskInfo.TaskStatus,
		Remark:     taskInfo.Remark,
	}

	jsonStr, err := utils.ToJSON(reply)
	if err != nil {
		return fmt.Errorf("转换为JSON失败: %v", err)
	}

	return ttu.mq.SendToService(ctx, ttu.tsn, jsonStr)
}

// StartTraining 开始训练任务
func (ttu *TrainingTaskUsecase) StartTraining(ctx context.Context, taskInfo *train.TrainingTaskInfo) error {
	eventInfo := taskInfo.TrainTaskReqMessage
	if eventInfo == nil || eventInfo.TaskId == "" {
		ttu.log.Error("开始训练失败：任务ID为空")
		return errors.New("任务ID不允许为空")
	}
	ttu.log.Infof("开始训练任务，任务ID: %s", eventInfo.TaskId)

	// 设置任务状态为"训练中"
	taskInfo.TaskStatus = ct.TrainJobStatusStarting.Code
	if err := ttu.ttm.AddTask(ctx, taskInfo); err != nil {
		return fmt.Errorf("添加任务失败: %v", err)
	}
	// 发送任务开始的状态到系统消息队列
	if err := ttu.sendStatusChangeMessage(ctx, taskInfo); err != nil {
		return fmt.Errorf("发送状态变更消息失败: %v", err)
	}

	// 创建必要的目录
	trainBasePath := filepath.Join(ttu.filePath, file.TRAIN)
	imagePath := filepath.Join(ttu.filePath, file.IMAGE)
	utils.EnsureDirectoryExists(trainBasePath)
	utils.EnsureDirectoryExists(imagePath)

	// 检查Docker镜像是否存在
	ttu.log.Infof("检查训练镜像是否存在: %s", eventInfo.AlgorithmTrainImageName)
	imageExists, err := ttu.d.FindImageByName(ctx, eventInfo.AlgorithmTrainImageName)
	if err != nil || !imageExists {
		ttu.log.Info("本地未找到训练镜像，准备下载")
		taskInfo.TaskStatus = ct.TrainJobStatusDownloadTrainImage.Code
		if err := ttu.editTaskAndSendMq(ctx, taskInfo); err != nil {
			return fmt.Errorf("更新任务状态失败: %v", err)
		}

		// 下载镜像文件
		imageFileName := strings.Replace(eventInfo.AlgorithmTrainImageName, ":", "-", -1)
		tarPath := filepath.Join(imagePath, imageFileName+".tar")
		err = ttu.oss.DownloadSingleFile(
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
		err = ttu.d.ImportAndTagImage(ctx, tarPath, eventInfo.AlgorithmTrainImageName)
		if err != nil {
			return fmt.Errorf("导入镜像失败: %v", err)
		}
	} else {
		ttu.log.Info("找到已存在的训练镜像")
	}

	// 设置任务路径
	taskPath := filepath.Join(trainBasePath, eventInfo.TaskId)

	// 下载算法脚本
	ttu.log.Info("开始下载算法脚本")
	taskInfo.TaskStatus = ct.TrainJobStatusDownloadAlgoScripts.Code
	if err := ttu.editTaskAndSendMq(ctx, taskInfo); err != nil {
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	algorithmScriptZipFile := "script_" + eventInfo.TaskId + ".zip"
	// 下载算法脚本
	err = ttu.oss.DownloadSingleFile(
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
	ttu.log.Info("解压算法脚本")
	scriptZipPath := filepath.Join(taskPath, algorithmScriptZipFile)
	scriptDestPath := filepath.Join(taskPath, file.SCRIPT)
	if err := utils.Unzip(scriptZipPath, scriptDestPath); err != nil {
		return fmt.Errorf("解压算法脚本失败: %v", err)
	}

	// 下载训练数据
	ttu.log.Info("开始下载训练数据")
	taskInfo.TaskStatus = ct.TrainJobStatusDownloadTrainData.Code
	if err := ttu.editTaskAndSendMq(ctx, taskInfo); err != nil {
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
		ttu.log.Infof("从文件夹下载数据集: %s", folder)
		files, err := ttu.oss.ListFiles(ctx, eventInfo.DatasetBucket, folder)
		if err != nil {
			return fmt.Errorf("列出数据集文件失败: %v", err)
		}

		for _, file := range files {
			err = ttu.oss.DownloadSingleFile(
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
		ttu.log.Info("解压数据集文件")
		files, err := filepath.Glob(filepath.Join(datasetDataPath, "*.zip"))
		if err == nil && len(files) > 0 {
			for _, zipFile := range files {
				if err := utils.Unzip(zipFile, datasetDataPath); err != nil {
					ttu.log.Warnf("解压数据集文件失败: %v", err)
				}
			}
		}
	}

	// 下载标注文件
	ttu.log.Info("开始下载标注文件")
	datasetAnnotationPath := filepath.Join(datasetPath, file.DATASET_ANNOTATION)
	utils.EnsureDirectoryExists(datasetAnnotationPath)

	if len(eventInfo.AnnotationFolders) == 0 {
		return errors.New("标注文件夹列表不能为空")
	}

	for _, folder := range eventInfo.AnnotationFolders {
		ttu.log.Infof("从文件夹下载标注: %s", folder)
		files, err := ttu.oss.ListFiles(ctx, eventInfo.DatasetBucket, folder)
		if err != nil {
			return fmt.Errorf("列出标注文件失败: %v", err)
		}

		for _, file := range files {
			err = ttu.oss.DownloadSingleFile(
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
		ttu.log.Info("解压标注文件")
		files, err := filepath.Glob(filepath.Join(datasetAnnotationPath, "*.zip"))
		if err == nil && len(files) > 0 {
			for _, zipFile := range files {
				if err := utils.Unzip(zipFile, datasetAnnotationPath); err != nil {
					ttu.log.Warnf("解压标注文件失败: %v", err)
				}
			}
		}
	}

	// 下载初始权重
	ttu.log.Info("检查预训练权重")
	modelDir := filepath.Join(taskPath, file.MODEL)
	utils.EnsureDirectoryExists(modelDir)
	preModel := false
	if eventInfo.PreModelBucket != "" && eventInfo.PreModelFileURL != "" {
		ttu.log.Info("下载预训练权重")
		preModel = true
		taskInfo.TaskStatus = ct.TrainJobStatusDownloadPreWeights.Code
		if err := ttu.editTaskAndSendMq(ctx, taskInfo); err != nil {
			return fmt.Errorf("更新任务状态失败: %v", err)
		}

		err = ttu.oss.DownloadSingleFile(
			ctx,
			eventInfo.PreModelBucket,
			eventInfo.PreModelFileURL,
			modelDir,
			filepath.Base(eventInfo.PreModelFileURL),
		)
		if err != nil {
			return fmt.Errorf("下载预训练权重失败: %v", err)
		}
	}

	// 设置容器运行参数
	hostPath, _ := filepath.Abs(taskPath)
	containerPath := file.UNIX_SEPARATOR + file.WORKSPACE
	scriptPath := containerPath + file.UNIX_SEPARATOR + file.SCRIPT + file.UNIX_SEPARATOR + ttu.trainScriptName

	ttu.log.Info("设置容器参数")
	args := []string{
		taskArgs.ArgTaskID, eventInfo.TaskId,
	}

	// 添加训练任务参数
	args = append(args, eventInfo.Args...)
	args = append(args,
		taskArgs.ArgAnnotationDir, containerPath+file.UNIX_SEPARATOR+file.DATASET+file.UNIX_SEPARATOR+file.DATASET_ANNOTATION+file.UNIX_SEPARATOR,
		taskArgs.ArgDataDir, containerPath+file.UNIX_SEPARATOR+file.DATASET+file.UNIX_SEPARATOR+file.DATASET_DATA+file.UNIX_SEPARATOR,
	)

	if preModel {
		fileName := filepath.Base(eventInfo.PreModelFileURL)
		args = append(args,
			taskArgs.ArgPreWeightPath, containerPath+file.UNIX_SEPARATOR+file.MODEL+file.UNIX_SEPARATOR+fileName,
		)
	}

	// 添加数据集标签
	if eventInfo.DatasetLabel != "" {
		utils.AddLabels(eventInfo.DatasetLabel, &args)
	}

	// 启动容器
	ttu.log.Infof("使用镜像启动容器: %s", eventInfo.AlgorithmTrainImageName)
	containerInfo, err := ttu.d.RunAndStartContainer(
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
	ttu.d.StartLogStream(ctx, containerInfo.ContainerId, func(logText string) {
		ttu.sendDockerLogData(ctx, eventInfo.TaskId, logText)
	})

	// 更新任务状态为"运行中"
	ttu.log.Infof("容器成功启动，更新任务状态为'运行中'，容器名称: %s", containerInfo.ContainerName)
	taskInfo.TaskStatus = ct.TrainJobStatusRunning.Code
	taskInfo.TrainingContainerName = containerInfo.ContainerName
	return ttu.editTaskAndSendMq(ctx, taskInfo)
}

// 发送日志数据
func (ttu *TrainingTaskUsecase) sendDockerLogData(ctx context.Context, taskId string, logText string) {
	taskIdInt, err := strconv.ParseInt(taskId, 10, 64)
	if err != nil {
		ttu.log.Errorf("转换任务ID失败: %v", err)
		return
	}

	logMsg := &event.DockerLogRespMessage{
		TaskId:   taskIdInt,
		Log:      logText,
		TaskType: 0, // 训练类型
	}

	jsonStr, err := utils.ToJSON(logMsg)
	if err != nil {
		ttu.log.Errorf("转换日志消息为JSON失败: %v", err)
		return
	}

	if err := ttu.mq.SendToService(ctx, ttu.tsn, jsonStr); err != nil {
		ttu.log.Errorf("发送日志消息失败: %v", err)
	}
}

// JustStop 只停止容器，不删除容器和相关文件
func (ttu *TrainingTaskUsecase) JustStop(ctx context.Context, taskId string, remove bool) {
	ttu.log.Infof("开始停止任务，任务ID: %s", taskId)
	task := ttu.ttm.FindTaskById(ctx, taskId)
	if task == nil {
		ttu.log.Warnf("停止失败，未找到任务。taskId: %s", taskId)
		return
	}
	ttu.log.Infof("找到任务信息: taskId=%s, containerName=%s", taskId, task.TrainingContainerName)

	// 停止日志
	containerName := task.TrainingContainerName
	if containerName != "" {
		ttu.d.StopLogStream(ctx, containerName, func(logText string) {
			ttu.sendDockerLogData(ctx, taskId, logText)
		})

		// 停止容器
		ttu.log.Infof("停止容器: %s, remove=%v", containerName, remove)
		err := ttu.d.StopContainerByName(ctx, containerName, remove)
		if err != nil {
			ttu.log.Errorf("停止容器失败: %v", err)
		} else {
			ttu.log.Infof("容器已成功停止: %s", containerName)
		}
	}
}

// CleanupTaskDirectoryAndRecord 清理任务相关的目录和记录
func (ttu *TrainingTaskUsecase) CleanupTaskDirectoryAndRecord(ctx context.Context, taskId string) {
	// 删除目录
	storePath := filepath.Join(ttu.filePath, file.TRAIN)
	taskPath := filepath.Join(storePath, taskId)
	ttu.log.Infof("删除任务目录: %s", taskPath)
	err := utils.RemoveDirectory(taskPath)
	if err != nil {
		ttu.log.Errorf("删除训练任务目录失败。taskId: %s, 错误: %v", taskId, err)
	} else {
		ttu.log.Infof("删除训练任务目录成功。taskId: %s", taskId)
	}

	// 删除记录
	ttu.log.Infof("从管理器中移除任务, taskId: %s", taskId)
	ttu.ttm.RemoveTask(ctx, taskId)
	ttu.log.Infof("任务已成功从管理器中移除, taskId: %s", taskId)
}

// DestroyAndDelete 销毁容器，并删除相关文件
func (ttu *TrainingTaskUsecase) DestroyAndDelete(ctx context.Context, taskId string) {
	ttu.log.Infof("开始销毁和删除任务，任务ID: %s", taskId)

	// 停止容器
	ttu.JustStop(ctx, taskId, true)
	ttu.log.Infof("容器已成功停止并移除, taskId: %s", taskId)

	// 清理任务相关的目录和记录
	ttu.CleanupTaskDirectoryAndRecord(ctx, taskId)
	ttu.log.Infof("任务已成功销毁和删除, taskId: %s", taskId)
}

// EpochInfoHandle 处理每轮的训练信息
func (ttu *TrainingTaskUsecase) EpochInfoHandle(ctx context.Context, epochInfo *train.TrainingEpochInfo) {
	ttu.log.Infof("处理训练任务周期信息, taskId: %s, epoch: %d", epochInfo.TaskId, epochInfo.Epoch)

	metrics := &event.Metrics{
		Epoch:            epochInfo.Epoch,
		EstimateTimeLeft: epochInfo.EstimatedTimeLeft,
		CreateTime:       time.Now(),
		DynamicFields:    make(map[string]interface{}),
	}

	for k, v := range epochInfo.DynamicFields {
		metrics.DynamicFields[k] = v
	}

	reply := &event.TrainTaskRespMessage{
		TaskId:       epochInfo.TaskId,
		TaskStatus:   ct.TrainJobStatusRunning.Code,
		Epoch:        epochInfo.Epoch,
		IsCheckpoint: false,
		Metrics:      metrics,
	}

	jsonStr, err := utils.ToJSON(reply)
	if err != nil {
		ttu.log.Errorf("转换训练周期信息为JSON失败: %v", err)
		return
	}

	if err := ttu.mq.SendToService(ctx, ttu.tsn, jsonStr); err != nil {
		ttu.log.Errorf("发送训练周期信息失败: %v", err)
	}
}

// CheckpointHandle 处理检查点
func (ttu *TrainingTaskUsecase) CheckpointHandle(ctx context.Context, checkpoint *train.TrainingCheckpoint) error {
	ttu.log.Infof("处理检查点, taskId: %s, epoch: %d", checkpoint.TaskId, checkpoint.Epoch)
	taskBasePath := filepath.Join(ttu.filePath, file.TRAIN, checkpoint.TaskId)
	taskInfo := ttu.ttm.FindTaskById(ctx, checkpoint.TaskId)
	if taskInfo == nil {
		return fmt.Errorf("处理检查点失败，未找到任务。taskId: %s", checkpoint.TaskId)
	}

	checkpointPath := checkpoint.CheckpointPath
	// 检查路径是否以 /workspace 开头，如果是，则需要去掉
	prefix := file.UNIX_SEPARATOR + file.WORKSPACE
	if strings.HasPrefix(checkpointPath, prefix) {
		ttu.log.Infof("检查点路径以 %s 开头。taskId: %s, checkpointPath: %s", prefix, checkpoint.TaskId, checkpointPath)
		checkpointPath = checkpointPath[len(prefix):]
		ttu.log.Infof("处理后的检查点路径: %s", checkpointPath)
	}

	fileName := filepath.Base(checkpointPath)
	objectUrl := ttu.checkpointPathPrefix + checkpoint.TaskId + file.UNIX_SEPARATOR + file.CHECKPOINT + file.UNIX_SEPARATOR + fileName
	path := filepath.Join(taskBasePath, checkpointPath)

	// 上传检查点文件
	err := ttu.oss.UploadFile(
		ctx,
		taskInfo.TrainTaskReqMessage.AlgorithmScriptMinioBucket,
		&File{
			Name: fileName,
			Path: path,
		},
		objectUrl,
	)
	if err != nil {
		return fmt.Errorf("上传检查点文件失败: %v", err)
	}

	// 发送检查点信息
	reply := &event.TrainTaskRespMessage{
		TaskId:             checkpoint.TaskId,
		TaskStatus:         ct.TrainJobStatusRunning.Code,
		Epoch:              checkpoint.Epoch,
		IsCheckpoint:       true,
		CheckpointFilePath: objectUrl,
	}

	jsonStr, err := utils.ToJSON(reply)
	if err != nil {
		return fmt.Errorf("转换检查点信息为JSON失败: %v", err)
	}

	return ttu.mq.SendToService(ctx, ttu.tsn, jsonStr)
}

// FinishHandle 处理训练完成
func (ttu *TrainingTaskUsecase) FinishHandle(ctx context.Context, result *train.TrainingTaskResult) error {
	ttu.log.Infof("处理训练完成, taskId: %s, bestEpoch: %d", result.TaskId, result.BestEpoch)
	taskInfo := ttu.ttm.FindTaskById(ctx, result.TaskId)
	if taskInfo == nil {
		return fmt.Errorf("处理训练完成失败，未找到任务。taskId: %s", result.TaskId)
	}

	bestPath := result.BestModelPath
	finalPath := result.FinalModelPath

	// 检查路径是否以 /workspace 开头，如果是，则需要去掉
	prefix := file.UNIX_SEPARATOR + file.WORKSPACE
	if strings.HasPrefix(bestPath, prefix) {
		ttu.log.Infof("最佳模型路径以 %s 开头。taskId: %s, bestPath: %s", prefix, result.TaskId, bestPath)
		bestPath = bestPath[len(prefix):]
		ttu.log.Infof("处理后的最佳模型路径: %s", bestPath)
	}
	if strings.HasPrefix(finalPath, prefix) {
		ttu.log.Infof("最终模型路径以 %s 开头。taskId: %s, finalPath: %s", prefix, result.TaskId, finalPath)
		finalPath = finalPath[len(prefix):]
		ttu.log.Infof("处理后的最终模型路径: %s", finalPath)
	}

	taskBasePath := filepath.Join(ttu.filePath, file.TRAIN, result.TaskId)
	bestFilePath := filepath.Join(taskBasePath, bestPath)
	finalFilePath := filepath.Join(taskBasePath, finalPath)

	// 验证文件存在
	if !utils.FileExists(bestFilePath) {
		errMsg := fmt.Sprintf("最佳模型文件不存在: %s", bestFilePath)
		ttu.log.Error(errMsg)
		taskInfo.TaskStatus = ct.TrainJobStatusFail.Code
		taskInfo.Remark = "训练失败: 最佳模型文件不存在"
		ttu.editTaskAndSendMq(ctx, taskInfo)
		ttu.JustStop(ctx, result.TaskId, false)
		return errors.New(errMsg)
	}

	if !utils.FileExists(finalFilePath) {
		errMsg := fmt.Sprintf("最终模型文件不存在: %s", finalFilePath)
		ttu.log.Error(errMsg)
		taskInfo.TaskStatus = ct.TrainJobStatusFail.Code
		taskInfo.Remark = "训练失败: 最终模型文件不存在"
		ttu.editTaskAndSendMq(ctx, taskInfo)
		ttu.JustStop(ctx, result.TaskId, false)
		return errors.New(errMsg)
	}

	bestFileName := filepath.Base(bestFilePath)
	finalFileName := filepath.Base(finalFilePath)
	bestObjectUrl := ttu.modelPathPrefix + result.TaskId + file.UNIX_SEPARATOR + bestFileName
	finalObjectUrl := ttu.modelPathPrefix + result.TaskId + file.UNIX_SEPARATOR + finalFileName

	// 上传模型文件
	err := ttu.oss.UploadFile(
		ctx,
		taskInfo.TrainTaskReqMessage.AlgorithmScriptMinioBucket,
		&File{
			Name: bestFileName,
			Path: bestFilePath,
		},
		bestObjectUrl,
	)
	if err != nil {
		errMsg := fmt.Sprintf("上传最佳模型文件失败: %v", err)
		ttu.log.Error(errMsg)
		taskInfo.TaskStatus = ct.TrainJobStatusFail.Code
		taskInfo.Remark = "训练失败: " + errMsg
		ttu.editTaskAndSendMq(ctx, taskInfo)
		ttu.JustStop(ctx, result.TaskId, false)
		return errors.New(errMsg)
	}

	err = ttu.oss.UploadFile(
		ctx,
		taskInfo.TrainTaskReqMessage.AlgorithmScriptMinioBucket,
		&File{
			Name: finalFileName,
			Path: finalFilePath,
		},
		finalObjectUrl,
	)
	if err != nil {
		errMsg := fmt.Sprintf("上传最终模型文件失败: %v", err)
		ttu.log.Error(errMsg)
		taskInfo.TaskStatus = ct.TrainJobStatusFail.Code
		taskInfo.Remark = "训练失败: " + errMsg
		ttu.editTaskAndSendMq(ctx, taskInfo)
		ttu.JustStop(ctx, result.TaskId, false)
		return errors.New(errMsg)
	}

	// 更新任务状态
	taskInfo.TaskStatus = ct.TrainJobStatusSucceed.Code
	ttu.ttm.UpdateTask(ctx, taskInfo)

	// 发送训练完成信息
	reply := &event.TrainTaskRespMessage{
		TaskId:         result.TaskId,
		TaskStatus:     ct.TrainJobStatusSucceed.Code,
		BestWeightPath: bestObjectUrl,
		LastWeightPath: finalObjectUrl,
		Epoch:          result.BestEpoch,
	}

	jsonStr, err := utils.ToJSON(reply)
	if err != nil {
		ttu.log.Errorf("转换训练完成信息为JSON失败: %v", err)
		return err
	}

	if err := ttu.mq.SendToService(ctx, ttu.tsn, jsonStr); err != nil {
		ttu.log.Errorf("发送训练完成信息失败: %v", err)
		return err
	}

	// 停止任务
	ttu.JustStop(ctx, result.TaskId, false)
	return nil
}

// HandleTrainingTask 处理训练任务消息
func (ttu *TrainingTaskUsecase) HandleTrainingTask(ctx context.Context, task *event.TrainTaskReqMessage) {
	taskId := task.TaskId
	taskInfo := train.NewTrainingTaskInfoWithMessage(task)

	if task.Op == ct.TrainJobOpStart.Code {
		ttu.log.Infof("收到启动训练任务请求，任务ID: %s", taskId)
		err := ttu.StartTraining(ctx, taskInfo)
		if err != nil {
			ttu.log.Errorf("启动训练任务失败! taskId: %s, 错误: %v", taskId, err)
			taskInfo.TaskStatus = ct.TrainJobStatusFail.Code
			taskInfo.Remark = "启动错误: " + err.Error()
			ttu.editTaskAndSendMq(ctx, taskInfo)
			ttu.JustStop(ctx, taskId, false)
		}
	} else if task.Op == ct.TrainJobOpStop.Code {
		ttu.log.Infof("收到停止训练任务请求，任务ID: %s", taskId)
		existingTask := ttu.ttm.FindTaskById(ctx, taskId)
		if existingTask == nil {
			ttu.log.Warnf("停止失败，未找到任务。taskId: %s", taskId)
			emptyTask := &train.TrainingTaskInfo{
				TaskId:     taskId,
				TaskStatus: ct.TrainJobStatusFail.Code,
				Remark:     "未找到任务!",
			}
			ttu.sendStatusChangeMessage(ctx, emptyTask)
			return
		}

		ttu.JustStop(ctx, taskId, false)
		existingTask.TaskStatus = ct.TrainJobStatusStopped.Code
		ttu.sendStatusChangeMessage(ctx, existingTask)
	} else {
		ttu.log.Errorf("未知的任务操作! op: %d, taskId: %s", task.Op, taskId)
	}
}

// CheckTask 检查所有任务的状态
func (ttu *TrainingTaskUsecase) CheckTask(ctx context.Context) {
	ttu.log.Info("开始检查任务状态...")
	tasks := ttu.ttm.GetTaskList(ctx)
	if len(tasks) == 0 {
		ttu.log.Info("没有需要检查的任务")
		return
	}

	for _, task := range tasks {
		// 使用匿名函数避免循环中的错误影响到所有任务的检查
		func(task *train.TrainingTaskInfo) {
			defer func() {
				if r := recover(); r != nil {
					ttu.log.Errorf("检查任务时发生异常: %v", r)
				}
			}()

			ttu.checkSingleTask(ctx, task)
		}(task)
	}
	ttu.log.Info("任务状态检查完成")
}

// checkSingleTask 检查单个任务的状态
func (ttu *TrainingTaskUsecase) checkSingleTask(ctx context.Context, task *train.TrainingTaskInfo) {
	taskId := task.TaskId
	status := task.TaskStatus

	// 只处理特定状态的任务
	if !ttu.isStatusNeedCheck(status) {
		return
	}

	containerName := task.TrainingContainerName
	if containerName == "" {
		ttu.log.Errorf("任务 %s 没有关联的容器名称，清理任务目录和记录", taskId)
		ttu.CleanupTaskDirectoryAndRecord(ctx, taskId)
		return
	}

	containerInfo, err := ttu.d.FindContainerByName(ctx, containerName)
	if err != nil {
		ttu.log.Errorf("查找容器信息失败: %v", err)
		return
	}

	if status == ct.TrainJobStatusRunning.Code {
		ttu.handleRunningTask(ctx, task, containerInfo)
	} else {
		ttu.handleFinishedTask(ctx, task, containerInfo)
	}
}

// isStatusNeedCheck 判断任务状态是否需要检查
func (ttu *TrainingTaskUsecase) isStatusNeedCheck(status int32) bool {
	return status == ct.TrainJobStatusRunning.Code ||
		status == ct.TrainJobStatusStopped.Code ||
		status == ct.TrainJobStatusFail.Code ||
		status == ct.TrainJobStatusSucceed.Code
}

// handleRunningTask 处理运行中的任务状态
func (ttu *TrainingTaskUsecase) handleRunningTask(
	ctx context.Context,
	task *train.TrainingTaskInfo,
	containerInfo *mc.ContainerInfo,
) {
	taskId := task.TaskId

	if containerInfo == nil {
		ttu.log.Errorf("任务 %s 标记为运行中但容器不存在，清理任务目录和记录", taskId)
		task.TaskStatus = ct.TrainJobStatusFail.Code
		task.Remark = "容器不存在!"
		ttu.editTaskAndSendMq(ctx, task)
		ttu.CleanupTaskDirectoryAndRecord(ctx, taskId)
		return
	}

	// 获取容器状态
	containerState, err := ttu.d.GetContainerState(ctx, containerInfo.ContainerId)
	if err != nil {
		ttu.log.Errorf("获取容器状态失败: %v", err)
		return
	}

	if containerState.State != string(cc.RUNNING) {
		ttu.log.Errorf("任务 %s 标记为运行中但容器已停止", taskId)
		logs, err := ttu.d.GetContainerLastLogs(ctx, containerInfo.ContainerId, 10)
		if err != nil {
			ttu.log.Errorf("获取容器日志失败: %v", err)
		}

		task.TaskStatus = ct.TrainJobStatusFail.Code
		task.Remark = "容器意外停止: " + logs
		ttu.editTaskAndSendMq(ctx, task)
	}
}

// handleFinishedTask 处理已完成的任务状态
func (ttu *TrainingTaskUsecase) handleFinishedTask(ctx context.Context, task *train.TrainingTaskInfo, containerInfo *mc.ContainerInfo) {
	taskId := task.TaskId
	status := task.TaskStatus

	if containerInfo == nil {
		ttu.log.Infof("任务 %s (状态:%d) 没有关联的容器", taskId, status)
		ttu.CleanupTaskDirectoryAndRecord(ctx, taskId)
		return
	}

	// 获取容器状态
	containerState, err := ttu.d.GetContainerState(ctx, containerInfo.ContainerId)
	if err != nil {
		ttu.log.Errorf("获取容器状态失败: %v", err)
		return
	}

	// 如果容器还在运行，先停止容器
	if containerState.State == string(cc.RUNNING) {
		ttu.log.Infof("任务 %s 的容器仍在运行，先停止容器", taskId)
		ttu.d.StopContainerByName(ctx, task.TrainingContainerName, false)
		return
	}

	ttu.checkAndCleanupStoppedContainer(ctx, task, containerInfo)
}

// checkAndCleanupStoppedContainer 检查并清理停止的容器
func (ttu *TrainingTaskUsecase) checkAndCleanupStoppedContainer(
	ctx context.Context,
	task *train.TrainingTaskInfo,
	containerInfo *mc.ContainerInfo,
) {
	taskId := task.TaskId
	status := task.TaskStatus
	stopTime, err := ttu.d.GetContainerStopTime(ctx, containerInfo.ContainerId)
	if err != nil {
		ttu.log.Errorf("获取容器停止时间失败: %v", err)
		return
	}

	currentTime := time.Now().Unix()
	timeDifferenceHours := float64(currentTime-stopTime) / 3600.0

	ttu.log.Infof("任务 %s (状态:%d) 的容器已停止 %.2f 小时",
		taskId,
		status,
		timeDifferenceHours)

	// 如果停止超过6小时，清理任务和容器
	if timeDifferenceHours >= 6.0 {
		ttu.log.Infof("清理任务 %s (状态:%d) 的容器，停止时间超过6小时",
			taskId,
			status)
		ttu.DestroyAndDelete(ctx, taskId)
	}
}
