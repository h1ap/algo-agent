package train

// TrainJobOp （训练/评估）作业操作指令
type TrainJobOp struct {
	Code int
	Name string
}

var (
	// TrainJobOpStart 启动
	TrainJobOpStart = TrainJobOp{
		Code: 0,
		Name: "启动",
	}

	// TrainJobOpStop 终止
	TrainJobOpStop = TrainJobOp{
		Code: 1,
		Name: "终止",
	}

	// 注释掉的枚举值：
	// TrainJobOpResume = TrainJobOp{
	//     Code: 2,
	//     Name: "恢复",
	// }
	//
	// TrainJobOpPublish = TrainJobOp{
	//     Code: 3,
	//     Name: "发布",
	// }
)

// GetCode 获取操作编码
func (t TrainJobOp) GetCode() int {
	return t.Code
}

// GetName 获取操作名称
func (t TrainJobOp) GetName() string {
	return t.Name
}
