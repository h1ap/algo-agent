package service

import (
	"algo-agent/internal/biz"
	"algo-agent/internal/utils"
	"context"
	"net/url"
	"strconv"
)

// RabbitMQConsumerServer 实现RabbitMQ消费者
type RabbitMQConsumerServer struct {
	uc *biz.RabbitMQUsecase
}

// Name 实现server.Server接口的Name方法
func (rmc *RabbitMQConsumerServer) Name() string {
	return "rabbitmq-consumer"
}

// Start 实现server.Server接口的Start方法
func (rmc *RabbitMQConsumerServer) Start(ctx context.Context) error {
	rmc.uc.Subscribe(ctx)
	return nil
}

// Stop 实现server.Server接口的Stop方法
func (rmc *RabbitMQConsumerServer) Stop(ctx context.Context) error {
	rmc.uc.Unsubscribe()
	return nil
}

// Endpoint 实现server.Server接口的Endpoint方法
func (rmc *RabbitMQConsumerServer) Endpoint() (*url.URL, error) {
	ip := utils.GetLocalIP()
	u := &url.URL{
		Scheme: "http",
		Host:   ip + ":" + strconv.Itoa(int(8003)),
	}
	return u, nil
}
