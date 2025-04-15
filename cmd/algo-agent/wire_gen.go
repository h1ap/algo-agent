// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"algo-agent/internal/biz"
	"algo-agent/internal/conf"
	"algo-agent/internal/data"
	"algo-agent/internal/server"
	"algo-agent/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	ossStore, err := data.NewOSSRepo(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	ossUsecase := biz.NewOSSUsecase(ossStore, logger)
	ossServer := service.NewOSSServer(ossUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, ossServer, logger)
	httpServer := server.NewHTTPServer(confServer, ossServer, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
	}, nil
}
