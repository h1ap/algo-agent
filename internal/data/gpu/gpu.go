package gpu

import (
	"algo-agent/internal/cons/gpu"
	"algo-agent/internal/mq/event"
)

// GpuManager GPU管理器接口
type GpuManager interface {
	// GetGpuInfo 获取GPU信息
	GetGpuInfo() []event.GpuInfo

	// GetIdentify 获取GPU唯一标识
	GetIdentify() gpu.GpuVendor
}
