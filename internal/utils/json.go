package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// 自定义时间格式
const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
)

// 自定义时间类型
type CustomTime time.Time

// MarshalJSON 实现自定义时间序列化
func (t CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(t).Format(DateTimeFormat))), nil
}

// UnmarshalJSON 实现自定义时间反序列化
func (t *CustomTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = strings.Trim(str, `"`)

	parsed, err := time.Parse(DateTimeFormat, str)
	if err != nil {
		return err
	}
	*t = CustomTime(parsed)
	return nil
}

// ToJSON 将对象转换为JSON字符串
func ToJSON(v interface{}) (string, error) {
	// 如果是字符串类型，直接返回
	if str, ok := v.(string); ok {
		return str, nil
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("marshal to JSON failed: %v", err)
	}
	return string(bytes), nil
}

// ParseToMap 将JSON字符串转换为map
func ParseToMap(jsonStr string) (map[string]interface{}, error) {
	if strings.TrimSpace(jsonStr) == "" {
		return nil, nil
	}

	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("parse JSON to map failed: %v", err)
	}
	return result, nil
}

// Parse 将JSON字符串转换为interface{}
func Parse(jsonStr string) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("parse JSON failed: %v", err)
	}
	return result, nil
}

// ReadValue 将JSON字符串转换为指定类型
func ReadValue(jsonStr string, v interface{}) error {
	return json.Unmarshal([]byte(jsonStr), v)
}

// IsCharSequence 判断是否为字符序列类型
func IsCharSequence(v interface{}) bool {
	if v == nil {
		return false
	}

	kind := reflect.TypeOf(v).Kind()
	return kind == reflect.String
}
