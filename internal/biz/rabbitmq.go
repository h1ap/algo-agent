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
	SendMessage(ctx context.Context, exchangeName, routingKey, message string) error

	// SendToQueue 发送消息到队列
	SendToQueue(ctx context.Context, queueName, message string) error

	// SendToService 发送消息到特定服务
	SendToService(ctx context.Context, service string, message string) error

	// Close 关闭连接
	Close()
}

// RabbitMQUsecase 是RabbitMQ用例
type RabbitMQUsecase struct {
	sender MQSender
	log    *log.Helper
}

// SendMessage 发送消息
func (uc *RabbitMQUsecase) SendMessage(ctx context.Context, exchangeName, routingKey, message string) error {
	uc.log.WithContext(ctx).Infof("SendMessage: exchangeName=%s, routingKey=%s, message=%s", exchangeName, routingKey, message)
	return uc.sender.SendMessage(ctx, exchangeName, routingKey, message)
}

// SendToQueue 发送消息到队列
func (uc *RabbitMQUsecase) SendToQueue(ctx context.Context, queueName, message string) error {
	uc.log.WithContext(ctx).Infof("SendToQueue: queueName=%s, message=%s", queueName, message)
	return uc.sender.SendToQueue(ctx, queueName, message)
}

// SendToService 发送消息到特定服务
func (uc *RabbitMQUsecase) SendToService(ctx context.Context, service string, message string) error {
	uc.log.WithContext(ctx).Infof("SendToService: service=%s, message=%s", service, message)
	return uc.sender.SendToService(ctx, service, message)
}
