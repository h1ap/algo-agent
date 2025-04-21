package server

import (
	"algo-agent/internal/job"
	"algo-agent/internal/service"
)

func NewJobServer(service *service.JobServer) *job.JobServer {
	js := job.NewJobServer(service)
	return js
}
