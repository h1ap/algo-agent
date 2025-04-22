package train

// TrainingCheckpoint 训练任务保存检查点信息
type TrainingCheckpoint struct {
	// Epoch 当前训练轮次
	Epoch int32 `json:"epoch"`

	// TaskId 任务ID
	TaskId string `json:"taskId"`

	// CheckpointPath 检查点路径 【绝对路径，或是/workspace 目录下的相对路径】
	CheckpointPath string `json:"checkpointPath"`
}
