package biz

import (
	"algo-agent/internal/conf"
	"algo-agent/internal/data"

	ddr "algo-agent/internal/data/deploy"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewOSSUsecase,
	NewRabbitMQUsecase,
	NewDockerUsecase,
)

// NewOSSUsecase 创建新的OSS用例
func NewOSSUsecase(store OSSStore, logger log.Logger) *OSSUsecase {
	return &OSSUsecase{
		store: store,
		log:   log.NewHelper(logger),
	}
}

// NewRabbitMQUsecase 创建新的RabbitMQ用例
func NewRabbitMQUsecase(sender MQSender, logger log.Logger) *RabbitMQUsecase {
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
	dsm *ddr.DeployServiceManager,
	rmr *data.RabbitMQRepo,
	dr *data.DockerRepo,
	or *data.OSSRepo,
	logger log.Logger,
) *DeployUsecase {
	return &DeployUsecase{
		dsm: dsm,
		rmr: rmr,
		dr:  dr,
		or:  or,
		log: log.NewHelper(logger),

		filePath:        cfg.MappingFilePath,
		mappingFilePath: cfg.MappingFilePath,
		dsn:             cfg.Services.Deploy,
		isn:             "inference.py",
	}
}
