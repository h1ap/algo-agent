package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Message 消息模型
type Message struct {
	Content string
}

// MQSender 消息队列发送接口
type MQSender interface {
	// SendMessage 发送字符串消息
	SendMessage(ctx context.Context, routingKey, message string) error

	// SendObjectMessage 发送对象消息
	SendObjectMessage(ctx context.Context, routingKey string, object interface{}) error

	// SendToService 发送消息到特定服务
	SendToService(ctx context.Context, service string, object interface{}) error

	// Close 关闭连接
	Close()
}

// RabbitMQUsecase 是RabbitMQ用例
type RabbitMQUsecase struct {
	sender MQSender
	log    *log.Helper
}

// NewRabbitMQUsecase 创建新的RabbitMQ用例
func NewRabbitMQUsecase(sender MQSender, logger log.Logger) *RabbitMQUsecase {
	return &RabbitMQUsecase{
		sender: sender,
		log:    log.NewHelper(logger),
	}
}

// SendMessage 发送消息
func (uc *RabbitMQUsecase) SendMessage(ctx context.Context, routingKey, message string) error {
	uc.log.WithContext(ctx).Infof("SendMessage: routingKey=%s", routingKey)
	return uc.sender.SendMessage(ctx, routingKey, message)
}

// SendObjectMessage 发送对象消息
func (uc *RabbitMQUsecase) SendObjectMessage(ctx context.Context, routingKey string, object interface{}) error {
	uc.log.WithContext(ctx).Infof("SendObjectMessage: routingKey=%s", routingKey)
	return uc.sender.SendObjectMessage(ctx, routingKey, object)
}

// SendToService 发送消息到特定服务
func (uc *RabbitMQUsecase) SendToService(ctx context.Context, service string, object interface{}) error {
	uc.log.WithContext(ctx).Infof("SendToService: service=%s", service)
	return uc.sender.SendToService(ctx, service, object)
}
