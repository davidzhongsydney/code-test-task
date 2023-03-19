package service_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	errors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"qantas.com/task/internal/biz"
	"qantas.com/task/internal/encoder"
	"qantas.com/task/internal/service"
	"qantas.com/task/mocks"
	"qantas.com/task/model"
)

func TestCreateAndUpdateTaskService(t *testing.T) {
	type testCases struct {
		description      string
		mockMethod       string
		callMethod       string
		mockInputTask    *model.Task
		mockSingleReturn *model.T_Task
		mockReturnError  *errors.Error
	}

	for _, scenario := range []testCases{
		{
			description:   "task creation success",
			mockMethod:    "Create",
			callMethod:    "CreateTask",
			mockInputTask: &model.Task{TaskID: 2, Name: "user", Content: "content"},
			mockSingleReturn: &model.T_Task{Task: model.Task{TaskID: 2, Name: "user", Content: "content"},
				T_Internal: model.T_Internal{CreatedAt: timestamppb.Now()}},
		},
		{
			description:     "task creation failed",
			mockMethod:      "Create",
			callMethod:      "CreateTask",
			mockInputTask:   &model.Task{TaskID: 2, Name: "user", Content: "content"},
			mockReturnError: model.ErrorTaskCreationError(string(encoder.TASK_CREATION_ERROR)),
		},
		{
			description:   "update task by id success",
			mockMethod:    "Update",
			callMethod:    "UpdateTaskByID",
			mockInputTask: &model.Task{TaskID: 2, Name: "user", Content: "content"},
			mockSingleReturn: &model.T_Task{Task: model.Task{TaskID: 2, Name: "user", Content: "content"},
				T_Internal: model.T_Internal{CreatedAt: timestamppb.Now()}},
		},
		{
			description:     "update task by id failed",
			mockMethod:      "Update",
			callMethod:      "UpdateTaskByID",
			mockInputTask:   &model.Task{TaskID: 2, Name: "user", Content: "content"},
			mockReturnError: model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			requires := require.New(t)

			logger := log.With(log.NewStdLogger(os.Stdout))
			context := context.Background()

			taskRepoMock := mocks.TaskRepo{}

			if scenario.mockReturnError != nil {
				taskRepoMock.On(scenario.mockMethod, mock.Anything, mock.Anything).Return(nil, scenario.mockReturnError)
			} else {
				taskRepoMock.On(scenario.mockMethod, mock.Anything, mock.Anything).Return(scenario.mockSingleReturn, nil)
			}

			taskUseCase := biz.NewTaskUsecase(&taskRepoMock, logger)
			taskService := service.NewTaskService(taskUseCase)

			method := reflect.ValueOf(taskService).MethodByName(scenario.callMethod)

			var params []reflect.Value

			params = append(params, reflect.ValueOf(context))
			params = append(params, reflect.ValueOf(scenario.mockInputTask))

			ret := method.Call(params)

			if scenario.mockReturnError != nil {
				retTask := ret[0].Interface()
				retError := ret[1].Interface().(*errors.Error)

				requires.Nil(retTask)
				requires.Equal(*scenario.mockReturnError, *retError)
			} else {
				retTask := ret[0].Interface().(*model.T_Task)
				retError := ret[1].Interface()

				requires.Nil(retError)
				requires.Equal(*scenario.mockSingleReturn, *retTask)
			}
		})
	}
}

