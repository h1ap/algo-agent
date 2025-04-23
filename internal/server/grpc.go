package server

import (
	dv1 "algo-agent/api/deploy/v1"
	dcv1 "algo-agent/api/docker/v1"
	ov1 "algo-agent/api/oss/v1"
	tv1 "algo-agent/api/train/v1"
	"algo-agent/internal/conf"
	"algo-agent/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Server,
	s *service.OSSServer,
	d *service.DeployServer,
	ds *service.DockerServer,
	ts *service.TrainServer,
	logger log.Logger,
) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	ov1.RegisterOSSServiceServer(srv, s)
	dv1.RegisterDeployServiceServer(srv, d)
	dcv1.RegisterDockerServiceServer(srv, ds)
	tv1.RegisterTrainInfoServiceServer(srv, ts)
	return srv
}
