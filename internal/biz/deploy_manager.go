package biz

import (
	"algo-agent/internal/model/deploy"
	"context"
)

// DeployServiceManager 部署服务管理接口
type DeployServiceManager interface {
	// AddService 添加服务
	AddService(ctx context.Context, service *deploy.DeployServiceInfo) error

	// RemoveService 移除服务
	RemoveService(ctx context.Context, id string) bool

	// GetServiceList 获取服务列表
	GetServiceList(ctx context.Context) []*deploy.DeployServiceInfo

	// FindServiceById 根据ID查找服务
	FindServiceById(ctx context.Context, serviceId string) *deploy.DeployServiceInfo

	// UpdateService 更新服务
	UpdateService(ctx context.Context, updatedService *deploy.DeployServiceInfo) error

	// SaveToFile 保存到文件
	SaveToFile(ctx context.Context)

	// Stop 停止管理器
	Stop(ctx context.Context)
}
