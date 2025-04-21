package data

import (
	"algo-agent/internal/biz"
	"algo-agent/internal/conf"
	"algo-agent/internal/cons/file"
	"context"
	"errors"
	"net/url"
	"strconv"

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
	NewDeployServiceManagerRepo,

	NewNvidiaGpuManager,
)

// Data .
type Data struct {
	or   *OSSRepo
	rmr  *RabbitMQRepo
	dr   *DockerRepo
	dsmr *DeployServiceManagerRepo
	log  *log.Helper
}

// NewData .
func NewData(c *conf.Data, ossRepo *OSSRepo, rabbitMQRepo *RabbitMQRepo, dockerRepo *DockerRepo, deployServiceManager *DeployServiceManagerRepo, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "data"))
	d := &Data{
		or:   ossRepo,
		rmr:  rabbitMQRepo,
		dr:   dockerRepo,
		dsmr: deployServiceManager,
		log:  log,
	}

	cleanup := func() {
		d.or.Close()
		d.rmr.Close()
		d.dr.Close()
		d.dsmr.Stop(context.Background())
	}
	return d, cleanup, nil
}

// NewOSSService 创建一个新的OSS客户端，适配biz.OSSStore接口
func NewOSSRepo(c *conf.Data, logger log.Logger) (biz.OSSService, error) {
	l := log.NewHelper(log.With(logger, "module", "data/oss"))

	mc := c.Oss
	if mc == nil || mc.Endpoint == "" {
		l.Error("MinIO endpoint is not set")
		return nil, errors.New("MinIO endpoint is not set")
	}

	// 创建MinIO客户端
	client, err := minio.New(mc.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(mc.AccessKey, mc.SecretKey, ""),
		Secure: false, // 根据环境设置，这里设为false表示使用HTTP
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
func NewRabbitMQRepo(c *conf.Data, logger log.Logger) (biz.MqService, error) {
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
		buildAMQPUrl(rc.Username, rc.Password, rc.Host, rc.Port, rc.Vhost),
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
		getDynamicQueueName(rc),
		rabbitmq.WithConsumerOptionsRoutingKey(rc.DefaultRoutingKey),
		rabbitmq.WithConsumerOptionsExchangeName(rc.DefaultExchangeName),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsExchangeDurable,
	)
	if err != nil {
		l.Errorf("failed to create consumer: %v", err)
		publisher.Close()
		if conn != nil {
			_ = conn.Close()
		}
		return nil, err
	}

	// 启动消费者
	err = consumer.Run(handler)
	if err != nil {
		l.Errorf("failed to run consumer: %v", err)
		publisher.Close()
		consumer.Close()
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
func NewDockerRepo(c *conf.Data, logger log.Logger) (biz.DockerService, error) {
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

func NewDeployServiceManagerRepo(c *conf.Data, logger log.Logger) (biz.DeployServiceManager, error) {
	manager := NewDeployServiceManager(
		file.DEPLOY+file.SEPARATOR+"deploy.json",
		c.MappingFilePath,
		logger,
	)

	return manager, nil
}

// buildAMQPUrl 创建安全的AMQP URL
func buildAMQPUrl(username, password, host string, port int32, vhost string) string {
	u := &url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(username, password),
		Host:   host + ":" + strconv.Itoa(int(port)),
	}

	// vhost需要特殊处理，因为可能包含/
	if vhost != "/" {
		// 移除开头的/
		if vhost[0] == '/' {
			vhost = vhost[1:]
		}
		u.Path = vhost
	}

	return u.String()
}

// NewNvidiaGpuManager 创建一个新的NVIDIA GPU管理器
func NewNvidiaGpuManager(logger log.Logger) biz.GpuManager {
	return &NvidiaGpuManager{
		log: log.NewHelper(log.With(logger, "module", "data/gpu/nvidia")),
	}
}
