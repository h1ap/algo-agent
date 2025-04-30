package task

// TaskType 表示任务类型
type TaskType struct {
	code int
	name string
}

var (
	// TRAIN 训练类型
	TRAIN = TaskType{
		code: 0,
		name: "训练",
	}

	// EVALUATE 评估类型
	EVALUATE = TaskType{
		code: 1,
		name: "评估",
	}
)

// 所有任务类型的集合，方便遍历
var AllTaskTypes = []TaskType{
	TRAIN,
	EVALUATE,
}

// Code 获取编码
func (e *TaskType) Code() int {
	return e.code
}

// Name 获取枚举名
func (e *TaskType) Name() string {
	return e.name
}
