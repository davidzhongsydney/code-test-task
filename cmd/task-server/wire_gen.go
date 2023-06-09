// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"qantas.com/task/internal/biz"
	"qantas.com/task/internal/conf"
	"qantas.com/task/internal/data"
	"qantas.com/task/internal/server"
	"qantas.com/task/internal/service"
)

// Injectors from wire.go:

func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger, ctx context.Context) (server.IServer, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	iTaskRepo := data.NewTaskRepo(dataData, logger)
	taskUsecase := biz.NewTaskUsecase(iTaskRepo, logger)
	taskService := service.NewTaskService(taskUsecase, logger)
	iTaskHTTPHandler := server.NewTaskHTTPHandler(taskService, logger, ctx)
	iServer := server.NewHTTPServer(confServer, logger, iTaskHTTPHandler)
	return iServer, func() {
		cleanup()
	}, nil
}
