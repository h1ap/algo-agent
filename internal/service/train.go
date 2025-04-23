package service

import (
	"context"

	pb "algo-agent/api/train/v1"
	"algo-agent/internal/biz"
	"algo-agent/internal/model/train"

	"github.com/go-kratos/kratos/v2/log"
)

// TrainServer 实现训练信息服务
type TrainServer struct {
	pb.UnimplementedTrainInfoServiceServer
	ttu *biz.TrainingTaskUsecase
	log *log.Helper
}

// EpochInfo 处理训练周期信息
func (s *TrainServer) EpochInfo(ctx context.Context, req *pb.TrainingEpochInfoRequest) (*pb.TrainingResponse, error) {
	s.log.Infof("收到训练周期信息，任务ID: %s，周期: %d", req.TaskId, req.Epoch)

	// 转换动态字段
	dynamicFields := make(map[string]interface{})
	for k, v := range req.DynamicFields {
		dynamicFields[k] = v
	}

	// 构建内部模型
	epochInfo := &train.TrainingEpochInfo{
		TaskId:            req.TaskId,
		Epoch:             req.Epoch,
		EstimatedTimeLeft: req.EstimatedTimeLeft,
		DynamicFields:     dynamicFields,
	}

	// 处理周期信息
	s.ttu.EpochInfoHandle(ctx, epochInfo)

	return &pb.TrainingResponse{
		Code:    0,
		Message: "success",
	}, nil
}

// CheckpointInfo 处理检查点信息
func (s *TrainServer) CheckpointInfo(ctx context.Context, req *pb.TrainingCheckpointRequest) (*pb.TrainingResponse, error) {
	s.log.Infof("收到检查点信息，任务ID: %s，周期: %d", req.TaskId, req.Epoch)

	// 构建内部模型
	checkpoint := &train.TrainingCheckpoint{
		TaskId:         req.TaskId,
		Epoch:          req.Epoch,
		CheckpointPath: req.CheckpointPath,
	}

	// 处理检查点信息
	err := s.ttu.CheckpointHandle(ctx, checkpoint)
	if err != nil {
		s.log.Errorf("处理检查点信息失败: %v", err)
		return &pb.TrainingResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &pb.TrainingResponse{
		Code:    0,
		Message: "success",
	}, nil
}

// FinishInfo 处理训练完成信息
func (s *TrainServer) FinishInfo(ctx context.Context, req *pb.TrainingTaskResultRequest) (*pb.TrainingResponse, error) {
	s.log.Infof("收到训练完成信息，任务ID: %s，最优周期: %d", req.TaskId, req.BestEpoch)

	// 构建内部模型
	result := &train.TrainingTaskResult{
		TaskId:         req.TaskId,
		BestEpoch:      req.BestEpoch,
		BestModelPath:  req.BestModelPath,
		FinalModelPath: req.FinalModelPath,
	}

	// 处理训练完成信息
	err := s.ttu.FinishHandle(ctx, result)
	if err != nil {
		s.log.Errorf("处理训练完成信息失败: %v", err)
		return &pb.TrainingResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &pb.TrainingResponse{
		Code:    0,
		Message: "success",
	}, nil
}
