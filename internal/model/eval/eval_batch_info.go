package eval

// EvalBatchInfo 评估信息
type EvalBatchInfo struct {
	// TaskId 任务ID
	TaskId string `json:"taskId"`

	// Details 评估详情列表
	Details []EvalDetail `json:"details"`
}

// EvalDetail 评估详情
type EvalDetail struct {
	// DataUuid 数据UUID
	DataUuid string `json:"dataUuid"`

	// EvalData 评估数据
	EvalData string `json:"evalData"`
}
