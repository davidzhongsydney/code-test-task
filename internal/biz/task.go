package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"qantas.com/task/internal/encoder"
	"qantas.com/task/model"
)

type TaskRepo interface {
	Create(context.Context, *model.Task) (*model.T_Task, error)
	Get(context.Context, uint64) (*model.T_Task, error)
	Update(context.Context, *model.Task) (*model.T_Task, error)
	Delete(context.Context, uint64) error
	List(context.Context) ([]model.T_Task, error)
	Empty(context.Context) error
}

type TaskUsecase struct {
	repo TaskRepo
	log  *log.Helper
}

func NewTaskUsecase(repo TaskRepo, logger log.Logger) *TaskUsecase {
	return &TaskUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *TaskUsecase) CreateTask(ctx context.Context, t *model.Task) (*model.T_Task, error) {
	uc.log.WithContext(ctx).Infof("TaskUsecase: CreateTask: %v", *t)
	return uc.repo.Create(ctx, t)
}

func (uc *TaskUsecase) GetTaskByID(ctx context.Context, id uint64) (*model.T_Task, error) {
	uc.log.WithContext(ctx).Infof("TaskUsecase: GetTaskByID: %v", id)
	if id == 0 {
		uc.log.WithContext(ctx).Error("TaskUsecase: GetTaskByID - Task ID not specified")
		return nil, model.ErrorTaskIdUnspecified(string(encoder.TASK_ID_NOT_SPECIFIED))
	}

	return uc.repo.Get(ctx, id)
}

func (uc *TaskUsecase) DeleteTaskByID(ctx context.Context, id uint64) error {
	uc.log.WithContext(ctx).Infof("TaskUsecase: DeleteTaskByID: %v", id)
	if id == 0 {
		uc.log.WithContext(ctx).Error("TaskUsecase: DeleteTaskByID - Task ID not specified")
		return model.ErrorTaskIdUnspecified(string(encoder.TASK_ID_NOT_SPECIFIED))
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *TaskUsecase) UpdateTaskByID(ctx context.Context, t *model.Task) (*model.T_Task, error) {
	uc.log.WithContext(ctx).Infof("TaskUsecase: UpdateTaskByID: %v", *t)
	if t.TaskID == 0 {
		uc.log.WithContext(ctx).Error("TaskUsecase: UpdateTaskByID - Task ID not specified")
		return nil, model.ErrorTaskIdUnspecified(string(encoder.TASK_ID_NOT_SPECIFIED))
	}
	return uc.repo.Update(ctx, t)
}

func (uc *TaskUsecase) ListTasks(ctx context.Context) ([]model.T_Task, error) {
	uc.log.WithContext(ctx).Infof("TaskUsecase: ListTasks")
	return uc.repo.List(ctx)
}
