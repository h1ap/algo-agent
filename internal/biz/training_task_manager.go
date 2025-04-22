package biz

import (
	"algo-agent/internal/model/train"
	"context"
)

// TrainingTaskManager 训练任务管理接口
type TrainingTaskManager interface {
	// AddTask 添加任务
	AddTask(ctx context.Context, task *train.TrainingTaskInfo) error

	// RemoveTask 移除任务
	RemoveTask(ctx context.Context, id string) bool

	// GetTaskList 获取任务列表
	GetTaskList(ctx context.Context) []*train.TrainingTaskInfo

	// FindTaskById 根据ID查找任务
	FindTaskById(ctx context.Context, taskId string) *train.TrainingTaskInfo

	// UpdateTask 更新任务
	UpdateTask(ctx context.Context, updatedTask *train.TrainingTaskInfo) error

	// SaveToFile 保存到文件
	SaveToFile(ctx context.Context)

	// Stop 停止管理器
	Stop(ctx context.Context)
}
