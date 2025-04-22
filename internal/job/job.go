package job

import (
	"algo-agent/internal/service"
	"algo-agent/internal/utils"
	"context"
	"net/url"
	"strconv"
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
	ip := utils.GetLocalIP()
	u := &url.URL{
		Scheme: "http",
		Host:   ip + ":" + strconv.Itoa(int(8001)),
	}
	return u, nil
}
