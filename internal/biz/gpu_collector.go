package biz

import (
	"algo-agent/internal/cons/gpu"
	"algo-agent/internal/cons/mq"
	"algo-agent/internal/mq/event"
	"context"
	"runtime"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// GpuManager GPU管理器接口
type GpuManager interface {
	// GetGpuInfo 获取GPU信息
	GetGpuInfo() []event.GpuInfo

	// GetIdentify 获取GPU唯一标识
	GetIdentify() gpu.GpuVendor
}

type GpuUsecase struct {
	g  GpuManager
	mq MqService

	tsn       string
	nodeName  string
	ipAddress string

	log *log.Helper
}

// GetSystemMetrics 获取系统指标
func (uc *GpuUsecase) GetSystemMetrics() *event.SystemMetricsEvent {
	// 获取CPU信息
	cpuPercent, err := cpu.Percent(time.Second, false)
	cpuLoad := 0.0
	if err == nil && len(cpuPercent) > 0 {
		cpuLoad = cpuPercent[0]
	}

	physicalCoreCount, err := cpu.Counts(false)
	if err != nil {
		physicalCoreCount = runtime.NumCPU()
	}

	logicalCoreCount, err := cpu.Counts(true)
	if err != nil {
		logicalCoreCount = runtime.NumCPU()
	}

	// 获取内存信息
	var totalMemory, freeMemory, usedMemory int64
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		totalMemory = int64(memInfo.Total)
		freeMemory = int64(memInfo.Free)
		usedMemory = int64(memInfo.Used)
	}

	// 获取GPU信息
	gpuList := uc.g.GetGpuInfo()

	// 创建系统指标事件
	return event.NewSystemMetricsEvent(
		uc.nodeName,
		uc.ipAddress,
		int64(physicalCoreCount),
		int64(logicalCoreCount),
		cpuLoad,
		totalMemory,
		freeMemory,
		usedMemory,
		gpuList,
	)
}

// ReportSystemMetrics 上报系统指标
func (uc *GpuUsecase) ReportSystemMetrics(ctx context.Context) {
	metrics := uc.GetSystemMetrics()
	// 发送到训练服务
	err := uc.mq.SendToService(ctx, uc.tsn, &event.ReqMessage{
		Type:    mq.SYSTEM_METRICS.Code(),
		Payload: metrics,
	})
	if err != nil {
		uc.log.Errorf("系统指标上报失败: %v", err)
		return
	}
	uc.log.Debugf("系统指标已上报, CPU: %.2f%%, 内存: %d/%d, GPU数量: %d",
		metrics.CpuLoad, metrics.UsedMemorySize, metrics.TotalMemorySize, len(metrics.GpuList))
}
