package biz

import (
	v1 "algo-agent/api/deploy/v1"
	"algo-agent/internal/cons/file"
	taskArgs "algo-agent/internal/cons/task"
	"algo-agent/internal/utils"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	di "algo-agent/internal/model/deploy"
	ds "algo-agent/internal/model/deploy"

	"github.com/go-kratos/kratos/v2/log"
)

type DeployUsecase struct {
	dsm DeployServiceManager

	mq  MqService
	d   DockerService
	oss OSSService

	log *log.Helper

	filePath        string
	mappingFilePath string

	isn string // 推理脚本名称
	dsn string // 部署服务名称
}

func (duc *DeployUsecase) editServiceAndSendMq(ctx context.Context, serivce *di.DeployServiceInfo) error {
	duc.dsm.UpdateService(ctx, serivce)
	return duc.sendStatusChangeMessage(ctx, serivce)
}

func (duc *DeployUsecase) sendStatusChangeMessage(ctx context.Context, serivce *di.DeployServiceInfo) error {
	reply := &v1.DeployReply{
		ServiceId:     serivce.ServiceId,
		ServiceStatus: serivce.ServiceStatus,
		Remark:        serivce.Remark,
	}

	jsonStr, err := utils.ToJSON(reply)
	if err != nil {
		return fmt.Errorf("转换为JSON失败: %v", err)
	}

	return duc.mq.SendToService(ctx, duc.dsn, jsonStr)
}

