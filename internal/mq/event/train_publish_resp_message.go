package event

// TrainPublishRespMessage 训练发布响应消息体
type TrainPublishRespMessage struct {
	// 训练详情id
	TrainDetailId string `json:"trainDetailId"`

	// 模型版本Id
	ModelVersionId int64 `json:"modelVersionId"`

	// 当前任务状态
	// see cons/train/train_job_status_enum.go
	TaskStatus int32 `json:"taskStatus"`

	// 错误信息
	Remark string `json:"remark"`

	// 权重文件路径
	ModelWeightsPath string `json:"modelWeightsPath"`
}
