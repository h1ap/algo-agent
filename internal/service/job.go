package service

import (
	"algo-agent/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

// JobServer 定时任务管理器
type JobServer struct {
	gpuUC  *biz.GpuUsecase
	log    *log.Helper
	ctx    context.Context
	cancel context.CancelFunc
}

// Start 启动所有定时任务
func (jm *JobServer) Start() {
	jm.log.Info("正在启动定时任务管理器...")

	// 启动GPU系统指标上报定时任务
	jm.startGpuMetricsReporter()

	jm.log.Info("所有定时任务已启动")
}

// Stop 停止所有定时任务
func (jm *JobServer) Stop() {
	jm.log.Info("正在停止定时任务管理器...")

	// 取消上下文，让所有任务都能优雅退出
	jm.cancel()

	// 停止GPU系统指标上报定时任务
	jm.stopGpuMetricsReporter()

	jm.log.Info("所有定时任务已停止")
}

// startGpuMetricsReporter 启动GPU系统指标上报定时任务
func (jm *JobServer) startGpuMetricsReporter() {
	jm.log.Info("启动GPU系统指标上报定时任务")
	jm.gpuUC.Start()
}

// stopGpuMetricsReporter 停止GPU系统指标上报定时任务
func (jm *JobServer) stopGpuMetricsReporter() {
	jm.log.Info("停止GPU系统指标上报定时任务")
	jm.gpuUC.Stop()
}
