package service

import (
	"context"

	pb "algo-agent/api/docker/v1"
	"algo-agent/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// DockerServer 实现Docker服务API接口
type DockerServer struct {
	pb.UnimplementedDockerServiceServer
	uc  *biz.DockerUsecase
	log *log.Helper
}

// NewDockerServer 创建Docker服务实例
func NewDockerServer(uc *biz.DockerUsecase, logger log.Logger) *DockerServer {
	return &DockerServer{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// FindContainerByName 根据容器名查找容器
func (s *DockerServer) FindContainerByName(ctx context.Context, req *pb.FindContainerByNameRequest) (*pb.FindContainerByNameReply, error) {
	s.log.WithContext(ctx).Infof("FindContainerByName: containerName=%s", req.ContainerName)

	container, err := s.uc.FindContainerByName(ctx, req.ContainerName)
	if err != nil {
		s.log.WithContext(ctx).Errorf("FindContainerByName failed: %v", err)
		return nil, err
	}

	return &pb.FindContainerByNameReply{
		Container: &pb.ContainerInfo{
			ContainerId:   container.ContainerId,
			ContainerName: container.ContainerName,
		},
	}, nil
}

// RunContainer 通过镜像和自定义参数运行一个容器
func (s *DockerServer) RunContainer(ctx context.Context, req *pb.RunContainerRequest) (*pb.RunContainerReply, error) {
	s.log.WithContext(ctx).Infof("RunContainer: imageName=%s", req.ImageName)

	containerID, err := s.uc.RunContainer(ctx, req.ImageName, req.CustomArgs)
	if err != nil {
		s.log.WithContext(ctx).Errorf("RunContainer failed: %v", err)
		return nil, err
	}

	return &pb.RunContainerReply{
		ContainerId: containerID,
	}, nil
}

// RunAndStartContainer 运行并启动容器，返回容器信息
func (s *DockerServer) RunAndStartContainer(ctx context.Context, req *pb.RunAndStartContainerRequest) (*pb.RunAndStartContainerReply, error) {
	s.log.WithContext(ctx).Infof("RunAndStartContainer: imageName=%s, hostPath=%s, containerPath=%s",
		req.ImageName, req.HostPath, req.ContainerPath)

	container, err := s.uc.RunAndStartContainer(ctx, req.ImageName, req.HostPath, req.ContainerPath, req.ScriptPath, req.Args)
	if err != nil {
		s.log.WithContext(ctx).Errorf("RunAndStartContainer failed: %v", err)
		return nil, err
	}

	return &pb.RunAndStartContainerReply{
		Container: &pb.ContainerInfo{
			ContainerId:   container.ContainerId,
			ContainerName: container.ContainerName,
		},
	}, nil
}

// RunAndStartContainerWithPort 运行并启动带端口映射的容器
func (s *DockerServer) RunAndStartContainerWithPort(ctx context.Context, req *pb.RunAndStartContainerWithPortRequest) (*pb.RunAndStartContainerWithPortReply, error) {
	s.log.WithContext(ctx).Infof("RunAndStartContainerWithPort: imageName=%s, hostPath=%s, containerPath=%s, hostPort=%d",
		req.ImageName, req.HostPath, req.ContainerPath, req.HostPort)

	containerID, err := s.uc.RunAndStartContainerWithPort(ctx, req.ImageName, req.HostPath, req.ContainerPath, req.Command, req.Args, int(req.HostPort))
	if err != nil {
		s.log.WithContext(ctx).Errorf("RunAndStartContainerWithPort failed: %v", err)
		return nil, err
	}

	return &pb.RunAndStartContainerWithPortReply{
		ContainerId: containerID,
	}, nil
}

// GetContainerLastLogs 获取指定容器的最后若干行日志
func (s *DockerServer) GetContainerLastLogs(ctx context.Context, req *pb.GetContainerLastLogsRequest) (*pb.GetContainerLastLogsReply, error) {
	s.log.WithContext(ctx).Infof("GetContainerLastLogs: containerID=%s, tail=%d", req.ContainerId, req.Tail)

	logs, err := s.uc.GetContainerLastLogs(ctx, req.ContainerId, int(req.Tail))
	if err != nil {
		s.log.WithContext(ctx).Errorf("GetContainerLastLogs failed: %v", err)
		return nil, err
	}

	return &pb.GetContainerLastLogsReply{
		Logs: logs,
	}, nil
}

// StopContainer 停止容器
func (s *DockerServer) StopContainer(ctx context.Context, req *pb.StopContainerRequest) (*pb.StopContainerReply, error) {
	s.log.WithContext(ctx).Infof("StopContainer: containerID=%s", req.ContainerId)

	err := s.uc.StopContainer(ctx, req.ContainerId)
	if err != nil {
		s.log.WithContext(ctx).Errorf("StopContainer failed: %v", err)
		return nil, err
	}

	return &pb.StopContainerReply{
		Success: true,
	}, nil
}

// StopContainerByName 通过名称停止容器
func (s *DockerServer) StopContainerByName(ctx context.Context, req *pb.StopContainerByNameRequest) (*pb.StopContainerByNameReply, error) {
	s.log.WithContext(ctx).Infof("StopContainerByName: containerName=%s, remove=%v", req.ContainerName, req.Remove)

	err := s.uc.StopContainerByName(ctx, req.ContainerName, req.Remove)
	if err != nil {
		s.log.WithContext(ctx).Errorf("StopContainerByName failed: %v", err)
		return nil, err
	}

	return &pb.StopContainerByNameReply{
		Success: true,
	}, nil
}

// GetContainerStopTime 获取容器的停止时间戳
func (s *DockerServer) GetContainerStopTime(ctx context.Context, req *pb.GetContainerStopTimeRequest) (*pb.GetContainerStopTimeReply, error) {
	s.log.WithContext(ctx).Infof("GetContainerStopTime: containerID=%s", req.ContainerId)

	stopTime, err := s.uc.GetContainerStopTime(ctx, req.ContainerId)
	if err != nil {
		s.log.WithContext(ctx).Errorf("GetContainerStopTime failed: %v", err)
		return nil, err
	}

	return &pb.GetContainerStopTimeReply{
		StopTime: stopTime,
	}, nil
}

// FindImageByName 查找指定名称的镜像
func (s *DockerServer) FindImageByName(ctx context.Context, req *pb.FindImageByNameRequest) (*pb.FindImageByNameReply, error) {
	s.log.WithContext(ctx).Infof("FindImageByName: imageName=%s", req.ImageName)

	exists, err := s.uc.FindImageByName(ctx, req.ImageName)
	if err != nil {
		s.log.WithContext(ctx).Errorf("FindImageByName failed: %v", err)
		return nil, err
	}

	return &pb.FindImageByNameReply{
		Exists: exists,
	}, nil
}

// ImportAndTagImage 导入并标记镜像
func (s *DockerServer) ImportAndTagImage(ctx context.Context, req *pb.ImportAndTagImageRequest) (*pb.ImportAndTagImageReply, error) {
	s.log.WithContext(ctx).Infof("ImportAndTagImage: tarFilePath=%s, fullImageName=%s", req.TarFilePath, req.FullImageName)

	err := s.uc.ImportAndTagImage(ctx, req.TarFilePath, req.FullImageName)
	if err != nil {
		s.log.WithContext(ctx).Errorf("ImportAndTagImage failed: %v", err)
		return nil, err
	}

	return &pb.ImportAndTagImageReply{
		Success: true,
	}, nil
}
