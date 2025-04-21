package job

import (
	"algo-agent/internal/service"
	"context"
	"net/url"
)

type JobServer struct {
	j *service.JobServer
}

func NewJobServer(j *service.JobServer) *JobServer {
	return &JobServer{j: j}
}

func (j *JobServer) Start(context.Context) error {
	j.j.Start()
	return nil
}

func (j *JobServer) Stop(context.Context) error {
	j.j.Stop()
	return nil
}

func (j *JobServer) Endpoint() (*url.URL, error) {
	// todo need implement
	return nil, nil
}
