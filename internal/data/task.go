package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/types/known/timestamppb"
	"qantas.com/task/internal/biz"
	"qantas.com/task/model"
	"robpike.io/filter"
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

func (r *taskRepo) List(context.Context) ([]model.T_Task, error) {
	tasks := maps.Values(r.data.tasks)
	result := filter.Choose(tasks, func(task model.T_Task) bool {
		return task.DeletedAt == nil
	}).([]model.T_Task)

	return result, nil
}

func (r *taskRepo) Get(ctx context.Context, id uint64) (*model.T_Task, error) {
	val, ok := r.data.tasks[id]

	// Task not exist
	if !ok {

	}

	// Task has been deleted
	if val.DeletedAt != nil {

	}

	return &val, nil
}

func (r *taskRepo) Create(ctx context.Context, task *model.Task) (*model.T_Task, error) {
	r.index++
	task.TaskID = r.index

	newEntry := model.T_Task{Task: *task, T_Internal: model.T_Internal{CreatedAt: timestamppb.Now()}}

	r.data.tasks[task.TaskID] = newEntry
	return &newEntry, nil
}

func (r *taskRepo) Update(ctx context.Context, task *model.Task) (*model.T_Task, error) {

	val, ok := r.data.tasks[task.TaskID]

	// Task not exist
	if !ok {

	}

	// Task has been deleted
	if val.DeletedAt != nil {

	}

	val.Task = *task
	val.T_Internal.UpdatedAt = timestamppb.Now()

	r.data.tasks[task.TaskID] = val

	return &val, nil
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
