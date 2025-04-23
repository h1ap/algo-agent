package biz

import (
	"algo-agent/internal/model/extract"
	"context"
)

// ExtractTaskManager 提取任务管理接口
type ExtractTaskManager interface {
	// AddTask 添加任务
	AddTask(ctx context.Context, task *extract.ExtractTaskInfo) error

	// RemoveTask 移除任务
	RemoveTask(ctx context.Context, id string) bool

	// GetTaskList 获取任务列表
	GetTaskList(ctx context.Context) []*extract.ExtractTaskInfo

	// FindTaskById 根据ID查找任务
	FindTaskById(ctx context.Context, taskId string) *extract.ExtractTaskInfo

	// UpdateTask 更新任务
	UpdateTask(ctx context.Context, updatedTask *extract.ExtractTaskInfo) error

	// SaveToFile 保存到文件
	SaveToFile(ctx context.Context)

	// Stop 停止管理器
	Stop(ctx context.Context)
}
