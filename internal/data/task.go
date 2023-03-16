package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"qantas.com/task/internal/biz"
	"qantas.com/task/model"
)

type taskRepo struct {
	index uint64
	data  *Data
	log   *log.Helper
}

func NewTaskRepo(data *Data, logger log.Logger) biz.TaskRepo {
	return &taskRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *taskRepo) List(context.Context) ([]model.Task, error) {
	return r.data.tasks, nil
}

func (r *taskRepo) Create(ctx context.Context, task *model.Task) (*model.Task, error) {
	task.TaskID = r.index
	r.index++

	task.CreatedAt = time.Now()

	r.data.tasks = append(r.data.tasks, *task)
	return task, nil
}
