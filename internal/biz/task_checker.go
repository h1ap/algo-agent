package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// TaskCheckerUsecase 任务检查器用例
type TaskCheckerUsecase struct {
	trainingUsecase *TrainingTaskUsecase
	evalUsecase     *EvalTaskUsecase
	deployUsecase   *DeployUsecase
	extractUsecase  *ExtractTaskUsecase

	ctx            context.Context
	cancel         context.CancelFunc
	trainingTicker *time.Ticker
	evalTicker     *time.Ticker
	deployTicker   *time.Ticker
	extractTicker  *time.Ticker

	log *log.Helper
}

// Start 启动所有定时检查任务
func (uc *TaskCheckerUsecase) Start() {
	// 训练任务检查 - 每30秒检查一次
	uc.trainingTicker = time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-uc.trainingTicker.C:
				uc.checkTrainingTask()
			case <-uc.ctx.Done():
				return
			}
		}
	}()

	// 评估任务检查 - 每30秒检查一次
	uc.evalTicker = time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-uc.evalTicker.C:
				uc.checkEvalTask()
			case <-uc.ctx.Done():
				return
			}
		}
	}()

	// 推理服务检查 - 每30秒检查一次
	uc.deployTicker = time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-uc.deployTicker.C:
				uc.checkDeployService()
			case <-uc.ctx.Done():
				return
			}
		}
	}()

	// 提取任务检查 - 每2分钟检查一次
	uc.extractTicker = time.NewTicker(2 * time.Minute)
	go func() {
		for {
			select {
			case <-uc.extractTicker.C:
				uc.checkExtractTask()
			case <-uc.ctx.Done():
				return
			}
		}
	}()

	uc.log.Info("所有任务定时检查服务已启动")
}

// Stop 停止所有定时检查任务
func (uc *TaskCheckerUsecase) Stop() {
	if uc.trainingTicker != nil {
		uc.trainingTicker.Stop()
	}
	if uc.evalTicker != nil {
		uc.evalTicker.Stop()
	}
	if uc.deployTicker != nil {
		uc.deployTicker.Stop()
	}
	if uc.extractTicker != nil {
		uc.extractTicker.Stop()
	}
	uc.cancel()
	uc.log.Info("所有任务定时检查服务已停止")
}

// 检查训练任务
func (uc *TaskCheckerUsecase) checkTrainingTask() {
	defer func() {
		if r := recover(); r != nil {
			uc.log.Errorf("检查训练任务时发生异常: %v", r)
		}
	}()

	uc.log.Debug("开始检查训练任务状态")
	uc.trainingUsecase.CheckTask(uc.ctx)
}

// 检查评估任务
func (uc *TaskCheckerUsecase) checkEvalTask() {
	defer func() {
		if r := recover(); r != nil {
			uc.log.Errorf("检查评估任务时发生异常: %v", r)
		}
	}()

	uc.log.Debug("开始检查评估任务状态")
	uc.evalUsecase.CheckTask(uc.ctx)
}

// 检查推理服务
func (uc *TaskCheckerUsecase) checkDeployService() {
	defer func() {
		if r := recover(); r != nil {
			uc.log.Errorf("检查推理服务时发生异常: %v", r)
		}
	}()

	uc.log.Debug("开始检查推理服务状态")
	uc.deployUsecase.CheckTask(uc.ctx)
}

// 检查提取任务
func (uc *TaskCheckerUsecase) checkExtractTask() {
	defer func() {
		if r := recover(); r != nil {
			uc.log.Errorf("检查提取任务时发生异常: %v", r)
		}
	}()

	uc.log.Debug("开始检查提取任务状态")
	uc.extractUsecase.CheckTask(uc.ctx)
}
