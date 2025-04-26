package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

// TaskCheckerUsecase 任务检查器用例
type TaskCheckerUsecase struct {
	trainingUsecase *TrainingTaskUsecase
	evalUsecase     *EvalTaskUsecase
	deployUsecase   *DeployUsecase
	extractUsecase  *ExtractTaskUsecase

	log *log.Helper
}

// 检查训练任务
func (uc *TaskCheckerUsecase) CheckTrainingTask(ctx context.Context) {
	uc.log.Debug("检查训练任务状态开始...")
	uc.trainingUsecase.CheckTask(ctx)
	uc.log.Debug("检查训练任务状态结束...")
}

// 检查评估任务
func (uc *TaskCheckerUsecase) CheckEvalTask(ctx context.Context) {
	uc.log.Debug("检查评估任务状态开始...")
	uc.evalUsecase.CheckTask(ctx)
	uc.log.Debug("检查评估任务状态结束...")
}

// 检查推理服务
func (uc *TaskCheckerUsecase) CheckDeployService(ctx context.Context) {
	uc.log.Debug("检查推理服务状态开始...")
	uc.deployUsecase.CheckTask(ctx)
	uc.log.Debug("检查推理服务状态结束...")
}

// 检查提取任务
func (uc *TaskCheckerUsecase) CheckExtractTask(ctx context.Context) {
	uc.log.Debug("检查提取任务状态开始...")
	uc.extractUsecase.CheckTask(ctx)
	uc.log.Debug("检查提取任务状态结束...")
}
