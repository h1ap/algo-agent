package service

import (
	"context"

	v1 "algo-agent/api/rabbitmq/v1"
	"algo-agent/internal/biz"
)

// RabbitMQService 是一个RabbitMQ服务
type RabbitMQService struct {
	v1.UnimplementedRabbitMQServer

	uc *biz.RabbitMQUsecase
}

// NewRabbitMQService 创建新的RabbitMQ服务
func NewRabbitMQService(uc *biz.RabbitMQUsecase) *RabbitMQService {
	return &RabbitMQService{uc: uc}
}

// SendMessage 实现rabbitmq.RabbitMQServer
func (s *RabbitMQService) SendMessage(ctx context.Context, in *v1.SendMessageRequest) (*v1.SendMessageReply, error) {
	err := s.uc.SendMessage(ctx, in.RoutingKey, in.Message)
	if err != nil {
		return nil, err
	}
	return &v1.SendMessageReply{Success: true}, nil
}

// SendToService 实现rabbitmq.RabbitMQServer
func (s *RabbitMQService) SendToService(ctx context.Context, in *v1.SendToServiceRequest) (*v1.SendToServiceReply, error) {
	message := &biz.Message{
		Content: in.Message,
	}
	err := s.uc.SendToService(ctx, in.Service, message)
	if err != nil {
		return nil, err
	}
	return &v1.SendToServiceReply{Success: true}, nil
}
