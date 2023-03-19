package biz_test

import (
	"context"
	"os"
	"testing"

	errors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"qantas.com/task/internal/biz"
	"qantas.com/task/internal/encoder"
	"qantas.com/task/mocks"
	"qantas.com/task/model"
)

func (uts *BizTestSuite) SetupTest() {
	taskRepoMock := mocks.TaskRepo{}
	logger := log.With(log.NewStdLogger(os.Stdout))

	uts.taskRepoMock = taskRepoMock
	uts.context = context.Background()
	uts.logger = logger
}

type BizTestSuite struct {
	suite.Suite
	taskRepoMock mocks.TaskRepo
	context      context.Context
	logger       log.Logger
}

func TestBizTestSuite(t *testing.T) {
	suite.Run(t, &BizTestSuite{})
}

func (uts *BizTestSuite) Test_CreateTask_Success() {
	mT_Task := model.T_Task{Task: model.Task{TaskID: 2, Name: "user", Content: "content"},
		T_Internal: model.T_Internal{CreatedAt: timestamppb.Now()}}

	uts.taskRepoMock.On("Create", mock.Anything, mock.Anything).Return(
		&mT_Task, nil)

	taskUseCase := biz.NewTaskUsecase(&uts.taskRepoMock, uts.logger)
	retTask, err := taskUseCase.CreateTask(uts.context,
		&model.Task{TaskID: 2, Name: "user", Content: "content"})

	uts.Require().Nil(err)
	uts.Require().Equal(mT_Task, *retTask)
}

func (uts *BizTestSuite) Test_CreateTask_DatabaseCreationError() {
	uts.taskRepoMock.On("Create", mock.Anything, mock.Anything).Return(
		nil, model.ErrorTaskCreationError(string(encoder.TASK_CREATION_ERROR)))

	taskUseCase := biz.NewTaskUsecase(&uts.taskRepoMock, uts.logger)
	retTask, err := taskUseCase.CreateTask(uts.context,
		&model.Task{TaskID: 2, Name: "user", Content: "content"})

	se := new(errors.Error)
	uts.Require().True(errors.As(err, &se))
	uts.Require().True(model.IsTaskCreationError(se))
	uts.Require().Nil(retTask)
}

func (uts *BizTestSuite) Test_GetTaskByID_Success() {
	mT_Task := model.T_Task{Task: model.Task{TaskID: 2, Name: "user", Content: "content"},
		T_Internal: model.T_Internal{CreatedAt: timestamppb.Now()}}

	uts.taskRepoMock.On("Get", mock.Anything, mock.Anything).Return(
		&mT_Task, nil)

	taskUseCase := biz.NewTaskUsecase(&uts.taskRepoMock, uts.logger)
	retTask, err := taskUseCase.GetTaskByID(uts.context, 2)

	uts.Require().Nil(err)
	uts.Require().Equal(mT_Task, *retTask)
}

func (uts *BizTestSuite) Test_GetTaskByID_DatabaseNotFoundError() {
	uts.taskRepoMock.On("Get", mock.Anything, mock.Anything).Return(
		nil, model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)))

	taskUseCase := biz.NewTaskUsecase(&uts.taskRepoMock, uts.logger)
	retTask, err := taskUseCase.GetTaskByID(uts.context, 2)

	se := new(errors.Error)
	uts.Require().True(errors.As(err, &se))
	uts.Require().True(model.IsTaskNotFound(se))
	uts.Require().Nil(retTask)
}

func (uts *BizTestSuite) Test_GetTaskByID_TaskIdNotSpecified() {

	taskUseCase := biz.NewTaskUsecase(&uts.taskRepoMock, uts.logger)
	retTask, err := taskUseCase.GetTaskByID(uts.context, 0)

	se := new(errors.Error)
	uts.Require().True(errors.As(err, &se))
	uts.Require().True(model.IsTaskIdUnspecified(se))
	uts.Require().Nil(retTask)
}

func (uts *BizTestSuite) Test_UpdateTaskByID_Success() {
	mT_Task := model.T_Task{Task: model.Task{TaskID: 2, Name: "user", Content: "content"},
		T_Internal: model.T_Internal{CreatedAt: timestamppb.Now(), UpdatedAt: timestamppb.Now()}}

	uts.taskRepoMock.On("Update", mock.Anything, mock.Anything).Return(
		&mT_Task, nil)

	taskUseCase := biz.NewTaskUsecase(&uts.taskRepoMock, uts.logger)
	retTask, err := taskUseCase.UpdateTaskByID(uts.context,
		&model.Task{TaskID: 2, Name: "user", Content: "content"})

	uts.Require().Nil(err)
	uts.Require().Equal(mT_Task, *retTask)
}

