package service

import (
	v1 "algo-agent/api/extract/v1"
	"algo-agent/internal/biz"
	"algo-agent/internal/model/extract"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// ExtractServer 提取任务服务实现
type ExtractServer struct {
	v1.UnimplementedExtractInfoServiceServer

	uc  *biz.ExtractTaskUsecase
	log *log.Helper
}

// ResultInfo 处理提取任务结果
func (s *ExtractServer) ResultInfo(ctx context.Context, req *v1.ExtractTaskResultRequest) (*v1.ExtractResponse, error) {
	s.log.Infof("收到提取任务结果请求: taskId=%s, modelPath=%s", req.TaskId, req.ModelPath)

	result := &extract.ExtractTaskResult{
		TaskId:    req.TaskId,
		ModelPath: req.ModelPath,
	}

	err := s.uc.ExtractTaskResultHandle(ctx, result)
	if err != nil {
		s.log.Errorf("处理提取任务结果失败: %v", err)
		return &v1.ExtractResponse{
			Code:    500,
			Message: "处理提取任务结果失败: " + err.Error(),
		}, nil
	}

	return &v1.ExtractResponse{
		Code:    200,
		Message: "处理提取任务结果成功",
	}, nil
}
