package service

import (
	"algo-agent/internal/utils"
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"algo-agent/internal/biz"
	"algo-agent/internal/mq/event"

	"github.com/go-kratos/kratos/v2/log"
)

// NodeOfflineServer 实现节点下线服务
type NodeOfflineServer struct {
	// 常量定义
	NODE_OFFLINE_SERVICE_KEY string

	mqService    biz.MqService
	nodeName     string
	trainService string
	log          *log.Helper
}

// Name 实现server.Server接口的Name方法
func (s *NodeOfflineServer) Name() string {
	return "node-offline"
}

// Start 实现server.Server接口的Start方法
func (s *NodeOfflineServer) Start(ctx context.Context) error {
	s.log.Info("节点下线服务已启动")
	return nil
}

// Stop 实现server.Server接口的Stop方法
func (s *NodeOfflineServer) Stop(ctx context.Context) error {
	s.log.Info("节点 " + s.nodeName + " 正在关闭，发送下线消息...")

	// 创建下线消息
	trainRPoolCloseMessage := &event.TrainRPoolCloseMessage{
		NodeName: s.nodeName,
	}

	// 序列化消息
	msgBytes, err := json.Marshal(trainRPoolCloseMessage)
	if err != nil {
		s.log.Error("序列化下线消息失败: ", err)
		return err
	}

	// 发送消息到训练服务
	err = s.mqService.SendToService(ctx, s.trainService, string(msgBytes))
	if err != nil {
		s.log.Error("发送下线消息失败: ", err)
		return err
	}

	s.log.Info("下线消息发送成功")
	return nil
}

func (j *NodeOfflineServer) Endpoint() (*url.URL, error) {
	ip := utils.GetLocalIP()
	u := &url.URL{
		Scheme: "http",
		Host:   ip + ":" + strconv.Itoa(int(8002)),
	}
	return u, nil
}
