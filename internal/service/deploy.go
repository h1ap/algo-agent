package service

import (
	pb "algo-agent/api/deploy/v1"
	"algo-agent/internal/biz"
	di "algo-agent/internal/model/deploy"
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

	deployServiceInfo := di.NewDeployServiceInfo(req)

	err := d.uc.Deploy(ctx, deployServiceInfo)
	if err != nil {
		d.log.WithContext(ctx).Errorf("Deploy failed: %v", err)
		return nil, err
	}

	return &pb.DeployReply{
		ServiceId:     deployServiceInfo.ServiceId,
		ServiceStatus: deployServiceInfo.ServiceStatus,
		Remark:        deployServiceInfo.Remark,
	}, nil
}

func (d *DeployServer) Destroy(ctx context.Context, req *pb.DestroyRequest) (*pb.DestroyReply, error) {
	jsonStr, _ := json.ToJSON(req)
	d.log.WithContext(ctx).Infof("Destroy: %v", jsonStr)

	err := d.uc.DestroyAndDelete(ctx, req.ServiceId)
	if err != nil {
		d.log.WithContext(ctx).Errorf("Destroy failed: %v", err)
		return nil, err
	}

	return &pb.DestroyReply{
		Message: "success",
	}, nil
}
