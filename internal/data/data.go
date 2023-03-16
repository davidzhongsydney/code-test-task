package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"qantas.com/task/internal/conf"
	"qantas.com/task/model"
)

var ProviderSet = wire.NewSet(NewData, NewTaskRepo)

type Data struct {
	tasks map[uint64]model.Task
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{tasks: make(map[uint64]model.Task)}, cleanup, nil
}
