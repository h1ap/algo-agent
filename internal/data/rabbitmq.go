package data

import (
	"context"
	"encoding/json"
	"errors"

	"algo-agent/internal/biz"
	"algo-agent/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/wagslane/go-rabbitmq"
)

// RabbitMQRepo 实现RabbitMQ客户端，满足biz.MQSender接口
type RabbitMQRepo struct {
	conn      *rabbitmq.Conn
	publisher *rabbitmq.Publisher
	consumer  *rabbitmq.Consumer
	conf      *conf.Data_RabbitMQ
	log       *log.Helper
}

// NewRabbitMQRepo 创建一个新的RabbitMQ客户端，适配biz.MQSender接口
func NewRabbitMQRepo(c *conf.Data, logger log.Logger) (biz.MQSender, error) {
	l := log.NewHelper(log.With(logger, "module", "data/rabbitmq"))
	rabbitLogger := &rabbitLogger{log: l}

	rabbitmqConf := c.Rabbitmq
	if rabbitmqConf == nil || rabbitmqConf.Url == "" {
		l.Error("RabbitMQ URL is not set")
		return nil, errors.New("RabbitMQ URL is not set")
	}

	// 填充配置
	cfg := rabbitmq.Config{
		Vhost: rabbitmqConf.Vhost,
	}

	// 建立连接
	conn, err := rabbitmq.NewConn(
		rabbitmqConf.Url,
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
		rabbitmq.WithPublisherOptionsExchangeName(rabbitmqConf.ExchangeName),
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
		getDynamicQueueName(rabbitmqConf),
		rabbitmq.WithConsumerOptionsRoutingKey(rabbitmqConf.RoutingKey),
		rabbitmq.WithConsumerOptionsExchangeName(rabbitmqConf.ExchangeName),
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

	return &RabbitMQRepo{
		conn:      conn,
		publisher: publisher,
		consumer:  consumer,
		conf:      rabbitmqConf,
		log:       l,
	}, nil
}

// SendMessage 发送字符串消息
func (r *RabbitMQRepo) SendMessage(ctx context.Context, routingKey, message string) error {
	rk := r.conf.Group + r.conf.ServiceQueuePrefix + routingKey

	err := r.publisher.Publish(
		[]byte(message),
		[]string{rk},
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)

	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to publish message: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("message sent: routingKey=%s", rk)
	return nil
}

// SendObjectMessage 发送对象消息
func (r *RabbitMQRepo) SendObjectMessage(ctx context.Context, routingKey string, object interface{}) error {
	rk := r.conf.Group + r.conf.ServiceQueuePrefix + routingKey

	jsonData, err := json.Marshal(object)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to marshal object message: %v", err)
		return err
	}

	err = r.publisher.Publish(
		jsonData,
		[]string{rk},
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to publish object message: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("object message sent: routingKey=%s, type=%T", rk, object)
	return nil
}

// SendToService 发送消息到特定服务
func (r *RabbitMQRepo) SendToService(ctx context.Context, service string, object interface{}) error {
	rk := r.conf.Group + r.conf.ServiceQueuePrefix + service

	jsonData, err := json.Marshal(object)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to marshal object message: %v", err)
		return err
	}

	err = r.publisher.Publish(
		jsonData,
		[]string{rk},
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to publish object message: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("object message sent to service: service=%s, type=%T", service, object)
	return nil
}

// Close 关闭连接
func (r *RabbitMQRepo) Close() {
	if r.publisher != nil {
		r.publisher.Close()
	}

	if r.consumer != nil {
		r.consumer.Close()
	}

	if r.conn != nil {
		err := r.conn.Close()
		if err != nil {
			r.log.Errorf("failed to close RabbitMQ connection: %v", err)
		}
	}
}

// getDynamicQueueName 生成动态队列名称
func getDynamicQueueName(c *conf.Data_RabbitMQ) string {
	return c.Group + c.NodeQueuePrefix + c.NodeName
}

// rabbitLogger 适配器，将kratos日志转换为rabbitmq日志
type rabbitLogger struct {
	log *log.Helper
}

func (l *rabbitLogger) Tracef(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *rabbitLogger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *rabbitLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *rabbitLogger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *rabbitLogger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *rabbitLogger) Fatalf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}
