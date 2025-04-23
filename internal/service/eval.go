package service

import (
	"context"

	v1 "algo-agent/api/eval/v1"
	"algo-agent/internal/biz"
	"algo-agent/internal/model/eval"

	"github.com/go-kratos/kratos/v2/log"
)

// EvalServer 评估服务实现
type EvalServer struct {
	v1.UnimplementedEvalInfoServiceServer

	uc  *biz.EvalTaskUsecase
	log *log.Helper
}

// BatchInfo 处理评估批次信息
func (s *EvalServer) BatchInfo(ctx context.Context, req *v1.EvalBatchInfoRequest) (*v1.EvalResponse, error) {
	s.log.Infof("收到批次信息请求: taskId=%s, 详情数量=%d", req.TaskId, len(req.Details))

	// 转换请求为内部模型
	batchInfo := &eval.EvalBatchInfo{
		TaskId:  req.TaskId,
		Details: make([]eval.EvalDetail, 0, len(req.Details)),
	}
	for _, detail := range req.Details {
		batchInfo.Details = append(batchInfo.Details, eval.EvalDetail{
			DataUuid: detail.DataUuid,
			EvalData: detail.EvalData,
		})
	}

	// 调用业务逻辑处理批次信息
	go s.uc.BatchInfoHandle(ctx, batchInfo)

	return &v1.EvalResponse{
		Code:    0,
		Message: "成功",
	}, nil
}

// EpochInfo 处理评估周期信息（已废弃，使用BatchInfo）
func (s *EvalServer) EpochInfo(ctx context.Context, req *v1.EvalBatchInfoRequest) (*v1.EvalResponse, error) {
	s.log.Infof("收到周期信息请求(已废弃): taskId=%s", req.TaskId)
	return s.BatchInfo(ctx, req)
}

// FinishInfo 处理评估完成信息
func (s *EvalServer) FinishInfo(ctx context.Context, req *v1.EvalTaskResultRequest) (*v1.EvalResponse, error) {
	s.log.Infof("收到完成信息请求: taskId=%s", req.TaskId)

	// 转换请求为内部模型
	resultInfo := &eval.EvalTaskResult{
		TaskId:          req.TaskId,
		OverallMetrics:  make(map[string]interface{}),
		ClassifyMetrics: make([]map[string]interface{}, 0, len(req.ClassifyMetrics)),
	}

	// 转换整体指标
	for k, v := range req.OverallMetrics {
		resultInfo.OverallMetrics[k] = v
	}

	// 转换分类指标
	for _, metric := range req.ClassifyMetrics {
		metricMap := make(map[string]interface{})
		for k, v := range metric.Metrics {
			metricMap[k] = v
		}
		resultInfo.ClassifyMetrics = append(resultInfo.ClassifyMetrics, metricMap)
	}

	// 调用业务逻辑处理完成信息
	go s.uc.FinishHandle(ctx, resultInfo)

	return &v1.EvalResponse{
		Code:    0,
		Message: "成功",
	}, nil
}
