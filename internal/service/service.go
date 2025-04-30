package service

import (
	"algo-agent/internal/biz"
	"algo-agent/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"math/rand"
	"time"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewOSSServer,
	NewDockerServer,
	NewRabbitMQServer,
	NewDeployServer,
	NewTrainServer,
	NewEvalServer,
	NewExtractServer,
	NewNodeOfflineServer,
	NewRabbitMQConsumerServer,
	NewCronServer,
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

// NewTrainServer 创建训练信息服务实例
func NewTrainServer(ttu *biz.TrainingTaskUsecase, logger log.Logger) *TrainServer {
	return &TrainServer{
		ttu: ttu,
		log: log.NewHelper(logger),
	}
}

// NewEvalServer 创建评估服务实例
func NewEvalServer(uc *biz.EvalTaskUsecase, logger log.Logger) *EvalServer {
	return &EvalServer{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// NewExtractServer 创建提取任务服务实例
func NewExtractServer(uc *biz.ExtractTaskUsecase, logger log.Logger) *ExtractServer {
	return &ExtractServer{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// NewNodeOfflineServer 创建节点下线服务实例
func NewNodeOfflineServer(cfg *conf.Data, mqService biz.MqService, logger log.Logger) *NodeOfflineServer {
	// 获取节点名称，优先使用Node配置
	source := rand.NewSource(time.Now().UnixNano())
	// 创建一个新的 Rand 实例
	r := rand.New(source)
	nodeName := "unknow-node-" + string(rune(r.Intn(100)))
	if cfg.Node != nil {
		nodeName = cfg.Node.NodeName
	}

	return &NodeOfflineServer{
		NODE_OFFLINE_SERVICE_KEY: "nodeOfflineService",
		mqService:                mqService,
		nodeName:                 nodeName,
		trainService:             cfg.Services.Train,
		log:                      log.NewHelper(log.With(logger, "module", "service/node-offline")),
	}
}

// NewRabbitMQConsumerServer 创建一个新的RabbitMQ消费者
func NewRabbitMQConsumerServer(
	uc *biz.RabbitMQUsecase,
	logger log.Logger,
) (*RabbitMQConsumerServer, error) {
	l := log.NewHelper(log.With(logger, "module", "service/rabbitmq-consumer"))

	rmc := &RabbitMQConsumerServer{
		uc:  uc,
		log: l,
	}

	l.Info("RabbitMQ消费者服务器创建成功")
	return rmc, nil
}
