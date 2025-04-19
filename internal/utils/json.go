package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

// JSONUtil JSON工具类
type JSONUtil struct {
	// 自定义编码选项
	jsonEncoder *json.Encoder
	jsonDecoder *json.Decoder
}

var (
	// DefaultJSONUtil 默认的JSON工具实例
	DefaultJSONUtil = newJSONUtil()
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

// newJSONUtil 创建新的JSON工具实例
func newJSONUtil() *JSONUtil {
	return &JSONUtil{}
}

// ToJSON 将对象转换为JSON字符串
func (ju *JSONUtil) ToJSON(v interface{}) (string, error) {
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
func (ju *JSONUtil) ParseToMap(jsonStr string) (map[string]interface{}, error) {
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
func (ju *JSONUtil) Parse(jsonStr string) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("parse JSON failed: %v", err)
	}
	return result, nil
}

// ReadValue 将JSON字符串转换为指定类型
func (ju *JSONUtil) ReadValue(jsonStr string, v interface{}) error {
	return json.Unmarshal([]byte(jsonStr), v)
}

// IsCharSequence 判断是否为字符序列类型
func (ju *JSONUtil) IsCharSequence(v interface{}) bool {
	if v == nil {
		return false
	}

	kind := reflect.TypeOf(v).Kind()
	return kind == reflect.String
}

// 提供一些便捷的静态方法
func ToJSON(v interface{}) string {
	result, err := DefaultJSONUtil.ToJSON(v)
	if err != nil {
		log.Printf("Convert to JSON failed: %v", err)
		return ""
	}
	return result
}

func ParseToMap(jsonStr string) map[string]interface{} {
	result, err := DefaultJSONUtil.ParseToMap(jsonStr)
	if err != nil {
		log.Printf("Parse JSON to map failed: %v", err)
		return nil
	}
	return result
}

func Parse(jsonStr string) interface{} {
	result, err := DefaultJSONUtil.Parse(jsonStr)
	if err != nil {
		log.Printf("Parse JSON failed: %v", err)
		return nil
	}
	return result
}

// ReadValueTo 泛型方法，将JSON字符串转换为指定类型
func ReadValueTo[T any](jsonStr string) (*T, error) {
	var result T
	err := DefaultJSONUtil.ReadValue(jsonStr, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
