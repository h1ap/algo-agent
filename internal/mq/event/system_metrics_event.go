package event

import "time"

// GpuInfo Nvidia Gpu 信息
type GpuInfo struct {
	// 显卡序号
	Index int `json:"index"`

	// 设备唯一标识
	DeviceId string `json:"deviceId"`

	// 显卡名称
	Name string `json:"name"`

	// 显卡供应商
	Vendor string `json:"vendor"`

	// 显存使用率
	CudaUsage float64 `json:"cudaUsage"`

	// 显存总大小
	MemoryTotal int64 `json:"memoryTotal"`

	// 显存已用大小
	MemoryUsed int64 `json:"memoryUsed"`
}

// NewGpuInfo 创建一个新的GpuInfo实例
func NewGpuInfo(index int, deviceId, name, vendor string, cudaUsage float64, memoryTotal, memoryUsed int64) *GpuInfo {
	return &GpuInfo{
		Index:       index,
		DeviceId:    deviceId,
		Name:        name,
		Vendor:      vendor,
		CudaUsage:   cudaUsage,
		MemoryTotal: memoryTotal,
		MemoryUsed:  memoryUsed,
	}
}

// NewEmptyGpuInfo 创建一个空的GpuInfo实例
func NewEmptyGpuInfo() *GpuInfo {
	return &GpuInfo{}
}

// SystemMetricsEvent 服务器系统资源使用情况
type SystemMetricsEvent struct {
	// 从配置文件读取
	NodeName string `json:"nodeName"`

	// ip地址
	IpAddress string `json:"ipAddress"`

	// 物理核心数
	PhysicalCpuCoreNum int64 `json:"physicalCpuCoreNum"`

	// 逻辑核心数
	CpuCoreNum int64 `json:"cpuCoreNum"`

	// CPU使用率
	CpuLoad float64 `json:"cpuLoad"`

	// 总内存
	TotalMemorySize int64 `json:"totalMemorySize"`

	// 已用内存
	FreeMemorySize int64 `json:"freeMemorySize"`

	// 可用内存
	UsedMemorySize int64 `json:"usedMemorySize"`

	// Gpu信息
	GpuList []GpuInfo `json:"gpuList"`

	// 时间戳
	Timestamp int64 `json:"timestamp"`
}

// NewSystemMetricsEvent 创建一个新的SystemMetricsEvent实例
func NewSystemMetricsEvent(
	nodeName, ipAddress string,
	physicalCpuCoreNum, cpuCoreNum int64,
	cpuLoad float64,
	totalMemorySize, freeMemorySize, usedMemorySize int64,
	gpuList []GpuInfo,
) *SystemMetricsEvent {
	return &SystemMetricsEvent{
		NodeName:           nodeName,
		IpAddress:          ipAddress,
		PhysicalCpuCoreNum: physicalCpuCoreNum,
		CpuCoreNum:         cpuCoreNum,
		CpuLoad:            cpuLoad,
		TotalMemorySize:    totalMemorySize,
		FreeMemorySize:     freeMemorySize,
		UsedMemorySize:     usedMemorySize,
		GpuList:            gpuList,
		Timestamp:          time.Now().Unix(),
	}
}

// NewEmptySystemMetricsEvent 创建一个空的SystemMetricsEvent实例
func NewEmptySystemMetricsEvent() *SystemMetricsEvent {
	return &SystemMetricsEvent{
		GpuList: make([]GpuInfo, 0),
	}
}
