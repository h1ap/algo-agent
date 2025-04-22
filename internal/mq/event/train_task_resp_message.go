package event

import (
	"encoding/json"
	"time"
)

// TrainTaskRespMessage 训练任务响应消息
type TrainTaskRespMessage struct {
	// 训练任务ID
	TaskId string `json:"taskId"`

	// 当前任务状态
	// see cons/train/train_job_status_enum.go
	TaskStatus int32 `json:"taskStatus"`

	// 训练指标信息
	Metrics *Metrics `json:"metrics"`

	// 训练轮次
	Epoch int32 `json:"epoch"`

	// 是否为检查点
	IsCheckpoint bool `json:"isCheckpoint"`

	// 检查点文件路径
	CheckpointFilePath string `json:"checkpointFilePath"`

	// 最好的权重(训练结束)
	BestWeightPath string `json:"bestWeightPath"`

	// 最后的权重(训练结束)
	LastWeightPath string `json:"lastWeightPath"`

	// 备注
	Remark string `json:"remark"`
}

// Metrics 训练指标详细信息
type Metrics struct {
	// 训练轮次
	Epoch int32 `json:"epoch"`

	// 批次大小
	BatchSize int32 `json:"batchSize"`

	// 创建时间
	CreateTime time.Time `json:"createTime"`

	// 剩余预估时间(s)
	EstimateTimeLeft int64 `json:"estimateTimeLeft"`

	// 动态指标
	DynamicFields map[string]interface{} `json:"-"`
}

// MarshalJSON 用于JSON序列化时将DynamicFields的字段直接放在外层
func (m *Metrics) MarshalJSON() ([]byte, error) {
	type Alias Metrics

	// 创建一个包含所有字段的映射
	fields := make(map[string]interface{})

	// 添加结构体字段
	data, err := json.Marshal((*Alias)(m))
	if err != nil {
		return nil, err
	}

	// 解析原始结构体字段到map
	err = json.Unmarshal(data, &fields)
	if err != nil {
		return nil, err
	}

	// 添加动态字段
	for k, v := range m.DynamicFields {
		fields[k] = v
	}

	// 将合并后的map序列化
	return json.Marshal(fields)
}

// UnmarshalJSON 用于JSON反序列化时将外层未知字段放入DynamicFields
func (m *Metrics) UnmarshalJSON(data []byte) error {
	type Alias Metrics

	// 创建别名对象
	alias := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	// 解码基本字段
	if err := json.Unmarshal(data, alias); err != nil {
		return err
	}

	// 解码所有字段到map
	var fields map[string]interface{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	// 初始化DynamicFields
	m.DynamicFields = make(map[string]interface{})

	// 获取结构体中已定义的字段
	knownFields := map[string]bool{
		"epoch":            true,
		"batchSize":        true,
		"createTime":       true,
		"estimateTimeLeft": true,
	}

	// 添加未知字段到DynamicFields
	for k, v := range fields {
		if !knownFields[k] {
			m.DynamicFields[k] = v
		}
	}

	return nil
}
