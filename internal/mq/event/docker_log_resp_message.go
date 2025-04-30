package event

// DockerLogRespMessage 获取Docker logs消息响应
type DockerLogRespMessage struct {
	// 任务id
	TaskId int64 `json:"taskId"`

	// 日志数据
	Log string `json:"log"`

	// 任务类型
	TaskType int `json:"taskType"`
}
