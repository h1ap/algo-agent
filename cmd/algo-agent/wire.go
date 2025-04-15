//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"algo-agent/internal/biz"
	"algo-agent/internal/conf"

	"algo-agent/internal/data"
	"algo-agent/internal/server"
	"algo-agent/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		newApp,
	))
}
