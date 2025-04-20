package biz

import (
	"context"

	mc "algo-agent/internal/model/container"

	"github.com/go-kratos/kratos/v2/log"
)

// DockerService Docker服务接口
type DockerService interface {
	// FindContainerByName 根据容器名查找容器
	FindContainerByName(ctx context.Context, containerName string) (*mc.ContainerInfo, error)

	// RunContainer 通过镜像和自定义参数运行一个容器
	RunContainer(ctx context.Context, imageName string, customArgs []string) (string, error)

	// RunAndStartContainer 运行并启动容器，返回容器信息
	RunAndStartContainer(ctx context.Context, imageName, hostPath, containerPath, scriptPath string, args []string) (*mc.ContainerInfo, error)

	// RunAndStartContainerWithPort 运行并启动带端口映射的容器
	RunAndStartContainerWithPort(ctx context.Context, imageName, hostPath, containerPath, command string, args []string, hostPort int) (string, error)

	// GetContainerLastLogs 获取指定容器的最后若干行日志
	GetContainerLastLogs(ctx context.Context, containerID string, tail int) (string, error)

	// StartLogStream 开始日志流
	StartLogStream(ctx context.Context, containerID string, logConsumer func(string))

	// StopLogStream 停止日志流
	StopLogStream(ctx context.Context, containerName string, logConsumer func(string))

	// StopContainer 停止容器
	StopContainer(ctx context.Context, containerID string) error

	// StopContainerByName 通过名称停止容器
	StopContainerByName(ctx context.Context, containerName string, remove bool) error

	// GetContainerStopTime 获取容器的停止时间戳
	GetContainerStopTime(ctx context.Context, containerID string) (int64, error)

	// FindImageByName 查找指定名称的镜像
	FindImageByName(ctx context.Context, imageName string) (bool, error)

	// ImportAndTagImage 导入并标记镜像
	ImportAndTagImage(ctx context.Context, tarFilePath, fullImageName string) error

	// Close 关闭连接
	Close()
}

// DockerUsecase 是Docker用例
type DockerUsecase struct {
	docker DockerService
	log    *log.Helper
}

// FindContainerByName 根据容器名查找容器
func (uc *DockerUsecase) FindContainerByName(ctx context.Context, containerName string) (*mc.ContainerInfo, error) {
	uc.log.WithContext(ctx).Infof("FindContainerByName: containerName=%s", containerName)
	return uc.docker.FindContainerByName(ctx, containerName)
}

// RunContainer 通过镜像和自定义参数运行一个容器
func (uc *DockerUsecase) RunContainer(ctx context.Context, imageName string, customArgs []string) (string, error) {
	uc.log.WithContext(ctx).Infof("RunContainer: imageName=%s", imageName)
	return uc.docker.RunContainer(ctx, imageName, customArgs)
}

// RunAndStartContainer 运行并启动容器，返回容器信息
func (uc *DockerUsecase) RunAndStartContainer(ctx context.Context, imageName, hostPath, containerPath, scriptPath string, args []string) (*mc.ContainerInfo, error) {
	uc.log.WithContext(ctx).Infof("RunAndStartContainer: imageName=%s, hostPath=%s, containerPath=%s", imageName, hostPath, containerPath)
	return uc.docker.RunAndStartContainer(ctx, imageName, hostPath, containerPath, scriptPath, args)
}

// RunAndStartContainerWithPort 运行并启动带端口映射的容器
func (uc *DockerUsecase) RunAndStartContainerWithPort(ctx context.Context, imageName, hostPath, containerPath, command string, args []string, hostPort int) (string, error) {
	uc.log.WithContext(ctx).Infof("RunAndStartContainerWithPort: imageName=%s, hostPath=%s, containerPath=%s, hostPort=%d", imageName, hostPath, containerPath, hostPort)
	return uc.docker.RunAndStartContainerWithPort(ctx, imageName, hostPath, containerPath, command, args, hostPort)
}

// GetContainerLastLogs 获取指定容器的最后若干行日志
func (uc *DockerUsecase) GetContainerLastLogs(ctx context.Context, containerID string, tail int) (string, error) {
	uc.log.WithContext(ctx).Infof("GetContainerLastLogs: containerID=%s, tail=%d", containerID, tail)
	return uc.docker.GetContainerLastLogs(ctx, containerID, tail)
}

// StartLogStream 开始日志流
func (uc *DockerUsecase) StartLogStream(ctx context.Context, containerID string, logConsumer func(string)) {
	uc.log.WithContext(ctx).Infof("StartLogStream: containerID=%s", containerID)
	uc.docker.StartLogStream(ctx, containerID, logConsumer)
}

// StopLogStream 停止日志流
func (uc *DockerUsecase) StopLogStream(ctx context.Context, containerName string, logConsumer func(string)) {
	uc.log.WithContext(ctx).Infof("StopLogStream: containerName=%s", containerName)
	uc.docker.StopLogStream(ctx, containerName, logConsumer)
}

// StopContainer 停止容器
func (uc *DockerUsecase) StopContainer(ctx context.Context, containerID string) error {
	uc.log.WithContext(ctx).Infof("StopContainer: containerID=%s", containerID)
	return uc.docker.StopContainer(ctx, containerID)
}

// StopContainerByName 通过名称停止容器
func (uc *DockerUsecase) StopContainerByName(ctx context.Context, containerName string, remove bool) error {
	uc.log.WithContext(ctx).Infof("StopContainerByName: containerName=%s, remove=%v", containerName, remove)
	return uc.docker.StopContainerByName(ctx, containerName, remove)
}

// GetContainerStopTime 获取容器的停止时间戳
func (uc *DockerUsecase) GetContainerStopTime(ctx context.Context, containerID string) (int64, error) {
	uc.log.WithContext(ctx).Infof("GetContainerStopTime: containerID=%s", containerID)
	return uc.docker.GetContainerStopTime(ctx, containerID)
}

// FindImageByName 查找指定名称的镜像
func (uc *DockerUsecase) FindImageByName(ctx context.Context, imageName string) (bool, error) {
	uc.log.WithContext(ctx).Infof("FindImageByName: imageName=%s", imageName)
	return uc.docker.FindImageByName(ctx, imageName)
}

// ImportAndTagImage 导入并标记镜像
func (uc *DockerUsecase) ImportAndTagImage(ctx context.Context, tarFilePath, fullImageName string) error {
	uc.log.WithContext(ctx).Infof("ImportAndTagImage: tarFilePath=%s, fullImageName=%s", tarFilePath, fullImageName)
	return uc.docker.ImportAndTagImage(ctx, tarFilePath, fullImageName)
}
