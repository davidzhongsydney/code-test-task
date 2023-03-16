package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"qantas.com/task/model"
)

type TaskRepo interface {
	// Save(context.Context, *Greeter) (*Greeter, error)
	// Update(context.Context, *Greeter) (*Greeter, error)
	// FindByID(context.Context, int64) (*Greeter, error)
	// ListByHello(context.Context, string) ([]*Greeter, error)
	ListAll(context.Context) ([]*model.Task, error)
}

type TaskUsecase struct {
	repo TaskRepo
	log  *log.Helper
}

func NewTaskUsecase(repo TaskRepo, logger log.Logger) *TaskUsecase {
	return &TaskUsecase{repo: repo, log: log.NewHelper(logger)}
}

// func (uc *TaskUsecase) CreateTask(ctx context.Context, t *Task) (*Task, error) {
// 	uc.log.WithContext(ctx).Infof("CreateTask: %v", g.Hello)
// 	return uc.repo.Save(ctx, g)
// }

func (uc *TaskUsecase) ListTasks(ctx context.Context) ([]*model.Task, error) {
	uc.log.WithContext(ctx).Infof("ListTasks")
	return uc.repo.ListAll(ctx)
}
