package biz

import (
	"algo-agent/internal/model/extract"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	cc "algo-agent/internal/cons/container"
	"algo-agent/internal/cons/file"
	taskArgs "algo-agent/internal/cons/task"
	ct "algo-agent/internal/cons/train"
	"algo-agent/internal/mq/event"
	"algo-agent/internal/utils"
)

// ExtractTaskUsecase 提取任务用例
type ExtractTaskUsecase struct {
	etm               ExtractTaskManager // 提取任务管理器
	mq                MqService          // 消息队列服务
	d                 DockerService      // Docker服务
	oss               OSSService         // 对象存储服务
	log               *log.Helper        // 日志
	filePath          string             // 文件路径
	tsn               string             // 训练服务名称
	extractScriptName string             // 提取脚本名称
}

// 更新状态并发送消息
func (uc *ExtractTaskUsecase) editTaskAndSendMq(ctx context.Context, task *extract.ExtractTaskInfo) error {
	uc.etm.UpdateTask(ctx, task)
	return uc.sendStatusChangeMessage(ctx, task)
}

// 发送状态变更消息
func (uc *ExtractTaskUsecase) sendStatusChangeMessage(ctx context.Context, taskInfo *extract.ExtractTaskInfo) error {
	reply := &event.TrainPublishRespMessage{
		TrainDetailId:  taskInfo.TaskId,
		ModelVersionId: taskInfo.TrainPublishReqMessage.ModelVersionId,
		TaskStatus:     taskInfo.TaskStatus,
		Remark:         taskInfo.Remark,
	}

	jsonStr, err := utils.ToJSON(reply)
	if err != nil {
		return fmt.Errorf("转换为JSON失败: %v", err)
	}

	return uc.mq.SendToService(ctx, uc.tsn, jsonStr)
}

// AddTask 添加提取任务
func (uc *ExtractTaskUsecase) AddTask(ctx context.Context, task *extract.ExtractTaskInfo) error {
	return uc.etm.AddTask(ctx, task)
}

// UpdateTask 更新提取任务
func (uc *ExtractTaskUsecase) UpdateTask(ctx context.Context, task *extract.ExtractTaskInfo) error {
	return uc.etm.UpdateTask(ctx, task)
}

// RemoveTask 移除提取任务
func (uc *ExtractTaskUsecase) RemoveTask(ctx context.Context, id string) bool {
	return uc.etm.RemoveTask(ctx, id)
}

// GetTaskList 获取提取任务列表
func (uc *ExtractTaskUsecase) GetTaskList(ctx context.Context) []*extract.ExtractTaskInfo {
	return uc.etm.GetTaskList(ctx)
}

// FindTaskById 根据ID查找提取任务
func (uc *ExtractTaskUsecase) FindTaskById(ctx context.Context, id string) *extract.ExtractTaskInfo {
	return uc.etm.FindTaskById(ctx, id)
}

