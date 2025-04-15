package service

import (
	"context"

	pb "algo-agent/api/rabbitmq/v1"
	"algo-agent/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// RabbitMQServer 实现RabbitMQ服务API接口
type RabbitMQServer struct {
	pb.UnimplementedRabbitMQServiceServer
	uc  *biz.RabbitMQUsecase
	log *log.Helper
}

// NewRabbitMQServer 创建RabbitMQ服务实例
func NewRabbitMQServer(uc *biz.RabbitMQUsecase, logger log.Logger) *RabbitMQServer {
	return &RabbitMQServer{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// SendMessage 发送消息到交换机和路由键
func (s *RabbitMQServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageReply, error) {
	s.log.WithContext(ctx).Infof("SendMessage: exchangeName=%s, routingKey=%s", req.ExchangeName, req.RoutingKey)

	err := s.uc.SendMessage(ctx, req.ExchangeName, req.RoutingKey, req.Message)
	if err != nil {
		s.log.WithContext(ctx).Errorf("SendMessage failed: %v", err)
		return nil, err
	}

	return &pb.SendMessageReply{Success: true}, nil
}

// SendToQueue 发送消息到队列
func (s *RabbitMQServer) SendToQueue(ctx context.Context, req *pb.SendToQueueRequest) (*pb.SendToQueueReply, error) {
	s.log.WithContext(ctx).Infof("SendToQueue: queueName=%s", req.QueueName)

	err := s.uc.SendToQueue(ctx, req.QueueName, req.Message)
	if err != nil {
		s.log.WithContext(ctx).Errorf("SendToQueue failed: %v", err)
		return nil, err
	}

	return &pb.SendToQueueReply{Success: true}, nil
}

// SendToService 发送消息到特定服务
func (s *RabbitMQServer) SendToService(ctx context.Context, req *pb.SendToServiceRequest) (*pb.SendToServiceReply, error) {
	s.log.WithContext(ctx).Infof("SendToService: service=%s", req.Service)

	err := s.uc.SendToService(ctx, req.Service, req.Message)
	if err != nil {
		s.log.WithContext(ctx).Errorf("SendToService failed: %v", err)
		return nil, err
	}

	return &pb.SendToServiceReply{Success: true}, nil
}
