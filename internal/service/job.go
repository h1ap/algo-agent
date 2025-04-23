package service

import (
	"algo-agent/internal/biz"
	"algo-agent/internal/utils"
	"context"
	"net/url"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
)

// JobServer 定时任务管理器
type JobServer struct {
	gpuUC       *biz.GpuUsecase
	taskChecker *biz.TaskCheckerUsecase
	log         *log.Helper
	ctx         context.Context
	cancel      context.CancelFunc
}

// Start 启动所有定时任务
func (jm *JobServer) Start(context.Context) error {
	jm.log.Info("正在启动定时任务管理器...")

	// 启动GPU系统指标上报定时任务
	jm.startGpuMetricsReporter()

	// 启动任务检查器
	jm.startTaskChecker()

	jm.log.Info("所有定时任务已启动")
	return nil
}

// Stop 停止所有定时任务
func (jm *JobServer) Stop(context.Context) error {
	jm.log.Info("正在停止定时任务管理器...")

	// 取消上下文，让所有任务都能优雅退出
	jm.cancel()

	// 停止GPU系统指标上报定时任务
	jm.stopGpuMetricsReporter()

	// 停止任务检查器
	jm.stopTaskChecker()

	jm.log.Info("所有定时任务已停止")
	return nil
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

// startTaskChecker 启动任务检查器
func (jm *JobServer) startTaskChecker() {
	jm.log.Info("启动任务检查定时任务")
	jm.taskChecker.Start()
}

// stopTaskChecker 停止任务检查器
func (jm *JobServer) stopTaskChecker() {
	jm.log.Info("停止任务检查定时任务")
	jm.taskChecker.Stop()
}

func (j *JobServer) Endpoint() (*url.URL, error) {
	ip := utils.GetLocalIP()
	u := &url.URL{
		Scheme: "http",
		Host:   ip + ":" + strconv.Itoa(int(8001)),
	}
	return u, nil
}
