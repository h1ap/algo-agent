package deploy

// DeployStatus 部署状态
type DeployStatus struct {
	Code int32
	Name string
}

var (
	// NOT_DEPLOYED 未部署
	NOT_DEPLOYED = DeployStatus{
		Code: 0,
		Name: "未部署",
	}

	// DEPLOYING 部署中
	DEPLOYING = DeployStatus{
		Code: 1,
		Name: "部署中",
	}

	// DOWNLOAD_DEPLOY_IMAGE 下载镜像
	DOWNLOAD_DEPLOY_IMAGE = DeployStatus{
		Code: 2,
		Name: "下载镜像",
	}

	// DOWNLOAD_ALGO_SCRIPTS 下载算法
	DOWNLOAD_ALGO_SCRIPTS = DeployStatus{
		Code: 3,
		Name: "下载算法",
	}

	// DOWNLOAD_WEIGHTS 下载权重文件
	DOWNLOAD_WEIGHTS = DeployStatus{
		Code: 4,
		Name: "下载权重文件",
	}

	// RUNNING 运行中
	RUNNING = DeployStatus{
		Code: 5,
		Name: "运行中",
	}

	// DESTROYED 已销毁
	DESTROYED = DeployStatus{
		Code: 6,
		Name: "已销毁",
	}

	// DEPLOYMENT_FAILED 部署失败
	DEPLOYMENT_FAILED = DeployStatus{
		Code: 9,
		Name: "部署失败",
	}
)

// 所有部署状态的集合，方便遍历
var AllDeployStatuses = []DeployStatus{
	NOT_DEPLOYED,
	DEPLOYING,
	DOWNLOAD_DEPLOY_IMAGE,
	DOWNLOAD_ALGO_SCRIPTS,
	DOWNLOAD_WEIGHTS,
	RUNNING,
	DESTROYED,
	DEPLOYMENT_FAILED,
}

// GetCode 获取部署状态的代码
func (d DeployStatus) GetCode() int32 {
	return d.Code
}

// GetName 获取部署状态的名称
func (d DeployStatus) GetName() string {
	return d.Name
}
