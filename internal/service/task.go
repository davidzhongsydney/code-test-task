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

func (t *TaskService) ListTasks(ctx context.Context) ([]*model.Task, error) {
	tasks, err := t.uc.ListTasks(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
