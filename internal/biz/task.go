package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"qantas.com/task/model"
)

type TaskRepo interface {
	Create(context.Context, *model.Task) (*model.Task, error)
	// Update(context.Context, *Greeter) (*Greeter, error)
	// FindByID(context.Context, int64) (*Greeter, error)
	// ListByHello(context.Context, string) ([]*Greeter, error)
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
	uc.log.WithContext(ctx).Infof("CreateTask: %v", t)
	return uc.repo.Create(ctx, t)
}

func (uc *TaskUsecase) ListTasks(ctx context.Context) ([]model.Task, error) {
	uc.log.WithContext(ctx).Infof("ListTasks")
	return uc.repo.List(ctx)
}
