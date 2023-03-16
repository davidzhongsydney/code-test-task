package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	return maps.Values(r.data.tasks), nil
}

func (r *taskRepo) Get(ctx context.Context, id uint64) (*model.Task, error) {
	val, ok := r.data.tasks[id]

	// Task not exist
	if !ok {

	}

	// Task has been deleted
	if val.DeletedAt != nil {

	}

	return &val, nil
}

func (r *taskRepo) Create(ctx context.Context, task *model.Task) (*model.Task, error) {
	r.index++
	task.TaskID = r.index
	task.CreatedAt = timestamppb.Now()
	task.UpdatedAt = nil
	task.DeletedAt = nil

	r.data.tasks[task.TaskID] = *task
	return task, nil
}

func (r *taskRepo) Update(ctx context.Context, task *model.Task) (*model.Task, error) {

	val, ok := r.data.tasks[task.TaskID]

	// Task not exist
	if !ok {

	}

	// Task has been deleted
	if val.DeletedAt != nil {

	}

	task.CreatedAt = val.CreatedAt
	task.UpdatedAt = timestamppb.Now()

	r.data.tasks[task.TaskID] = *task

	return task, nil
}

func (r *taskRepo) Delete(ctx context.Context, id uint64) error {
	val, ok := r.data.tasks[id]

	// Task not exist
	if !ok {

	}

	// Task has been deleted
	if val.DeletedAt != nil {

	}

	val.DeletedAt = timestamppb.Now()
	r.data.tasks[id] = val

	return nil
}
