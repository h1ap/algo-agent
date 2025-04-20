package service

import (
	pb "algo-agent/api/deploy/v1"
	"algo-agent/internal/biz"
	json "algo-agent/internal/utils"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type DeployServer struct {
	pb.UnimplementedDeployServiceServer
	uc  *biz.DeployUsecase
	log *log.Helper
}

func (s *DeployServer) Deploy(ctx context.Context, req *pb.DeployRequest) (*pb.DeployRequest, error) {
	jsonStr, _ := json.ToJSON(req)
	s.log.WithContext(ctx).Infof("Deploy: %v", jsonStr)

	return &pb.DeployRequest{}, nil
}
