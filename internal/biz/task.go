package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"qantas.com/task/model"
)

type TaskRepo interface {
	Create(context.Context, *model.Task) (*model.Task, error)
	Get(context.Context, uint64) (*model.Task, error)
	Update(context.Context, *model.Task) (*model.Task, error)
	Delete(context.Context, uint64) error
	List(context.Context) ([]model.Task, error)
}

type TaskUsecase struct {
	repo TaskRepo
	log  *log.Helper
}

func NewTaskUsecase(repo TaskRepo, logger log.Logger) *TaskUsecase {
	return &TaskUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *TaskUsecase) CreateTask(ctx context.Context, t *model.Task) (*model.Task, error) {
	uc.log.WithContext(ctx).Infof("CreateTask: %v", *t)
	return uc.repo.Create(ctx, t)
}

func (uc *TaskUsecase) GetTaskByID(ctx context.Context, id uint64) (*model.Task, error) {
	uc.log.WithContext(ctx).Infof("GetTaskByID: %v", id)
	return uc.repo.Get(ctx, id)
}

func (uc *TaskUsecase) DeleteTaskByID(ctx context.Context, id uint64) error {
	uc.log.WithContext(ctx).Infof("DeleteTaskByID: %v", id)
	return uc.repo.Delete(ctx, id)
}

func (uc *TaskUsecase) UpdateTaskByID(ctx context.Context, t *model.Task) (*model.Task, error) {
	uc.log.WithContext(ctx).Infof("UpdateTaskByID: %v", *t)
	return uc.repo.Update(ctx, t)
}

func (uc *TaskUsecase) ListTasks(ctx context.Context) ([]model.Task, error) {
	uc.log.WithContext(ctx).Infof("ListTasks")
	return uc.repo.List(ctx)
}
