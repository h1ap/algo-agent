package deploy

// 部署指令
type DeployOp struct {
	Code int32
	Name string
}

var (
	// DEPLOY 部署
	DEPLOY = DeployOp{
		Code: 0,
		Name: "部署",
	}

	// DESTROY 销毁
	DESTROY = DeployOp{
		Code: 1,
		Name: "销毁",
	}
)

// 所有部署指令的集合，方便遍历
var AllDeployOps = []DeployOp{
	DEPLOY,
	DESTROY,
}

// GetCode 获取部署指令的代码
func (d DeployOp) GetCode() int32 {
	return d.Code
}

// GetName 获取部署指令的名称
func (d DeployOp) GetName() string {
	return d.Name
}
