package service

import (
	"algo-agent/internal/biz"
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"runtime"
)

type CronServer struct {
	tcu *biz.TaskCheckerUsecase
}

func NewCronServer(tcu *biz.TaskCheckerUsecase) *CronServer {
	return &CronServer{
		tcu: tcu,
	}
}

func (s *CronServer) RunCheckTrainingTaskService(ctx context.Context) error {
	ctx, span := otel.Tracer("Service").Start(ctx, "RunCheckTrainingTaskService", trace.WithSpanKind(trace.SpanKindProducer))
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 2048)
			runtime.Stack(buf, false)
			span.SetStatus(codes.Error, fmt.Sprintf("%v", err))
		}
		span.End()
	}()

	s.tcu.CheckTrainingTask(ctx)
	return nil
}

func (s *CronServer) RunCheckEvalTaskService(ctx context.Context) error {
	ctx, span := otel.Tracer("Service").Start(ctx, "RunCheckEvalTaskService", trace.WithSpanKind(trace.SpanKindProducer))
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 2048)
			runtime.Stack(buf, false)
			span.SetStatus(codes.Error, fmt.Sprintf("%v", err))
		}
		span.End()
	}()

	s.tcu.CheckEvalTask(ctx)
	return nil
}

func (s *CronServer) RunCheckDeployService(ctx context.Context) error {
	ctx, span := otel.Tracer("Service").Start(ctx, "RunCheckDeployService", trace.WithSpanKind(trace.SpanKindProducer))
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 2048)
			runtime.Stack(buf, false)
			span.SetStatus(codes.Error, fmt.Sprintf("%v", err))
		}
		span.End()
	}()

	s.tcu.CheckDeployService(ctx)
	return nil
}

func (s *CronServer) RunCheckExtractTaskService(ctx context.Context) error {
	ctx, span := otel.Tracer("Service").Start(ctx, "RunCheckExtractTaskService", trace.WithSpanKind(trace.SpanKindProducer))
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 2048)
			runtime.Stack(buf, false)
			span.SetStatus(codes.Error, fmt.Sprintf("%v", err))
		}
		span.End()
	}()

	s.tcu.CheckExtractTask(ctx)
	return nil
}
