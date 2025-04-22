package event

// TrainRPoolCloseMessage 训练资源池关闭消息
type TrainRPoolCloseMessage struct {
	// 训练资源池名称
	NodeName string `json:"nodeName"`
}
