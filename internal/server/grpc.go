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
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/go-kratos/kratos/v2/middleware/tracing"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// 设置全局trace
func initTracer(endpoint string) error {
	// 创建 exporter
	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// 将基于父span的采样率设置为100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// 始终确保在生产中批量处理
		tracesdk.WithBatcher(exporter),
		// 在资源中记录有关此应用程序的信息
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String("kratos-trace"),
			attribute.String("exporter", "otlp"),
			attribute.Float64("float", 312.23),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Server,
	s *service.OSSServer,
	d *service.DeployServer,
	ds *service.DockerServer,
	ts *service.TrainServer,
	es *service.EvalServer,
	exs *service.ExtractServer,
	logger log.Logger,
) *grpc.Server {
	ip := utils.GetLocalIP()
	// 格式化字符串
	endpoint := fmt.Sprintf("%s:%d", ip, 4318)
	err := initTracer(endpoint)
	if err != nil {
		panic(err)
	}
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(),
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
	ev1.RegisterEvalInfoServiceServer(srv, es)
	exv1.RegisterExtractInfoServiceServer(srv, exs)
	return srv
}
