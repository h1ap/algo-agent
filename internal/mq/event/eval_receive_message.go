package event

// EvalReceiveMessage 评估作业接收消息 接收
type EvalReceiveMessage struct {
	// 评估作业主键
	TaskId string `json:"taskId"`

	// 状态：1 进行中， 2 已终止，3 失败， 4 成功
	Status int32 `json:"status"`

	// 备注
	Remark string `json:"remark"`

	// 更新详情
	DetailList []EvalDetail `json:"detailList"`

	// 更新结果
	Result string `json:"result"`
}

// EvalDetail 评估详情结构体
type EvalDetail struct {
	// 关联数据的uuid
	DataUUID string `json:"dataUuid"`

	// 评估数据
	EvalData string `json:"evalData"`
}
