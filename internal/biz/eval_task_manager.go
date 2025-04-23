package biz

import (
	"algo-agent/internal/model/eval"
	"context"
)

// EvalTaskManager 评估任务管理接口
type EvalTaskManager interface {
	// AddTask 添加任务
	AddTask(ctx context.Context, task *eval.EvalTaskInfo) error

	// RemoveTask 移除任务
	RemoveTask(ctx context.Context, id string) bool

	// GetTaskList 获取任务列表
	GetTaskList(ctx context.Context) []*eval.EvalTaskInfo

	// FindTaskById 根据ID查找任务
	FindTaskById(ctx context.Context, taskId string) *eval.EvalTaskInfo

	// UpdateTask 更新任务
	UpdateTask(ctx context.Context, updatedTask *eval.EvalTaskInfo) error

	// SaveToFile 保存到文件
	SaveToFile(ctx context.Context)

	// Stop 停止管理器
	Stop(ctx context.Context)
}
