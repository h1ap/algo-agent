package server

import (
	dv1 "algo-agent/api/deploy/v1"
	dcv1 "algo-agent/api/docker/v1"
	ev1 "algo-agent/api/eval/v1"
	exv1 "algo-agent/api/extract/v1"
	ov1 "algo-agent/api/oss/v1"
	tv1 "algo-agent/api/train/v1"
	"algo-agent/internal/conf"
	"algo-agent/internal/service"
	"algo-agent/internal/utils"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c *conf.Server,
	s *service.OSSServer,
	d *service.DeployServer,
	ds *service.DockerServer,
	ts *service.TrainServer,
	es *service.EvalServer,
	exs *service.ExtractServer,
	logger log.Logger,
) *http.Server {
	ip := utils.GetLocalIP()
	// 格式化字符串
	endpoint := fmt.Sprintf("%s:%d", ip, 4318)
	err := initTracer(endpoint)
	if err != nil {
		panic(err)
	}
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	ov1.RegisterOSSServiceHTTPServer(srv, s)
	dv1.RegisterDeployServiceHTTPServer(srv, d)
	dcv1.RegisterDockerServiceHTTPServer(srv, ds)
	tv1.RegisterTrainInfoServiceHTTPServer(srv, ts)
	ev1.RegisterEvalInfoServiceHTTPServer(srv, es)
	exv1.RegisterExtractInfoServiceHTTPServer(srv, exs)
	return srv
}
