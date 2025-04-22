package eval

import "algo-agent/internal/mq/event"

// EvalTaskInfo 评估任务信息
type EvalTaskInfo struct {
	// TaskId 任务id
	TaskId string `json:"taskId"`

	// TrainingContainerName 训练容器名称（唯一），不需要传入
	TrainingContainerName string `json:"trainingContainerName"`

	// TaskStatus 当前任务状态
	TaskStatus int32 `json:"taskStatus"`

	// Remark 备注
	Remark string `json:"remark"`

	// EvalSendMessage 评估任务信息
	EvalSendMessage *event.EvalSendMessage `json:"evalSendMessage"`
}

// NewEvalTaskInfo 创建EvalTaskInfo的构造函数
func NewEvalTaskInfo() *EvalTaskInfo {
	return &EvalTaskInfo{}
}

// NewEvalTaskInfoWithMessage 使用评估任务请求消息创建EvalTaskInfo
func NewEvalTaskInfoWithMessage(message *event.EvalSendMessage) *EvalTaskInfo {
	return &EvalTaskInfo{
		TaskId:          message.TaskId,
		EvalSendMessage: message,
	}
}
