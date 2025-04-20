package service

import (
	"algo-agent/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewOSSServer, NewDockerServer, NewRabbitMQServer)

// 当proto文件生成后，此部分注册OSS服务的API实现
func NewOSSServer(uc *biz.OSSUsecase, logger log.Logger) *OSSServer {
	return &OSSServer{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}
