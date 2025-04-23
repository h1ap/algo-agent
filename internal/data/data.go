package data

import (
	"algo-agent/internal/biz"
	"algo-agent/internal/conf"
	"algo-agent/internal/cons/file"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

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
	NewTrainingTaskManagerRepo,
	NewExtractTaskManagerRepo,
	NewEvalTaskManagerRepo,

	NewNvidiaGpuManager,
)

// Data .
type Data struct {
	or   *OSSRepo
	rmr  *RabbitMQRepo
	dr   *DockerRepo
	dsmr *DeployServiceManagerRepo
	etmr *EvalTaskManagerRepo
	log  *log.Helper
}

// NewData .
func NewData(
	c *conf.Data,
	ossRepo *OSSRepo,
	rabbitMQRepo *RabbitMQRepo,
	dockerRepo *DockerRepo,
	deployServiceManagerRepo *DeployServiceManagerRepo,
	evalTaskManagerRepo *EvalTaskManagerRepo,
	logger log.Logger,
) (*Data, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "data"))
	d := &Data{
		or:   ossRepo,
		rmr:  rabbitMQRepo,
		dr:   dockerRepo,
		dsmr: deployServiceManagerRepo,
		etmr: evalTaskManagerRepo,
		log:  log,
	}

	cleanup := func() {
		d.or.Close()
		d.rmr.Close()
		d.dr.Close()
		d.dsmr.Stop(context.Background())
		d.etmr.Stop(context.Background())
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

	// 使用固定的OrbStack套接字
	socketPath := "/Users/heap/.orbstack/run/docker.sock"

	// 验证套接字是否存在
	if _, err := os.Stat(socketPath); err != nil {
		l.Warnf("无法访问OrbStack套接字文件：%v", err)
		socketPath = c.Docker.Host
		l.Infof("尝试使用配置的套接字：%s", socketPath)

		if _, err := os.Stat(socketPath); err != nil {
			l.Errorf("配置的套接字也无法访问：%v", err)
			return nil, fmt.Errorf("找不到可用的Docker套接字")
		}
	}

	l.Infof("使用Docker套接字：%s", socketPath)
	dockerHost := "unix://" + socketPath

	// 设置默认超时
	var responseTimeoutDuration time.Duration = 60 * time.Second
	if c.Docker != nil && c.Docker.ResponseTimeout != nil {
		responseTimeoutDuration = c.Docker.ResponseTimeout.AsDuration()
		l.Infof("使用配置的响应超时：%v", responseTimeoutDuration)
	}

	// 创建HTTP传输
	httpTransport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			d := net.Dialer{Timeout: 5 * time.Second}
			l.Infof("直接连接Unix套接字：%s", socketPath)
			conn, err := d.DialContext(ctx, "unix", socketPath)
			if err != nil {
				l.Errorf("连接Unix套接字失败：%v", err)
			}
			return conn, err
		},
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	}

	// 创建HTTP客户端
	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   responseTimeoutDuration,
	}

	// 创建Docker客户端
	cli, err := client.NewClientWithOpts(
		client.WithHTTPClient(httpClient),
		client.WithHost(dockerHost),
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		l.Errorf("创建Docker客户端失败：%v", err)
		return nil, err
	}

	// 测试连接
	pingCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 直接尝试ping一次
	l.Info("尝试Ping Docker守护进程...")
	ping, err := cli.Ping(pingCtx)
	if err != nil {
		l.Errorf("Docker守护进程Ping失败：%v", err)
		cli.Close()
		return nil, err
	}

	l.Infof("成功连接到Docker守护进程，API版本：%s", ping.APIVersion)

	return &DockerRepo{
		client: cli,
		conf:   c.Docker,
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

func NewTrainingTaskManagerRepo(c *conf.Data, logger log.Logger) (biz.TrainingTaskManager, error) {
	manager := NewTrainingTaskManager(
		file.TRAIN+file.SEPARATOR+"training.json",
		c.MappingFilePath,
		logger,
	)

	return manager, nil
}

func NewExtractTaskManagerRepo(c *conf.Data, logger log.Logger) (biz.ExtractTaskManager, error) {
	manager := NewExtractTaskManager(
		file.EXTRACT+file.SEPARATOR+"extract.json",
		c.MappingFilePath,
		logger,
	)

	return manager, nil
}

func NewEvalTaskManagerRepo(c *conf.Data, logger log.Logger) (biz.EvalTaskManager, error) {
	manager := NewEvalTaskManager(
		file.EVAL+file.SEPARATOR+"eval.json",
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
