//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"qantas.com/task/internal/biz"
	conf "qantas.com/task/internal/conf"
	"qantas.com/task/internal/data"
	"qantas.com/task/internal/server"
	"qantas.com/task/internal/service"

	"github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Data, log.Logger) (server.Server, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet))
}
