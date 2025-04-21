package service

import (
	"algo-agent/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewOSSServer,
	NewDockerServer,
	NewRabbitMQServer,
	NewDeployServer,
	NewJobServer,
)

// 当proto文件生成后，此部分注册推理服务的API实现
func NewDeployServer(uc *biz.DeployUsecase, logger log.Logger) *DeployServer {
	return &DeployServer{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// 当proto文件生成后，此部分注册OSS服务的API实现
func NewOSSServer(uc *biz.OSSUsecase, logger log.Logger) *OSSServer {
	return &OSSServer{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// NewJobServer 创建定时任务管理器
func NewJobServer(gpuUC *biz.GpuUsecase, logger log.Logger) *JobServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &JobServer{
		gpuUC:  gpuUC,
		log:    log.NewHelper(log.With(logger, "module", "biz/job-manager")),
		ctx:    ctx,
		cancel: cancel,
	}
}