func (uts *BizTestSuite) Test_UpdateTaskByID_DatabaseTaskNotFound() {
	uts.taskRepoMock.On("Update", mock.Anything, mock.Anything).Return(
		nil, model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)))

	taskUseCase := biz.NewTaskUsecase(&uts.taskRepoMock, uts.logger)
	retTask, err := taskUseCase.UpdateTaskByID(uts.context,
		&model.Task{TaskID: 2, Name: "user", Content: "content"})

	se := new(errors.Error)
	uts.Require().True(errors.As(err, &se))
	uts.Require().True(model.IsTaskNotFound(se))
	uts.Require().Nil(retTask)
}

func (uts *BizTestSuite) Test_UpdateTaskByID_TaskIdNotSpecified() {
	taskUseCase := biz.NewTaskUsecase(&uts.taskRepoMock, uts.logger)
	retTask, err := taskUseCase.UpdateTaskByID(uts.context,
		&model.Task{Name: "user", Content: "content"})

	se := new(errors.Error)
	uts.Require().True(errors.As(err, &se))
	uts.Require().True(model.IsTaskIdUnspecified(se))
	uts.Require().Nil(retTask)
}

func (uts *BizTestSuite) Test_DeleteTaskByID_Success() {
	uts.taskRepoMock.On("Delete", mock.Anything, mock.Anything).Return(nil)

	taskUseCase := biz.NewTaskUsecase(&uts.taskRepoMock, uts.logger)
	err := taskUseCase.DeleteTaskByID(uts.context, 3)

	uts.Require().Nil(err)
}

func (uts *BizTestSuite) Test_DeleteTaskByID_DatabaseTaskNotFound() {
	uts.taskRepoMock.On("Delete", mock.Anything, mock.Anything).Return(
		model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)))

	taskUseCase := biz.NewTaskUsecase(&uts.taskRepoMock, uts.logger)
	err := taskUseCase.DeleteTaskByID(uts.context, 3)

	se := new(errors.Error)
	uts.Require().True(errors.As(err, &se))
	uts.Require().True(model.IsTaskNotFound(se))
}

func (uts *BizTestSuite) Test_DeleteTaskByID_TaskIdNotSpecified() {
	taskUseCase := biz.NewTaskUsecase(&uts.taskRepoMock, uts.logger)
	err := taskUseCase.DeleteTaskByID(uts.context, 0)

	se := new(errors.Error)
	uts.Require().True(errors.As(err, &se))
	uts.Require().True(model.IsTaskIdUnspecified(se))
}

func (uts *BizTestSuite) Test_ListTasks_Success() {
	mT_Task1 := model.T_Task{Task: model.Task{TaskID: 1, Name: "user1", Content: "content1"},
		T_Internal: model.T_Internal{CreatedAt: timestamppb.Now()}}
	mT_Task2 := model.T_Task{Task: model.Task{TaskID: 1, Name: "user2", Content: "content2"},
		T_Internal: model.T_Internal{CreatedAt: timestamppb.Now()}}

	uts.taskRepoMock.On("List", mock.Anything, mock.Anything).Return(
		[]model.T_Task{mT_Task1, mT_Task2}, nil)

	taskUseCase := biz.NewTaskUsecase(&uts.taskRepoMock, uts.logger)
	retTasks, err := taskUseCase.ListTasks(uts.context)

	uts.Require().Nil(err)
	if retTasks[0].TaskID == 1 {
		uts.Require().Equal(retTasks[0], mT_Task1)
		uts.Require().Equal(retTasks[1], mT_Task2)
	} else {
		uts.Require().Equal(retTasks[0], mT_Task2)
		uts.Require().Equal(retTasks[1], mT_Task1)
	}
}

func (uts *BizTestSuite) Test_ListTasks_DatabaseError() {
	uts.taskRepoMock.On("List", mock.Anything, mock.Anything).Return(
		nil, model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)))

	taskUseCase := biz.NewTaskUsecase(&uts.taskRepoMock, uts.logger)
	retTasks, err := taskUseCase.ListTasks(uts.context)

	uts.Require().Nil(retTasks)
	se := new(errors.Error)
	uts.Require().True(errors.As(err, &se))
	uts.Require().True(model.IsTaskNotFound(se))
}
