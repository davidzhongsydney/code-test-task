package service

import (
	"context"

	"qantas.com/task/internal/biz"
	"qantas.com/task/model"
)

type TaskService struct {
	uc *biz.TaskUsecase
}

func NewTaskService(uc *biz.TaskUsecase) *TaskService {
	return &TaskService{uc: uc}
}

func (t *TaskService) ListTasks(ctx context.Context) ([]model.T_Task, error) {
	tasks, err := t.uc.ListTasks(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *TaskService) CreateTask(ctx context.Context, task *model.Task) (*model.T_Task, error) {
	taskTable, err := t.uc.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}
	return taskTable, nil
}

func (t *TaskService) GetTaskByID(ctx context.Context, id uint64) (*model.T_Task, error) {
	task, err := t.uc.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *TaskService) UpdateTaskByID(ctx context.Context, task *model.Task) (*model.T_Task, error) {
	t_task, err := t.uc.UpdateTaskByID(ctx, task)
	if err != nil {
		return nil, err
	}
	return t_task, nil
}

func (t *TaskService) DeleteTaskByID(ctx context.Context, id uint64) error {
	err := t.uc.DeleteTaskByID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
