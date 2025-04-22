package extract

// ExtractTaskResult 提取任务结果
type ExtractTaskResult struct {
	// TaskId 任务ID
	TaskId string `json:"taskId"`

	// ModelPath 提取的模型文件路径【绝对路径，或是/workspace 目录下的相对路径】
	ModelPath string `json:"modelPath"`
}