// 部署一个推理服务
func (duc *DeployUsecase) Deploy(ctx context.Context, serviceInfo *di.DeployServiceInfo) error {
	duc.log.Infof("开始部署服务，服务ID: %s", serviceInfo.ServiceId)

	// 验证服务ID
	if serviceInfo.ServiceId == "" {
		duc.log.Error("服务ID验证失败：值为空")
		return errors.New("服务ID不允许为空")
	}

	deployRequest := serviceInfo.DeployRequest
	if deployRequest == nil {
		return errors.New("部署请求数据不能为空")
	}

	// 设置服务状态为"正在部署"
	duc.log.Info("设置服务状态为'正在部署'")
	serviceInfo.ServiceStatus = ds.DEPLOYING.Code // DEPLOYING 状态码
	if err := duc.dsm.AddService(ctx, serviceInfo); err != nil {
		return fmt.Errorf("添加服务失败: %v", err)
	}

	// 发送MQ消息
	if err := duc.sendStatusChangeMessage(ctx, serviceInfo); err != nil {
		return fmt.Errorf("发送状态变更消息失败: %v", err)
	}

	// 创建必要的目录
	deployBasePath := filepath.Join(duc.filePath, file.DEPLOY)
	imagePath := filepath.Join(duc.filePath, file.IMAGE)
	utils.EnsureDirectoryExists(deployBasePath)
	utils.EnsureDirectoryExists(imagePath)

	// 检查Docker镜像是否存在
	duc.log.Infof("检查Docker镜像是否存在: %s", deployRequest.InferImageName)
	imageExists, err := duc.d.FindImageByName(ctx, deployRequest.InferImageName)
	if err != nil || !imageExists {
		duc.log.Info("本地未找到Docker镜像，准备下载")
		serviceInfo.ServiceStatus = ds.DOWNLOAD_DEPLOY_IMAGE.Code // DOWNLOAD_DEPLOY_IMAGE 状态码
		if err := duc.editServiceAndSendMq(ctx, serviceInfo); err != nil {
			return fmt.Errorf("更新服务状态失败: %v", err)
		}

		// 下载镜像文件
		imageFileName := strings.Replace(deployRequest.InferImageName, ":", "-", -1)
		tarPath := filepath.Join(imagePath, imageFileName+".tar")
		err = duc.oss.DownloadSingleFile(
			ctx,
			deployRequest.InferImageBucket,
			deployRequest.InferImagePath,
			imagePath,
			imageFileName+".tar",
		)
		if err != nil {
			return fmt.Errorf("下载镜像文件失败: %v", err)
		}

		// 导入镜像
		err = duc.d.ImportAndTagImage(ctx, tarPath, deployRequest.InferImageName)
		if err != nil {
			return fmt.Errorf("导入镜像失败: %v", err)
		}
	} else {
		duc.log.Info("找到已存在的Docker镜像")
	}

	// 设置任务路径
	taskPath := filepath.Join(deployBasePath, serviceInfo.ServiceId)

	// 下载算法脚本
	duc.log.Info("开始下载算法脚本")
	serviceInfo.ServiceStatus = ds.DOWNLOAD_ALGO_SCRIPTS.Code // DOWNLOAD_ALGO_SCRIPTS 状态码
	if err := duc.editServiceAndSendMq(ctx, serviceInfo); err != nil {
		return fmt.Errorf("更新服务状态失败: %v", err)
	}

	scriptZipFile := "script_" + serviceInfo.ServiceId + ".zip"
	// 下载算法脚本
	err = duc.oss.DownloadSingleFile(
		ctx,
		deployRequest.AlgorithmScriptBucket,
		deployRequest.AlgorithmScriptPath,
		taskPath,
		scriptZipFile,
	)
	if err != nil {
		return fmt.Errorf("下载算法脚本失败: %v", err)
	}

	// 解压脚本文件
	duc.log.Info("解压算法脚本")
	scriptZipPath := filepath.Join(taskPath, scriptZipFile)
	scriptDestPath := filepath.Join(taskPath, file.SCRIPT)
	if err := utils.Unzip(scriptZipPath, scriptDestPath); err != nil {
		return fmt.Errorf("解压算法脚本失败: %v", err)
	}

	// 下载推理权重文件
	duc.log.Info("开始下载推理权重文件")
	serviceInfo.ServiceStatus = ds.DOWNLOAD_WEIGHTS.Code // DOWNLOAD_WEIGHTS 状态码
	if err := duc.editServiceAndSendMq(ctx, serviceInfo); err != nil {
		return fmt.Errorf("更新服务状态失败: %v", err)
	}

	modelDir := filepath.Join(taskPath, file.MODEL)
	utils.EnsureDirectoryExists(modelDir)

	if deployRequest.ModelBucket != "" && deployRequest.ModelPath != "" {
		duc.log.Info("从MinIO下载模型文件")
		err = duc.oss.DownloadSingleFile(
			ctx,
			deployRequest.ModelBucket,
			deployRequest.ModelPath,
			modelDir,
			filepath.Base(deployRequest.ModelPath),
		)
		if err != nil {
			return fmt.Errorf("下载模型文件失败: %v", err)
		}
	}

	modelFileName := filepath.Base(deployRequest.ModelPath)

	// 设置容器运行参数
	hostPath, _ := filepath.Abs(taskPath)
	containerPath := file.UNIX_SEPARATOR + file.WORKSPACE
	scriptPath := containerPath + file.UNIX_SEPARATOR + file.SCRIPT + file.UNIX_SEPARATOR + duc.isn

	duc.log.Info("设置容器参数")
	args := []string{
		taskArgs.ArgServiceID, serviceInfo.ServiceId,
		taskArgs.ArgModelPath, containerPath + file.UNIX_SEPARATOR + file.MODEL + file.UNIX_SEPARATOR + modelFileName,
	}

	// 添加额外参数
	if len(deployRequest.Args) > 0 {
		args = append(args, deployRequest.Args...)
	}

	// 添加数据集标签
	if deployRequest.DatasetLabel != "" {
		addLabels(deployRequest.DatasetLabel, &args)
	}

	// 启动容器
	duc.log.Infof("使用镜像启动容器: %s", deployRequest.InferImageName)
	containerInfo, err := duc.d.RunAndStartContainer(
		ctx,
		deployRequest.InferImageName,
		hostPath,
		containerPath,
		scriptPath,
		args,
	)
	if err != nil {
		return fmt.Errorf("启动容器失败: %v", err)
	}

	// 更新服务状态为"运行中"
	duc.log.Infof("容器成功启动，更新服务状态为'运行中'，容器名称: %s", containerInfo.ContainerName)
	serviceInfo.ServiceStatus = 5 // RUNNING 状态码
	serviceInfo.ServiceContainerName = containerInfo.ContainerName
	return duc.editServiceAndSendMq(ctx, serviceInfo)
}

