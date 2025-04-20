package data

import (
	"algo-agent/internal/conf"
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/wagslane/go-rabbitmq"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewOSSRepo,
	NewRabbitMQRepo,
	NewDockerRepo,
)

// Data .
type Data struct {
	or  *OSSRepo
	rmr *RabbitMQRepo
	dr  *DockerRepo
	log *log.Helper
}

// NewData .
func NewData(c *conf.Data, ossRepo *OSSRepo, rabbitMQRepo *RabbitMQRepo, dockerRepo *DockerRepo, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "data"))
	d := &Data{
		or:  ossRepo,
		rmr: rabbitMQRepo,
		dr:  dockerRepo,
		log: log,
	}

	cleanup := func() {
		d.or.Close()
		d.rmr.Close()
		d.dr.Close()
	}
	return d, cleanup, nil
}

// NewOSSService 创建一个新的OSS客户端，适配biz.OSSStore接口
func NewOSSRepo(c *conf.Data, logger log.Logger) (*OSSRepo, error) {
	l := log.NewHelper(log.With(logger, "module", "data/oss"))

	mc := c.Oss
	if mc == nil || mc.Endpoint == "" {
		l.Error("MinIO endpoint is not set")
		return nil, errors.New("MinIO endpoint is not set")
	}

	// 创建MinIO客户端
	client, err := minio.New(mc.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(mc.AccessKey, mc.SecretKey, ""),
	})
	if err != nil {
		l.Errorf("failed to create MinIO client: %v", err)
		return nil, err
	}

	repo := &OSSRepo{
		client: client,
		conf:   mc,
		log:    l,
	}

	return repo, nil
}

// NewRabbitMQRepo 创建一个新的RabbitMQ客户端，适配biz.MQSender接口
func NewRabbitMQRepo(c *conf.Data, logger log.Logger) (*RabbitMQRepo, error) {
	l := log.NewHelper(log.With(logger, "module", "data/rabbitmq"))
	rabbitLogger := &rabbitLogger{log: l}

	rc := c.Rabbitmq
	if rc == nil || rc.Host == "" {
		l.Error("RabbitMQ host is not set")
		return nil, errors.New("RabbitMQ host is not set")
	}

	// 填充配置
	cfg := rabbitmq.Config{
		Vhost: rc.Vhost,
	}

	// 建立连接
	conn, err := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s:%s", rc.Username, rc.Password, rc.Host, rc.Port),
		rabbitmq.WithConnectionOptionsLogger(rabbitLogger),
		rabbitmq.WithConnectionOptionsConfig(cfg),
	)
	if err != nil {
		l.Errorf("failed to connect to RabbitMQ: %v", err)
		if conn != nil {
			_ = conn.Close()
		}
		return nil, err
	}

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogger(rabbitLogger),
		rabbitmq.WithPublisherOptionsExchangeName(rc.DefaultExchangeName),
		rabbitmq.WithPublisherOptionsExchangeKind("direct"),
		rabbitmq.WithPublisherOptionsExchangeDurable,
	)
	if err != nil {
		l.Errorf("failed to create publisher: %v", err)
		if conn != nil {
			_ = conn.Close()
		}
		return nil, err
	}

	// 创建消费者处理器
	handler := func(d rabbitmq.Delivery) rabbitmq.Action {
		l.Infof("Received message: %s", string(d.Body))
		return rabbitmq.Ack
	}

	consumer, err := rabbitmq.NewConsumer(
		conn,
		handler,
		getDynamicQueueName(rc),
		rabbitmq.WithConsumerOptionsRoutingKey(rc.DefaultRoutingKey),
		rabbitmq.WithConsumerOptionsExchangeName(rc.DefaultExchangeName),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		l.Errorf("failed to create consumer: %v", err)
		publisher.Close()
		if conn != nil {
			_ = conn.Close()
		}
		return nil, err
	}

	repo := &RabbitMQRepo{
		conn:      conn,
		publisher: publisher,
		consumer:  consumer,
		conf:      rc,
		log:       l,
	}

	return repo, nil
}

// NewDockerRepo 创建一个新的Docker客户端，适配biz.DockerService接口
func NewDockerRepo(c *conf.Data, logger log.Logger) (*DockerRepo, error) {
	l := log.NewHelper(log.With(logger, "module", "data/docker"))

	dockerConf := c.Docker
	if dockerConf == nil || dockerConf.Host == "" {
		l.Warn("Docker host is not set, using default")
		dockerConf = &conf.Data_Docker{
			Host: "unix:///var/run/docker.sock",
		}
	}

	// 创建Docker客户端
	cli, err := client.NewClientWithOpts(
		client.WithHost(dockerConf.Host),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		l.Errorf("failed to create Docker client: %v", err)
		return nil, err
	}

	// 测试连接
	_, err = cli.Ping(context.Background())
	if err != nil {
		l.Errorf("failed to connect to Docker daemon: %v", err)
		cli.Close()
		return nil, err
	}

	l.Info("Docker connection successful")

	return &DockerRepo{
		client: cli,
		conf:   dockerConf,
		log:    l,
		logStreams: make(map[string]struct {
			close     func() error
			isRunning bool
		}),
	}, nil
}
