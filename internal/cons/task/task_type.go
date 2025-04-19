package tasks

// TaskType 表示任务类型
type TaskType struct {
	Code int
	Name string
}

var (
	// TRAIN 训练类型
	TRAIN = TaskType{
		Code: 0,
		Name: "训练",
	}

	// EVALUATE 评估类型
	EVALUATE = TaskType{
		Code: 1,
		Name: "评估",
	}
)

// 所有任务类型的集合，方便遍历
var AllTaskTypes = []TaskType{
	TRAIN,
	EVALUATE,
}

// GetCode 获取任务类型的代码
func (t TaskType) GetCode() int {
	return t.Code
}

// GetName 获取任务类型的名称
func (t TaskType) GetName() string {
	return t.Name
}
