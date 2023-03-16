package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"qantas.com/task/internal/biz"
	"qantas.com/task/model"
)

type taskRepo struct {
	data *Data
	log  *log.Helper
}

func NewTaskRepo(data *Data, logger log.Logger) biz.TaskRepo {
	return &taskRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *taskRepo) ListAll(context.Context) ([]*model.Task, error) {
	return nil, nil
}
