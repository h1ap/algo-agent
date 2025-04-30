package event

type ReqMessage struct {
	// 消息类型 mq.MqMessageTypeEnum
	Type int `json:"type"`

	// 消息体
	Payload interface{} `json:"payload"`
}