// StartExtract 开始提取任务
func (uc *ExtractTaskUsecase) StartExtract(ctx context.Context, taskInfo *extract.ExtractTaskInfo) error {
	eventInfo := taskInfo.TrainPublishReqMessage
	if eventInfo == nil || eventInfo.TrainDetailId == "" {
		uc.log.Error("开始提取失败：任务ID为空")
		return errors.New("任务ID不允许为空")
	}
	uc.log.Infof("开始提取任务，任务ID: %s", eventInfo.TrainDetailId)

	// 设置任务状态为"启动中"
	taskInfo.TaskStatus = ct.TrainJobStatusStarting.Code
	if err := uc.etm.AddTask(ctx, taskInfo); err != nil {
		return fmt.Errorf("添加任务失败: %v", err)
	}

	// 创建必要的目录
	extractBasePath := filepath.Join(uc.filePath, file.EXTRACT)
	imagePath := filepath.Join(uc.filePath, file.IMAGE)
	utils.EnsureDirectoryExists(extractBasePath)
	utils.EnsureDirectoryExists(imagePath)

	// 检查Docker镜像是否存在
	uc.log.Infof("检查提取镜像是否存在: %s", eventInfo.AlgorithmTrainImageName)
	imageExists, err := uc.d.FindImageByName(ctx, eventInfo.AlgorithmTrainImageName)
	if err != nil || !imageExists {
		uc.log.Info("本地未找到提取镜像，准备下载")
		taskInfo.TaskStatus = ct.TrainJobStatusDownloadTrainImage.Code
		if err := uc.etm.UpdateTask(ctx, taskInfo); err != nil {
			return fmt.Errorf("更新任务状态失败: %v", err)
		}

		// 下载镜像文件
		imageFileName := strings.Replace(eventInfo.AlgorithmTrainImageName, ":", "-", -1)
		tarPath := filepath.Join(imagePath, imageFileName+".tar")
		err = uc.oss.DownloadSingleFile(
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
		err = uc.d.ImportAndTagImage(ctx, tarPath, eventInfo.AlgorithmTrainImageName)
		if err != nil {
			return fmt.Errorf("导入镜像失败: %v", err)
		}
	} else {
		uc.log.Info("找到已存在的提取镜像")
	}

	// 设置任务路径
	taskPath := filepath.Join(extractBasePath, eventInfo.TrainDetailId)
	utils.EnsureDirectoryExists(taskPath)

	// 下载算法脚本
	uc.log.Info("开始下载算法脚本")
	taskInfo.TaskStatus = ct.TrainJobStatusDownloadAlgoScripts.Code
	if err := uc.etm.UpdateTask(ctx, taskInfo); err != nil {
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	algorithmScriptZipFile := "script_" + eventInfo.TrainDetailId + ".zip"
	// 下载算法脚本
	err = uc.oss.DownloadSingleFile(
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
	uc.log.Info("解压算法脚本")
	scriptZipPath := filepath.Join(taskPath, algorithmScriptZipFile)
	scriptDestPath := filepath.Join(taskPath, file.SCRIPT)
	if err := utils.Unzip(scriptZipPath, scriptDestPath); err != nil {
		return fmt.Errorf("解压算法脚本失败: %v", err)
	}

	// 下载检查点文件
	uc.log.Infof("下载检查点, 任务ID: %s", eventInfo.TrainDetailId)
	taskInfo.TaskStatus = ct.TrainJobStatusDownloadCheckpointFile.Code
	if err := uc.etm.UpdateTask(ctx, taskInfo); err != nil {
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	checkpointPath := filepath.Join(taskPath, file.CHECKPOINT)
	utils.EnsureDirectoryExists(checkpointPath)

	const checkpointName = "checkpoint.pth"
	err = uc.oss.DownloadSingleFile(
		ctx,
		eventInfo.CheckpointBucket,
		eventInfo.CheckpointFileURL,
		checkpointPath,
		checkpointName,
	)
	if err != nil {
		return fmt.Errorf("下载检查点文件失败: %v", err)
	}

	// 设置容器运行参数
	hostPath, _ := filepath.Abs(taskPath)
	containerPath := file.UNIX_SEPARATOR + file.WORKSPACE
	scriptPath := containerPath + file.UNIX_SEPARATOR + file.SCRIPT + file.UNIX_SEPARATOR + uc.extractScriptName

	uc.log.Info("设置容器参数")
	args := []string{
		taskArgs.ArgTaskID, eventInfo.TrainDetailId,
	}

	// 添加提取任务参数
	if eventInfo.Args != nil && len(eventInfo.Args) > 0 {
		args = append(args, eventInfo.Args...)
	}

	// 添加检查点路径参数
	args = append(args,
		taskArgs.ArgCheckpointPath,
		containerPath+file.UNIX_SEPARATOR+file.CHECKPOINT+file.UNIX_SEPARATOR+checkpointName,
	)

	// 启动容器
	uc.log.Infof("使用镜像启动容器: %s", eventInfo.AlgorithmTrainImageName)
	containerInfo, err := uc.d.RunAndStartContainer(
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
	uc.d.StartLogStream(ctx, containerInfo.ContainerId, func(logText string) {
		uc.sendDockerLogData(ctx, eventInfo.TrainDetailId, logText)
	})

	// 更新任务状态为"运行中"
	uc.log.Infof("容器成功启动，更新任务状态为'运行中'，容器名称: %s", containerInfo.ContainerName)
	taskInfo.TaskStatus = ct.TrainJobStatusRunning.Code
	taskInfo.ContainerName = containerInfo.ContainerName
	return uc.etm.UpdateTask(ctx, taskInfo)
}

// 发送日志数据
func (uc *ExtractTaskUsecase) sendDockerLogData(ctx context.Context, taskId string, logText string) {
	taskIdInt, err := strconv.ParseInt(taskId, 10, 64)
	if err != nil {
		uc.log.Errorf("转换任务ID失败: %v", err)
		return
	}

	logMsg := &event.DockerLogRespMessage{
		TaskId:   taskIdInt,
		Log:      logText,
		TaskType: 0, // 提取任务也使用类型0
	}

	jsonStr, err := utils.ToJSON(logMsg)
	if err != nil {
		uc.log.Errorf("转换日志消息为JSON失败: %v", err)
		return
	}

	if err := uc.mq.SendToService(ctx, uc.tsn, jsonStr); err != nil {
		uc.log.Errorf("发送日志消息失败: %v", err)
	}
}

// JustStop 只停止容器，不删除容器和相关文件
func (uc *ExtractTaskUsecase) JustStop(ctx context.Context, taskId string, remove bool) {
	uc.log.Infof("开始停止任务，任务ID: %s", taskId)
	task := uc.etm.FindTaskById(ctx, taskId)
	if task == nil {
		uc.log.Warnf("停止失败，未找到任务。taskId: %s", taskId)
		return
	}
	uc.log.Infof("找到任务信息: taskId=%s, containerName=%s", taskId, task.ContainerName)

	// 停止日志
	containerName := task.ContainerName
	if containerName != "" {
		uc.d.StopLogStream(ctx, containerName, func(logText string) {
			uc.sendDockerLogData(ctx, taskId, logText)
		})

		// 停止容器
		uc.log.Infof("停止容器: %s, remove=%v", containerName, remove)
		err := uc.d.StopContainerByName(ctx, containerName, remove)
		if err != nil {
			uc.log.Errorf("停止容器失败: %v", err)
		} else {
			uc.log.Infof("容器已成功停止: %s", containerName)
		}
	}
}

// CleanupTaskDirectoryAndRecord 清理任务相关的目录和记录
func (uc *ExtractTaskUsecase) CleanupTaskDirectoryAndRecord(ctx context.Context, taskId string) {
	// 删除目录
	storePath := filepath.Join(uc.filePath, file.EXTRACT)
	taskPath := filepath.Join(storePath, taskId)
	uc.log.Infof("删除任务目录: %s", taskPath)
	err := utils.RemoveDirectory(taskPath)
	if err != nil {
		uc.log.Errorf("删除提取任务目录失败。taskId: %s, 错误: %v", taskId, err)
	} else {
		uc.log.Infof("删除提取任务目录成功。taskId: %s", taskId)
	}

	// 删除记录
	uc.log.Infof("从管理器中移除任务, taskId: %s", taskId)
	uc.etm.RemoveTask(ctx, taskId)
	uc.log.Infof("任务已成功从管理器中移除, taskId: %s", taskId)
}

// DestroyAndDelete 销毁容器，并删除相关文件
func (uc *ExtractTaskUsecase) DestroyAndDelete(ctx context.Context, taskId string) {
	uc.log.Infof("开始销毁和删除任务，任务ID: %s", taskId)

	// 停止容器
	uc.JustStop(ctx, taskId, true)
	uc.log.Infof("容器已成功停止并移除, taskId: %s", taskId)

	// 清理任务相关的目录和记录
	uc.CleanupTaskDirectoryAndRecord(ctx, taskId)
	uc.log.Infof("任务已成功销毁和删除, taskId: %s", taskId)
}

// ExtractTaskResultHandle 处理提取任务结果
func (uc *ExtractTaskUsecase) ExtractTaskResultHandle(ctx context.Context, result *extract.ExtractTaskResult) error {
	uc.log.Infof("处理提取任务结果, taskId: %s", result.TaskId)
	taskInfo := uc.etm.FindTaskById(ctx, result.TaskId)
	if taskInfo == nil {
		return fmt.Errorf("处理提取任务结果失败，未找到任务。taskId: %s", result.TaskId)
	}

	modelPath := result.ModelPath
	// 检查路径是否以 /workspace 开头，如果是，则需要去掉
	prefix := file.UNIX_SEPARATOR + file.WORKSPACE
	if strings.HasPrefix(modelPath, prefix) {
		uc.log.Infof("模型路径以 %s 开头。taskId: %s, modelPath: %s", prefix, result.TaskId, modelPath)
		modelPath = modelPath[len(prefix):]
		uc.log.Infof("处理后的模型路径: %s", modelPath)
	}

	taskBasePath := filepath.Join(uc.filePath, file.EXTRACT, result.TaskId)
	modelFilePath := filepath.Join(taskBasePath, modelPath)
	fileName := filepath.Base(modelFilePath)

	// 模型MinIO路径前缀
	modelPathPrefix := "Model/train_detail-"
	objectUrl := modelPathPrefix + result.TaskId + file.UNIX_SEPARATOR + fileName

	// 上传模型文件
	err := uc.oss.UploadFile(
		ctx,
		taskInfo.TrainPublishReqMessage.AlgorithmScriptMinioBucket,
		&File{
			Name: fileName,
			Path: modelFilePath,
		},
		objectUrl,
	)
	if err != nil {
		return fmt.Errorf("上传模型文件失败: %v", err)
	}

	// 发送任务完成消息
	reply := &event.TrainPublishRespMessage{
		TrainDetailId:    result.TaskId,
		ModelVersionId:   taskInfo.TrainPublishReqMessage.ModelVersionId,
		TaskStatus:       ct.TrainJobStatusSucceed.Code,
		ModelWeightsPath: objectUrl,
	}

	jsonStr, err := utils.ToJSON(reply)
	if err != nil {
		return fmt.Errorf("转换为JSON失败: %v", err)
	}

	if err := uc.mq.SendToService(ctx, uc.tsn, jsonStr); err != nil {
		uc.log.Errorf("发送任务完成消息失败: %v", err)
		return err
	}

	// 销毁任务
	uc.DestroyAndDelete(ctx, result.TaskId)
	return nil
}

// HandleExtractTask 处理提取任务消息
func (uc *ExtractTaskUsecase) HandleExtractTask(ctx context.Context, task *event.TrainPublishReqMessage) {
	taskId := task.TrainDetailId
	taskInfo := extract.NewExtractTaskInfoWithMessage(task)

	uc.log.Infof("收到提取任务请求，任务ID: %s", taskId)
	err := uc.StartExtract(ctx, taskInfo)
	if err != nil {
		uc.log.Errorf("启动提取任务失败! taskId: %s, 错误: %v", taskId, err)
		taskInfo.TaskStatus = ct.TrainJobStatusFail.Code
		taskInfo.Remark = "启动错误: " + err.Error()
		uc.editTaskAndSendMq(ctx, taskInfo)
		uc.JustStop(ctx, taskId, false)
	}
}

// CheckTask 检查所有任务的状态
func (uc *ExtractTaskUsecase) CheckTask(ctx context.Context) {
	uc.log.Debug("开始检查任务状态...")
	tasks := uc.etm.GetTaskList(ctx)
	if len(tasks) == 0 {
		uc.log.Debug("没有需要检查的任务")
		return
	}

	for _, task := range tasks {
		// 使用匿名函数避免循环中的错误影响到所有任务的检查
		func(task *extract.ExtractTaskInfo) {
			defer func() {
				if r := recover(); r != nil {
					uc.log.Errorf("检查任务时发生异常: %v", r)
				}
			}()

			if task.ContainerName == "" {
				uc.log.Infof("任务 %s 没有关联的容器名称，跳过检查", task.TaskId)
				return
			}

			if task.TaskStatus != ct.TrainJobStatusRunning.Code {
				uc.log.Infof("任务 %s 不是运行中状态，跳过检查", task.TaskId)
				return
			}

			containerInfo, err := uc.d.FindContainerByName(ctx, task.ContainerName)
			if err != nil {
				uc.log.Errorf("查找容器信息失败: %v", err)
				return
			}

			if containerInfo == nil {
				uc.log.Infof("任务 %s 的容器不存在", task.TaskId)
				return
			}

			containerState, err := uc.d.GetContainerState(ctx, containerInfo.ContainerId)
			if err != nil {
				uc.log.Errorf("获取容器状态失败: %v", err)
				return
			}

			if containerState.State != string(cc.RUNNING) {
				// 获取容器停止时间
				stopTime, err := uc.d.GetContainerStopTime(ctx, containerInfo.ContainerId)
				if err != nil {
					uc.log.Errorf("获取容器停止时间失败: %v", err)
					return
				}

				currentTime := time.Now().Unix()
				timeDifferenceMinutes := (currentTime - stopTime) / 60

				if timeDifferenceMinutes < 10 {
					uc.log.Infof("容器停止不超过10分钟，跳过处理: %s", task.ContainerName)
					return
				}

				uc.log.Errorf("任务 %s 标记为运行中但容器已停止超过10分钟", task.TaskId)
				logs, err := uc.d.GetContainerLastLogs(ctx, containerInfo.ContainerId, 10)
				if err != nil {
					uc.log.Errorf("获取容器日志失败: %v", err)
				}

				task.TaskStatus = ct.TrainJobStatusFail.Code
				task.Remark = "容器意外停止: " + logs
				uc.editTaskAndSendMq(ctx, task)
				uc.DestroyAndDelete(ctx, task.TaskId)
			}
		}(task)
	}
	uc.log.Info("任务状态检查完成")
}
