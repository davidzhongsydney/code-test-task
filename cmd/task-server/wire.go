//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"qantas.com/task/internal/biz"
	conf "qantas.com/task/internal/conf"
	"qantas.com/task/internal/data"
	"qantas.com/task/internal/server"
	"qantas.com/task/internal/service"

	"github.com/google/wire"
)

func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger, ctx context.Context) (server.IServer, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet))
}