func addLabels(datasetLabel string, args *[]string) {
	// 解析JSON字符串为map
	stringObjectMap, err := utils.ParseToMap(datasetLabel)
	if err != nil {
		// 如果解析失败，直接返回
		return
	}

	if stringObjectMap != nil && len(stringObjectMap) > 0 {
		// 创建一个可排序的切片
		type KeyValue struct {
			Key   int
			Value interface{}
		}

		sortedList := make([]KeyValue, 0, len(stringObjectMap))

		// 将map转换为切片
		for k, v := range stringObjectMap {
			key, err := strconv.Atoi(k)
			if err != nil {
				continue
			}
			sortedList = append(sortedList, KeyValue{Key: key, Value: v})
		}

		// 按key排序
		sort.Slice(sortedList, func(i, j int) bool {
			return sortedList[i].Key < sortedList[j].Key
		})

		// 创建一个字符串切片存储值
		values := make([]string, 0, len(sortedList))
		for _, entry := range sortedList {
			values = append(values, fmt.Sprintf("%v", entry.Value))
		}

		// 添加类名参数
		*args = append(*args, taskArgs.ArgClassNames)
		*args = append(*args, strings.Join(values, ","))
	}
}

// DestroyAndDelete 销毁容器，并删除相关文件
func (duc *DeployUsecase) DestroyAndDelete(ctx context.Context, serviceId string) error {
	duc.log.Infof("开始销毁和删除服务，服务ID: %s", serviceId)

	// 查找服务信息
	service := duc.dsm.FindServiceById(ctx, serviceId)
	if service == nil {
		duc.log.Warnf("销毁失败，未找到服务。serviceId: %s", serviceId)
		return fmt.Errorf("销毁失败，未找到服务。serviceId: %s", serviceId)
	}
	duc.log.Infof("找到服务信息: %+v", service)

	// 停止容器
	containerName := service.ServiceContainerName
	if containerName == "" {
		duc.log.Warnf("销毁失败，服务容器名为空。serviceId: %s", serviceId)
	} else {
		duc.log.Infof("停止并移除容器: %s", containerName)
		err := duc.d.StopContainerByName(ctx, containerName, true)
		if err != nil {
			duc.log.Errorf("停止容器失败: %v", err)
		} else {
			duc.log.Info("容器已成功停止并移除")
		}
	}

	// 删除目录
	storePath := filepath.Join(duc.filePath, file.DEPLOY)
	taskPath := filepath.Join(storePath, serviceId)
	duc.log.Infof("删除服务目录: %s", taskPath)
	err := utils.RemoveDirectory(taskPath)
	if err != nil {
		duc.log.Errorf("删除推理服务目录失败。serviceId: %s, 错误: %v", serviceId, err)
	} else {
		duc.log.Infof("删除推理服务目录成功。serviceId: %s", serviceId)
	}

	// 删除服务记录
	duc.log.Infof("从管理器中移除服务, serviceId: %s", serviceId)
	duc.dsm.RemoveService(ctx, serviceId)
	duc.log.Infof("服务已成功销毁和删除, serviceId: %s", serviceId)
	return nil
}

