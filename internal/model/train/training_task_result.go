package train

// TrainingTaskResult 训练结果
type TrainingTaskResult struct {
	// TaskId 任务id
	TaskId string `json:"taskId"`

	// BestEpoch 最优epoch
	BestEpoch int32 `json:"bestEpoch,omitempty"`

	// BestModelPath 最优模型权重路径【绝对路径，或是/workspace 目录下的相对路径】
	BestModelPath string `json:"bestModelPath"`

	// FinalModelPath 最后一次模型权重路径【绝对路径，或是/workspace 目录下的相对路径】
	FinalModelPath string `json:"finalModelPath"`
}
