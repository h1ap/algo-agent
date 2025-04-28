package biz

import (
	"algo-agent/internal/conf"
	"algo-agent/internal/utils"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewOSSUsecase,
	NewRabbitMQUsecase,
	NewDockerUsecase,
	NewDeployUsecase,
	NewTrainingTaskUsecase,
	NewExtractTaskUsecase,
	NewEvalTaskUsecase,
	NewGpuUsecase,
	NewTaskCheckerUsecase,
)

// NewOSSUsecase 创建新的OSS用例
func NewOSSUsecase(store OSSService, logger log.Logger) *OSSUsecase {
	return &OSSUsecase{
		store: store,
		log:   log.NewHelper(logger),
	}
}

// NewRabbitMQUsecase 创建新的RabbitMQ用例
func NewRabbitMQUsecase(
	sender MqService,
	trainingUsecase *TrainingTaskUsecase,
	evalUsecase *EvalTaskUsecase,
	deployUsecase *DeployUsecase,
	extractUsecase *ExtractTaskUsecase,
	logger log.Logger,
) *RabbitMQUsecase {
	return &RabbitMQUsecase{
		mq:              sender,
		trainingUsecase: trainingUsecase,
		evalUsecase:     evalUsecase,
		deployUsecase:   deployUsecase,
		extractUsecase:  extractUsecase,
		log:             log.NewHelper(logger),
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

// NewTrainingTaskUsecase 创建新的训练任务用例
func NewTrainingTaskUsecase(
	cfg *conf.Data,
	ttm TrainingTaskManager,
	mq MqService,
	d DockerService,
	oss OSSService,
	logger log.Logger,
) *TrainingTaskUsecase {
	return &TrainingTaskUsecase{
		ttm:                  ttm,
		mq:                   mq,
		d:                    d,
		oss:                  oss,
		log:                  log.NewHelper(logger),
		filePath:             cfg.MappingFilePath,
		tsn:                  cfg.Services.Train,
		checkpointPathPrefix: "Checkpoint/train_id-",
		modelPathPrefix:      "Model/train_id-",
		trainScriptName:      "train.py",
	}
}

// NewExtractTaskUsecase 创建新的提取任务用例
func NewExtractTaskUsecase(
	cfg *conf.Data,
	etm ExtractTaskManager,
	mq MqService,
	d DockerService,
	oss OSSService,
	logger log.Logger,
) *ExtractTaskUsecase {
	return &ExtractTaskUsecase{
		etm:               etm,
		mq:                mq,
		d:                 d,
		oss:               oss,
		log:               log.NewHelper(logger),
		filePath:          cfg.MappingFilePath,
		tsn:               cfg.Services.Train,
		extractScriptName: "extract.py",
	}
}

// NewEvalTaskUsecase 创建新的评估任务用例
func NewEvalTaskUsecase(
	cfg *conf.Data,
	etm EvalTaskManager,
	mq MqService,
	d DockerService,
	oss OSSService,
	logger log.Logger,
) *EvalTaskUsecase {
	return &EvalTaskUsecase{
		etm:            etm,
		mq:             mq,
		d:              d,
		oss:            oss,
		log:            log.NewHelper(logger),
		filePath:       cfg.MappingFilePath,
		tsn:            cfg.Services.Train,
		evalScriptName: "eval.py",
	}
}

func NewGpuUsecase(cfg *conf.Data, g GpuManager, mq MqService, logger log.Logger) *GpuUsecase {
	// 获取节点名称，优先使用Node配置
	nodeName := ""
	if cfg.Node != nil {
		nodeName = cfg.Node.NodeName
	} else if cfg.Rabbitmq != nil {
		nodeName = cfg.Rabbitmq.NodeName
	}

	return &GpuUsecase{
		g:         g,
		mq:        mq,
		tsn:       cfg.Services.Train,
		nodeName:  nodeName,
		ipAddress: utils.GetLocalIP(),
		log:       log.NewHelper(logger),
	}
}

// NewTaskCheckerUsecase 创建新的任务检查器用例
func NewTaskCheckerUsecase(
	training *TrainingTaskUsecase,
	eval *EvalTaskUsecase,
	deploy *DeployUsecase,
	extract *ExtractTaskUsecase,
	gpu *GpuUsecase,
	logger log.Logger,
) *TaskCheckerUsecase {
	return &TaskCheckerUsecase{
		trainingUsecase: training,
		evalUsecase:     eval,
		deployUsecase:   deploy,
		extractUsecase:  extract,
		gpuUsecase:      gpu,
		log:             log.NewHelper(logger),
	}
}
