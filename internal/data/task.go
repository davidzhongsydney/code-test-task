package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/exp/maps"
	"qantas.com/task/internal/biz"
	"qantas.com/task/internal/encoder"
	"qantas.com/task/model"
	"robpike.io/filter"
)

type taskRepo struct {
	index uint64
	data  *Data
	log   *log.Helper
}

func NewTaskRepo(data *Data, logger log.Logger) biz.ITaskRepo {
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
		r.log.WithContext(ctx).Errorf("taskRepo: Get - Task ID not specified")
		return nil, model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST))
	}

	// Task has been deleted
	if val.DeletedAt != nil {
		return nil, model.ErrorTaskNotFound(string(encoder.TASK_DELETED))
	}

	return &val, nil
}

func (r *taskRepo) Create(ctx context.Context, task *model.Task) (*model.T_Task, error) {
	r.index++
	task.TaskID = r.index

	nt := time.Now()
	newEntry := model.T_Task{Task: *task, T_Internal: model.T_Internal{CreatedAt: &nt}}

	r.data.tasks[task.TaskID] = newEntry
	return &newEntry, nil
}

func (r *taskRepo) Update(ctx context.Context, task *model.Task) (*model.T_Task, error) {

	val, ok := r.data.tasks[task.TaskID]

	// Task not exist
	if !ok {
		return nil, model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST))
	}

	// Task has been deleted
	if val.DeletedAt != nil {
		return nil, model.ErrorTaskNotFound(string(encoder.TASK_DELETED))
	}

	val.Task = *task
	nt := time.Now()
	val.T_Internal.UpdatedAt = &nt

	r.data.tasks[task.TaskID] = val

	return &val, nil
}

func (r *taskRepo) Delete(ctx context.Context, id uint64) error {
	val, ok := r.data.tasks[id]

	// Task not exist
	if !ok {
		return model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST))
	}

	// Task has been deleted
	if val.DeletedAt != nil {
		return model.ErrorTaskNotFound(string(encoder.TASK_DELETED))
	}

	nt := time.Now()
	val.DeletedAt = &nt
	r.data.tasks[id] = val

	return nil
}

func (r *taskRepo) Empty(ctx context.Context) error {
	r.data.tasks = make(map[uint64]model.T_Task)
	r.index = 0

	return nil
}