func TestGetAndDeleteTaskService(t *testing.T) {
	type testCases struct {
		description      string
		mockMethod       string
		callMethod       string
		mockInputInteger uint64
		mockSingleReturn *model.T_Task
		mockReturnError  *errors.Error
		numOfReturnArg   uint
	}

	for _, scenario := range []testCases{
		{
			description:      "get task by id success",
			mockMethod:       "Get",
			callMethod:       "GetTaskByID",
			mockInputInteger: 2,
			numOfReturnArg:   2,
			mockSingleReturn: &model.T_Task{Task: model.Task{TaskID: 2, Name: "user", Content: "content"},
				T_Internal: model.T_Internal{CreatedAt: timestamppb.Now()}},
		},
		{
			description:      "get task by id failure",
			mockMethod:       "Get",
			callMethod:       "GetTaskByID",
			mockInputInteger: 2,
			numOfReturnArg:   2,
			mockReturnError:  model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
		},
		{
			description:      "delete task by id success",
			mockMethod:       "Delete",
			callMethod:       "DeleteTaskByID",
			mockInputInteger: 2,
			numOfReturnArg:   1,
		},
		{
			description:      "delete task by id failed",
			mockMethod:       "Delete",
			callMethod:       "DeleteTaskByID",
			mockInputInteger: 2,
			numOfReturnArg:   1,
			mockReturnError:  model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			requires := require.New(t)

			logger := log.With(log.NewStdLogger(os.Stdout))
			context := context.Background()

			taskRepoMock := mocks.TaskRepo{}

			if scenario.numOfReturnArg == 1 {
				if scenario.mockReturnError != nil {
					taskRepoMock.On(scenario.mockMethod, mock.Anything, mock.Anything).Return(scenario.mockReturnError)
				} else {
					taskRepoMock.On(scenario.mockMethod, mock.Anything, mock.Anything).Return(nil)
				}
			} else {
				if scenario.mockReturnError != nil {
					taskRepoMock.On(scenario.mockMethod, mock.Anything, mock.Anything).Return(nil, scenario.mockReturnError)
				} else {
					taskRepoMock.On(scenario.mockMethod, mock.Anything, mock.Anything).Return(scenario.mockSingleReturn, nil)
				}
			}

			taskUseCase := biz.NewTaskUsecase(&taskRepoMock, logger)
			taskService := service.NewTaskService(taskUseCase)

			method := reflect.ValueOf(taskService).MethodByName(scenario.callMethod)

			var params []reflect.Value

			params = append(params, reflect.ValueOf(context))
			params = append(params, reflect.ValueOf(scenario.mockInputInteger))

			ret := method.Call(params)

			if scenario.numOfReturnArg == 1 {
				if scenario.mockReturnError != nil {
					retError := ret[0].Interface().(*errors.Error)
					requires.Equal(*scenario.mockReturnError, *retError)
				} else {
					retError := ret[0].Interface()
					requires.Nil(retError)
				}
			} else {
				if scenario.mockReturnError != nil {
					retTask := ret[0].Interface()
					retError := ret[1].Interface().(*errors.Error)

					requires.Nil(retTask)
					requires.Equal(*scenario.mockReturnError, *retError)
				} else {
					retTask := ret[0].Interface().(*model.T_Task)
					retError := ret[1].Interface()

					requires.Nil(retError)
					requires.Equal(*scenario.mockSingleReturn, *retTask)
				}
			}
		})
	}
}

func TestListTasksService(t *testing.T) {
	type testCases struct {
		description     string
		mockMethod      string
		callMethod      string
		mockReturn      []model.T_Task
		mockReturnError *errors.Error
	}

	for _, scenario := range []testCases{
		{
			description: "List tasks failed",
			mockMethod:  "List",
			callMethod:  "ListTasks",
			mockReturn: []model.T_Task{
				{Task: model.Task{TaskID: 1, Name: "user1", Content: "content1"},
					T_Internal: model.T_Internal{CreatedAt: timestamppb.Now()}},
				{Task: model.Task{TaskID: 2, Name: "user2", Content: "content2"},
					T_Internal: model.T_Internal{CreatedAt: timestamppb.Now()}}},
			mockReturnError: model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
		},
		{
			description: "List tasks success",
			mockMethod:  "List",
			callMethod:  "ListTasks",
			mockReturn: []model.T_Task{
				{Task: model.Task{TaskID: 1, Name: "user1", Content: "content1"},
					T_Internal: model.T_Internal{CreatedAt: timestamppb.Now()}},
				{Task: model.Task{TaskID: 2, Name: "user2", Content: "content2"},
					T_Internal: model.T_Internal{CreatedAt: timestamppb.Now()}}},
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			requires := require.New(t)

			logger := log.With(log.NewStdLogger(os.Stdout))
			context := context.Background()

			taskRepoMock := mocks.TaskRepo{}

			if scenario.mockReturnError != nil {
				taskRepoMock.On(scenario.mockMethod, mock.Anything).Return(nil, scenario.mockReturnError)
			} else {
				taskRepoMock.On(scenario.mockMethod, mock.Anything).Return(scenario.mockReturn, nil)
			}

			taskUseCase := biz.NewTaskUsecase(&taskRepoMock, logger)
			taskService := service.NewTaskService(taskUseCase)

			method := reflect.ValueOf(taskService).MethodByName(scenario.callMethod)

			var params []reflect.Value

			params = append(params, reflect.ValueOf(context))

			ret := method.Call(params)

			if scenario.mockReturnError != nil {
				retTasks := ret[0].Interface()
				retError := ret[1].Interface().(*errors.Error)

				requires.Nil(retTasks)
				requires.Equal(*scenario.mockReturnError, *retError)
			} else {
				retTask := ret[0].Interface().([]model.T_Task)
				retError := ret[1].Interface()

				requires.Nil(retError)
				requires.Equal(scenario.mockReturn, retTask)
			}
		})
	}
}
