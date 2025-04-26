package cron

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type ServerOption func(o *Server)

func WithContext(ctx context.Context) ServerOption {
	return func(s *Server) {
		s.baseCtx = ctx
	}
}

func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(logger)
	}
}

func RegisterFunc(spec string, taskName string, fun func()) ServerOption {
	return func(s *Server) {
		entryId, err := s.AddFunc(spec, fun)
		if err != nil {
			s.log.WithContext(s.baseCtx).Errorw("register func is error")
		}
		s.RegisterMapEntry(entryId, taskName)
	}
}
