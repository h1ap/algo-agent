package eval

// EvalType 评估类型枚举
type EvalType string

const (
	// EvalTypeModel 模型
	EvalTypeModel EvalType = "1"
	// EvalTypeCheckpoint 检查点
	EvalTypeCheckpoint EvalType = "2"
)

// Code 返回评估类型的编码
func (e EvalType) Code() string {
	return string(e)
}
