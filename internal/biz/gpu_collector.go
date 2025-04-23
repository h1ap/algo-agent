package biz

import (
	"algo-agent/internal/cons/gpu"
	"algo-agent/internal/mq/event"
	"algo-agent/internal/utils"
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

	tsn            string
	nodeName       string
	ipAddress      string
	ctx            context.Context
	cancel         context.CancelFunc
	ticker         *time.Ticker
	reportInterval time.Duration

	log *log.Helper
}

// Start 启动定时任务
func (uc *GpuUsecase) Start() {
	uc.ticker = time.NewTicker(uc.reportInterval)
	go func() {
		for {
			select {
			case <-uc.ticker.C:
				uc.reportSystemMetrics()
			case <-uc.ctx.Done():
				return
			}
		}
	}()
	uc.log.Info("系统指标定时上报服务已启动")
}

// Stop 停止定时任务
func (uc *GpuUsecase) Stop() {
	if uc.ticker != nil {
		uc.ticker.Stop()
	}
	uc.cancel()
	uc.log.Info("系统指标定时上报服务已停止")
}

// SetReportInterval 设置上报间隔
func (uc *GpuUsecase) SetReportInterval(interval time.Duration) {
	uc.reportInterval = interval
	if uc.ticker != nil {
		uc.ticker.Reset(interval)
	}
	uc.log.Infof("系统指标上报间隔已设置为 %v", interval)
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

// reportSystemMetrics 上报系统指标
func (uc *GpuUsecase) reportSystemMetrics() {
	metrics := uc.GetSystemMetrics()
	// 将指标转换为JSON
	jsonStr, err := utils.ToJSON(metrics)
	if err != nil {
		uc.log.Errorf("系统指标转换JSON失败: %v", err)
		return
	}

	// 发送到训练服务
	err = uc.mq.SendToService(uc.ctx, uc.tsn, jsonStr)
	if err != nil {
		uc.log.Errorf("系统指标上报失败: %v", err)
		return
	}
	uc.log.Infof("系统指标已上报, CPU: %.2f%%, 内存: %d/%d, GPU数量: %d",
		metrics.CpuLoad, metrics.UsedMemorySize, metrics.TotalMemorySize, len(metrics.GpuList))
}
