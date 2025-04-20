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

func (d *DeployServer) Deploy(ctx context.Context, req *pb.DeployRequest) (*pb.DeployReply, error) {
	jsonStr, _ := json.ToJSON(req)
	d.log.WithContext(ctx).Infof("Deploy: %v", jsonStr)

	err := d.uc.Deploy(ctx, req)
	if err != nil {
		d.log.WithContext(ctx).Errorf("Deploy failed: %v", err)
		return nil, err
	}

	return &pb.DeployReply{
		Code: 200,
		Msg:  "success",
	}, nil
}
