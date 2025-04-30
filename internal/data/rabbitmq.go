package data

import (
	"algo-agent/internal/conf"
	"algo-agent/internal/mq/event"
	json "algo-agent/internal/utils"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/wagslane/go-rabbitmq"
)

// RabbitMQRepo 实现RabbitMQ客户端，满足biz.MQSender接口
type RabbitMQRepo struct {
	conn      *rabbitmq.Conn
	publisher *rabbitmq.Publisher
	consumer  *rabbitmq.Consumer
	rc        *conf.Data_RabbitMQ
	nc        *conf.Data_Node
	log       *log.Helper
}

// SendMessage 发送字符串消息
func (r *RabbitMQRepo) SendMessage(ctx context.Context, exchangeName, routingKey string, message *event.ReqMessage) error {
	rk := r.rc.Group + r.rc.ServiceQueuePrefix + routingKey

	messageJson, err := json.ToJSON(message)

	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to convert object to JSON: %v", err)
		return err
	}

	err = r.publisher.Publish(
		[]byte(messageJson),
		[]string{rk},
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)

	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to publish message: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("message sent: exchangeName=%s, routingKey=%s, message=%s", exchangeName, rk, message)
	return nil
}

// SendToQueue 发送消息到队列
func (r *RabbitMQRepo) SendToQueue(ctx context.Context, queueName string, message *event.ReqMessage) error {
	rk := r.rc.Group + r.rc.NodeQueuePrefix + queueName

	messageJson, err := json.ToJSON(message)

	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to convert object to JSON: %v", err)
		return err
	}

	err = r.publisher.Publish(
		[]byte(messageJson),
		[]string{rk},
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to publish message to queue: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("message sent to queue: queueName=%s, message=%s", r.rc.DefaultExchangeName, rk, message)
	return nil
}

// SendToService 发送消息到特定服务
func (r *RabbitMQRepo) SendToService(ctx context.Context, service string, message *event.ReqMessage) error {
	rk := r.rc.Group + r.rc.ServiceQueuePrefix + service
	messageJson, err := json.ToJSON(message)

	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to convert object to JSON: %v", err)
		return err
	}

	err = r.publisher.Publish(
		[]byte(messageJson),
		[]string{rk},
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)
	if err != nil {
		r.log.WithContext(ctx).Errorf("failed to publish object message: %v", err)
		return err
	}

	r.log.WithContext(ctx).Infof("message sent to service: service=%s, message=%s", service, message)
	return nil
}

func (r *RabbitMQRepo) GetOrCreateConsumer(ctx context.Context) (*rabbitmq.Consumer, error) {

	if r.consumer != nil {
		return r.consumer, nil
	}

	consumer, err := rabbitmq.NewConsumer(
		r.conn,
		getDynamicQueueName(r.rc, r.nc),
		rabbitmq.WithConsumerOptionsRoutingKey(r.rc.DefaultRoutingKey),
		rabbitmq.WithConsumerOptionsExchangeName(r.rc.DefaultExchangeName),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsExchangeDurable,
	)
	if err != nil {
		r.log.Errorf("创建RabbitMQ消费者失败: %v", err)
		return nil, err
	}
	return consumer, nil
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

// GetConnection 获取RabbitMQ连接
func (r *RabbitMQRepo) GetConnection() *rabbitmq.Conn {
	return r.conn
}

// getDynamicQueueName 生成动态队列名称
func getDynamicQueueName(c *conf.Data_RabbitMQ, nc *conf.Data_Node) string {
	return c.Group + c.NodeQueuePrefix + nc.NodeName
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
