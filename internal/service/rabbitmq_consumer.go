package service

import (
	"algo-agent/internal/biz"
	"algo-agent/internal/utils"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"net/url"
	"strconv"
)

// RabbitMQConsumerServer 实现RabbitMQ消费者
type RabbitMQConsumerServer struct {
	uc  *biz.RabbitMQUsecase
	log *log.Helper
}

// Name 实现server.Server接口的Name方法
func (rmc *RabbitMQConsumerServer) Name() string {
	return "rabbitmq-consumer"
}

// Start 实现server.Server接口的Start方法
func (rmc *RabbitMQConsumerServer) Start(ctx context.Context) error {
	err := rmc.uc.Subscribe(ctx)
	if err != nil {
		rmc.log.WithContext(ctx).Error("rabbitmqConsumerServer Start failed")
		return err
	}
	rmc.log.WithContext(ctx).Infof("rabbitmqConsumerServer Started")
	return nil
}

// Stop 实现server.Server接口的Stop方法
func (rmc *RabbitMQConsumerServer) Stop(ctx context.Context) error {
	rmc.uc.Unsubscribe()
	rmc.log.WithContext(ctx).Infof("rabbitmqConsumerServer Stop")
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
