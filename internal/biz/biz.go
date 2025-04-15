package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewOSSUsecase, NewRabbitMQUsecase, NewDockerUsecase)

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
