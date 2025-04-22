package train

import (
	"encoding/json"
)

// TrainingEpochInfo 训练任务每个epoch的返回值
type TrainingEpochInfo struct {
	// TaskId 任务ID
	TaskId string `json:"taskId"`

	// Epoch 当前训练轮次
	Epoch int32 `json:"epoch"`

	// EstimatedTimeLeft 预估剩余时间(秒)
	EstimatedTimeLeft int64 `json:"estimatedTimeLeft"`

	// DynamicFields 动态指标
	DynamicFields map[string]interface{} `json:"-"`
}

// MarshalJSON 自定义JSON序列化方法，将动态字段合并到输出
func (t TrainingEpochInfo) MarshalJSON() ([]byte, error) {
	type Alias TrainingEpochInfo
	return json.Marshal(struct {
		Alias
		*mapstruc
	}{
		Alias:    Alias(t),
		mapstruc: (*mapstruc)(&t.DynamicFields),
	})
}

// UnmarshalJSON 自定义JSON反序列化方法，处理动态字段
func (t *TrainingEpochInfo) UnmarshalJSON(data []byte) error {
	type Alias TrainingEpochInfo
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if t.DynamicFields == nil {
		t.DynamicFields = make(map[string]interface{})
	}

	// 解析所有额外字段
	var objMap map[string]*json.RawMessage
	if err := json.Unmarshal(data, &objMap); err != nil {
		return err
	}

	for k, v := range objMap {
		switch k {
		case "taskId", "epoch", "estimatedTimeLeft":
			// 跳过已处理的字段
			continue
		default:
			var value interface{}
			if err := json.Unmarshal(*v, &value); err != nil {
				return err
			}
			t.DynamicFields[k] = value
		}
	}

	return nil
}

// mapstruc 用于帮助序列化map
type mapstruc map[string]interface{}

// MarshalJSON 实现json.Marshaler接口
func (m *mapstruc) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("{}"), nil
	}
	return json.Marshal((map[string]interface{})(*m))
}
