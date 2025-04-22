package biz

import (
	"algo-agent/internal/conf"
	"algo-agent/internal/utils"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewOSSUsecase,
	NewRabbitMQUsecase,
	NewDockerUsecase,
	NewDeployUsecase,

	NewGpuUsecase,
)

// NewOSSUsecase 创建新的OSS用例
func NewOSSUsecase(store OSSService, logger log.Logger) *OSSUsecase {
	return &OSSUsecase{
		store: store,
		log:   log.NewHelper(logger),
	}
}

// NewRabbitMQUsecase 创建新的RabbitMQ用例
func NewRabbitMQUsecase(sender MqService, logger log.Logger) *RabbitMQUsecase {
	return &RabbitMQUsecase{
		sender: sender,
		log:    log.NewHelper(logger),
	}
}

// NewDockerUsecase 创建新的Docker用例
func NewDockerUsecase(docker DockerService, logger log.Logger) *DockerUsecase {
	return &DockerUsecase{
		docker: docker,
		log:    log.NewHelper(logger),
	}
}

func NewDeployUsecase(
	cfg *conf.Data,
	dsm DeployServiceManager,
	mq MqService,
	d DockerService,
	oss OSSService,
	logger log.Logger,
) *DeployUsecase {
	return &DeployUsecase{
		dsm: dsm,
		mq:  mq,
		d:   d,
		oss: oss,
		log: log.NewHelper(logger),

		filePath:        cfg.MappingFilePath,
		mappingFilePath: cfg.MappingFilePath,
		dsn:             cfg.Services.Deploy,
		isn:             "inference.py",
	}
}

func NewGpuUsecase(cfg *conf.Data, g GpuManager, mq MqService, logger log.Logger) *GpuUsecase {
	ctx, cancel := context.WithCancel(context.Background())
	// 获取节点名称，优先使用Node配置
	nodeName := ""
	if cfg.Node != nil {
		nodeName = cfg.Node.NodeName
	} else if cfg.Rabbitmq != nil {
		nodeName = cfg.Rabbitmq.NodeName
	}

	return &GpuUsecase{
		g:              g,
		mq:             mq,
		tsn:            cfg.Services.Train,
		nodeName:       nodeName,
		ipAddress:      utils.GetLocalIP(),
		ctx:            ctx,
		cancel:         cancel,
		reportInterval: 20 * time.Second, // 默认20秒上报一次
		log:            log.NewHelper(logger),
	}
}
