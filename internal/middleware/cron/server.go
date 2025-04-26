package cron

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/robfig/cron/v3"
	"golang.org/x/net/context"
	"sync"
)

var (
	// "归一化"写法:
	// Server结构体需要实现 transport.Server 这个interface对应的方法
	_ transport.Server = (*Server)(nil)
)

type Server struct {
	sync.RWMutex
	Cron       *cron.Cron
	log        *log.Helper
	baseCtx    context.Context
	MapEntryId map[string]cron.EntryID
	started    bool
	err        error
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		Cron:       cron.New(),
		MapEntryId: make(map[string]cron.EntryID),
		started:    false,
	}
	srv.init(opts...)
	return srv
}

func (s *Server) init(opts ...ServerOption) {
	for _, o := range opts {
		o(s)
	}
}

func (s *Server) Name() string {
	return "task1"
}

func (s *Server) RegisterMapEntry(id cron.EntryID, taskName string) {
	s.Lock()
	defer s.Unlock()
	if len(taskName) > 0 {
		s.MapEntryId[taskName] = id
		s.log.WithContext(s.baseCtx).Infow("[robfig-cron] register task1:", taskName, "id:", id)
	}
}

func (s *Server) GetMapEntry() map[string]cron.EntryID {
	return s.MapEntryId
}

func (s *Server) AddFunc(spec string, fun func()) (cron.EntryID, error) {
	return s.Cron.AddFunc(spec, fun)
}

func (s *Server) Start(ctx context.Context) error {
	if s.err != nil {
		return s.err
	}
	if s.started {
		return nil
	}
	s.Cron.Start()
	s.log.WithContext(ctx).Info("[robfig-corn] server starting")
	s.baseCtx = ctx
	s.started = true
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.log.WithContext(s.baseCtx).Info("[robfig-cron] server stopping")
	s.started = false
	s.Cron.Stop()
	return nil
}
