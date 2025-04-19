package model

import (
	v1 "algo-agent/api/deploy/v1"
)

// DeployServiceInfo 部署信息
type DeployServiceInfo struct {
	// ServiceID 任务id
	ServiceID string `json:"serviceId"`

	// ServiceContainerName 训练容器名称（唯一），不需要传入
	ServiceContainerName string `json:"serviceContainerName"`

	// ServiceStatus 当前任务状态
	// 参考 DeployStatusEnum
	ServiceStatus int `json:"serviceStatus"`

	// Remark 备注
	Remark string `json:"remark"`

	// DeployRequest 部署请求
	DeployRequest *v1.DeployRequest `json:"deployRequest"`
}

// NewDeployServiceInfo 创建新的DeployServiceInfo实例
func NewDeployServiceInfo(deployRequest *v1.DeployRequest) *DeployServiceInfo {
	return &DeployServiceInfo{
		ServiceID:     deployRequest.GetServiceId(),
		DeployRequest: deployRequest,
	}
}
