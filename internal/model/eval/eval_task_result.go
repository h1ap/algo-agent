package eval

import (
	"time"
)

// EvalTaskResult 评估结果
type EvalTaskResult struct {
	// TaskId 任务ID
	TaskId string `json:"taskId"`

	// OverallMetrics 整体指标
	OverallMetrics map[string]interface{} `json:"overallMetrics"`

	// ClassifyMetrics 分类指标
	// 每一个分类指标，都应应该包含一个 label字段，标识分类标签
	ClassifyMetrics []map[string]interface{} `json:"classifyMetrics"`

	// CreateTime 创建时间
	CreateTime time.Time `json:"createTime"`
}
