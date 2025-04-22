package extract

import (
	"algo-agent/internal/mq/event"
)

// ExtractTaskInfo 提取任务信息
type ExtractTaskInfo struct {
	// TaskId 任务id
	TaskId string `json:"taskId"`

	// ContainerName 训练容器名称（唯一），不需要传入
	ContainerName string `json:"containerName"`

	// TaskStatus 当前任务状态
	TaskStatus int32 `json:"taskStatus"`

	// Remark 任务备注
	Remark string `json:"remark"`

	// TrainPublishReqMessage 提取任务参数信息
	TrainPublishReqMessage *event.TrainPublishReqMessage `json:"trainPublishReqMessage"`
}

// NewExtractTaskInfo 创建ExtractTaskInfo的构造函数
func NewExtractTaskInfo() *ExtractTaskInfo {
	return &ExtractTaskInfo{}
}

// NewExtractTaskInfoWithMessage 使用训练发布请求消息创建ExtractTaskInfo
func NewExtractTaskInfoWithMessage(message *event.TrainPublishReqMessage) *ExtractTaskInfo {
	return &ExtractTaskInfo{
		TaskId:                 message.TrainDetailId,
		TrainPublishReqMessage: message,
	}
}