// CheckTask 检查所有正在运行的任务，监控容器状态
func (duc *DeployUsecase) CheckTask(ctx context.Context) {
	duc.log.Info("检查任务...")
	services := duc.dsm.GetServiceList(ctx)
	if len(services) == 0 {
		duc.log.Info("检查任务, 没有服务")
		return
	}

	for _, service := range services {
		containerName := service.ServiceContainerName
		if containerName == "" {
			duc.log.Infof("检查任务, 容器名为空. serviceId: %s", service.ServiceId)
			continue
		}
		if service.ServiceStatus != ds.RUNNING.Code {
			duc.log.Infof("检查任务, 服务未在运行中, 状态: %d, serviceId: %s", service.ServiceStatus, service.ServiceId)
			continue
		}

		containerInfo, err := duc.d.FindContainerByName(ctx, containerName)
		if err != nil || containerInfo == nil {
			duc.log.Infof("检查任务, 未找到容器. serviceId: %s", service.ServiceId)
			continue
		}

		// 获取容器状态
		inspect, err := duc.d.GetContainerState(ctx, containerInfo.ContainerID)
		if err != nil {
			duc.log.Errorf("获取容器状态失败: %v", err)
			continue
		}

		duc.log.Infof("检查任务, serviceId: %s, 容器状态: %s", service.ServiceId, inspect.State)

		// 检查容器是否在运行中
		if inspect.State != "running" {
			duc.log.Infof("容器状态异常，终止任务: %s", containerName)

			// 获取容器最后的日志
			logs, err := duc.d.GetContainerLastLogs(ctx, containerInfo.ContainerID, 10)
			if err != nil {
				duc.log.Errorf("获取容器日志失败: %v", err)
			} else {
				duc.log.Infof("错误日志: %s", logs)
			}

			// 更新服务状态为部署失败
			service.ServiceStatus = ds.DEPLOYMENT_FAILED.Code

			// 从日志中提取最后一行作为错误消息
			lastLine := ""
			if logs != "" {
				lines := strings.Split(logs, "\n")
				for i := len(lines) - 1; i >= 0; i-- {
					if lines[i] != "" {
						lastLine = lines[i]
						break
					}
				}
			}

			service.Remark = "任务执行错误: " + lastLine
			err = duc.editServiceAndSendMq(ctx, service)
			if err != nil {
				duc.log.Errorf("更新服务状态失败: %v", err)
			}

			// 销毁和删除服务
			duc.DestroyAndDelete(ctx, service.ServiceId)
		}
	}
}

// HandleEvent 处理部署消息
func (duc *DeployUsecase) HandleEvent(ctx context.Context, deployMessage *v1.DeployRequest) {
	serviceId := deployMessage.GetServiceId()

	// 创建服务信息
	serviceInfo := di.NewDeployServiceInfo(deployMessage)

	// 处理部署操作
	if deployMessage.GetOp() == ds.DEPLOY.Code {
		duc.log.Infof("收到部署操作请求，serviceId: %s", serviceId)
		err := duc.Deploy(ctx, serviceInfo)
		if err != nil {
			duc.log.Errorf("启动推理服务失败! serviceId: %s, 错误: %v", serviceId, err)
			serviceInfo.ServiceStatus = ds.DEPLOYMENT_FAILED.Code
			serviceInfo.Remark = "启动错误: " + err.Error()
			duc.editServiceAndSendMq(ctx, serviceInfo)

			duc.DestroyAndDelete(ctx, serviceId)
		}
	} else if deployMessage.GetOp() == ds.DESTROY.Code {
		duc.log.Infof("收到销毁操作请求，serviceId: %s", serviceId)

		service := duc.dsm.FindServiceById(ctx, serviceId)
		if service == nil {
			duc.log.Warnf("销毁失败，未找到服务。serviceId: %s", serviceId)
			emptyService := &di.DeployServiceInfo{
				ServiceId:     serviceId,
				ServiceStatus: ds.DESTROYED.Code,
				Remark:        "未找到服务!",
			}
			duc.sendStatusChangeMessage(ctx, emptyService)
			return
		}

		duc.DestroyAndDelete(ctx, serviceId)
		service.ServiceStatus = ds.DESTROYED.Code
		duc.sendStatusChangeMessage(ctx, service)
	} else {
		duc.log.Errorf("未知的任务操作! op: %d, serviceId: %s", deployMessage.GetOp(), serviceId)
	}
}
