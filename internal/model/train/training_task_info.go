package train

import "algo-agent/internal/mq/event"

// TrainingTaskInfo 训练任务信息
type TrainingTaskInfo struct {
	// TaskId 任务id
	TaskId string `json:"taskId"`

	// TrainingContainerName 训练容器名称（唯一），不需要传入
	TrainingContainerName string `json:"trainingContainerName"`

	// TaskStatus 当前任务状态
	TaskStatus int32 `json:"taskStatus"`

	// Remark 备注
	Remark string `json:"remark"`

	// TrainTaskReqMessage 训练任务参数信息
	TrainTaskReqMessage *event.TrainTaskReqMessage `json:"trainTaskReqMessage"`
}

// NewTrainingTaskInfo 创建TrainingTaskInfo的构造函数
func NewTrainingTaskInfo() *TrainingTaskInfo {
	return &TrainingTaskInfo{}
}

// NewTrainingTaskInfoWithMessage 使用训练任务请求消息创建TrainingTaskInfo
func NewTrainingTaskInfoWithMessage(message *event.TrainTaskReqMessage) *TrainingTaskInfo {
	return &TrainingTaskInfo{
		TaskId:              message.TaskId,
		TrainTaskReqMessage: message,
	}
}
