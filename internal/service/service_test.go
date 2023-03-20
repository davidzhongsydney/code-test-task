package service_test

import (
	"context"
	"os"
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

func TestTaskService(t *testing.T) {
	type testCases struct {
		description      string
		mockMethod       string
		callMethod       string
		mockInputTask    *model.Task
		mockInputInteger uint64
		mockReturnTask   *model.T_Task
		mockReturnTasks  []model.T_Task
		mockReturnError  *errors.Error
		expectedError    *errors.Error
	}

	for _, scenario := range []testCases{
		{
			description:   "task creation success",
			mockMethod:    "Create",
			callMethod:    "CreateTask",
			mockInputTask: &model.Task{TaskID: 2, Name: "user", Content: "content"},
			mockReturnTask: &model.T_Task{Task: model.Task{TaskID: 2, Name: "user", Content: "content"},
				T_Internal: model.T_Internal{CreatedAt: &timestamppb.Timestamp{Seconds: 1879287443, Nanos: 196685200}}},
		},
		{
			description:     "task creation failed - creation error",
			mockMethod:      "Create",
			callMethod:      "CreateTask",
			mockInputTask:   &model.Task{TaskID: 2, Name: "user", Content: "content"},
			mockReturnError: model.ErrorTaskCreationError(string(encoder.TASK_CREATION_ERROR)),
			expectedError:   model.ErrorTaskCreationError(string(encoder.TASK_CREATION_ERROR)),
		},
		{
			description:   "update task by id success",
			mockMethod:    "Update",
			callMethod:    "UpdateTaskByID",
			mockInputTask: &model.Task{TaskID: 2, Name: "user", Content: "content"},
			mockReturnTask: &model.T_Task{Task: model.Task{TaskID: 2, Name: "user", Content: "content"},
				T_Internal: model.T_Internal{CreatedAt: &timestamppb.Timestamp{Seconds: 1879287443, Nanos: 196685200},
					UpdatedAt: &timestamppb.Timestamp{Seconds: 1979287443, Nanos: 20685200}}},
		},
		{
			description:     "update task by id failed - task not found",
			mockMethod:      "Update",
			callMethod:      "UpdateTaskByID",
			mockInputTask:   &model.Task{TaskID: 2, Name: "user", Content: "content"},
			mockReturnError: model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
			expectedError:   model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
		},
		{
			description:      "get task by id success",
			mockMethod:       "Get",
			callMethod:       "GetTaskByID",
			mockInputInteger: 2,
			mockReturnTask: &model.T_Task{Task: model.Task{TaskID: 2, Name: "user", Content: "content"},
				T_Internal: model.T_Internal{CreatedAt: &timestamppb.Timestamp{Seconds: 1879287443, Nanos: 196685200}}},
		},
		{
			description:      "get task by id failed - task not found",
			mockMethod:       "Get",
			callMethod:       "GetTaskByID",
			mockInputInteger: 2,
			mockReturnError:  model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
			expectedError:    model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
		},
		{
			description:      "get task by id failed - id not specified",
			mockMethod:       "Get",
			callMethod:       "GetTaskByID",
			mockInputInteger: 0,
			mockReturnError:  model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
			expectedError:    model.ErrorTaskIdUnspecified(string(encoder.TASK_ID_NOT_SPECIFIED)),
		},
		{
			description:      "delete task by id success",
			mockMethod:       "Delete",
			callMethod:       "DeleteTaskByID",
			mockInputInteger: 2,
		},
		{
			description:      "delete task by id failed - task not found",
			mockMethod:       "Delete",
			callMethod:       "DeleteTaskByID",
			mockInputInteger: 2,
			mockReturnError:  model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
			expectedError:    model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
		},
		{
			description:      "delete task by id failed - id not specified",
			mockMethod:       "Delete",
			callMethod:       "DeleteTaskByID",
			mockInputInteger: 0,
			mockReturnError:  model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
			expectedError:    model.ErrorTaskIdUnspecified(string(encoder.TASK_ID_NOT_SPECIFIED)),
		},
		{
			description: "List tasks success",
			mockMethod:  "List",
			callMethod:  "ListTasks",
			mockReturnTasks: []model.T_Task{
				{Task: model.Task{TaskID: 1, Name: "user1", Content: "content1"},
					T_Internal: model.T_Internal{CreatedAt: timestamppb.Now()}},
				{Task: model.Task{TaskID: 2, Name: "user2", Content: "content2"},
					T_Internal: model.T_Internal{CreatedAt: timestamppb.Now()}}},
		},
		{
			description:     "List tasks failed - database timeout",
			mockMethod:      "List",
			callMethod:      "ListTasks",
			mockReturnError: model.ErrorTaskDbTimeout(string(encoder.TASK_DATABASE_TIMEOUT)),
			expectedError:   model.ErrorTaskDbTimeout(string(encoder.TASK_DATABASE_TIMEOUT)),
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			requires := require.New(t)

			logger := log.With(log.NewStdLogger(os.Stdout))
			context := context.Background()

			taskRepoMock := mocks.TaskRepo{}

			// Set up dabase mock
			switch scenario.callMethod {
			case "CreateTask", "UpdateTaskByID", "GetTaskByID":
				if scenario.mockReturnError != nil {
					taskRepoMock.On(scenario.mockMethod, mock.Anything, mock.Anything).Return(nil, scenario.mockReturnError)
				} else {
					taskRepoMock.On(scenario.mockMethod, mock.Anything, mock.Anything).Return(scenario.mockReturnTask, nil)
				}
			case "DeleteTaskByID":
				if scenario.mockReturnError != nil {
					taskRepoMock.On(scenario.mockMethod, mock.Anything, mock.Anything).Return(scenario.mockReturnError)
				} else {
					taskRepoMock.On(scenario.mockMethod, mock.Anything, mock.Anything).Return(nil)
				}
			case "ListTasks":
				if scenario.mockReturnError != nil {
					taskRepoMock.On(scenario.mockMethod, mock.Anything).Return(nil, scenario.mockReturnError)
				} else {
					taskRepoMock.On(scenario.mockMethod, mock.Anything).Return(scenario.mockReturnTasks, nil)
				}
			}

			taskUseCase := biz.NewTaskUsecase(&taskRepoMock, logger)
			taskService := service.NewTaskService(taskUseCase)

			var err error
			var responseTask *model.T_Task
			var responseTasks []model.T_Task
			switch scenario.callMethod {
			case "CreateTask":
				responseTask, err = taskService.CreateTask(context, scenario.mockInputTask)
			case "UpdateTaskByID":
				responseTask, err = taskService.UpdateTaskByID(context, scenario.mockInputTask)
			case "GetTaskByID":
				responseTask, err = taskService.GetTaskByID(context, scenario.mockInputInteger)
			case "DeleteTaskByID":
				err = taskService.DeleteTaskByID(context, scenario.mockInputInteger)
			case "ListTasks":
				responseTasks, err = taskService.ListTasks(context)
			}

			switch scenario.callMethod {
			case "CreateTask", "UpdateTaskByID", "GetTaskByID":
				if scenario.mockReturnError != nil {
					se := new(errors.Error)
					requires.True(errors.As(err, &se))
					requires.Equal(*scenario.expectedError, *se)
				} else {
					requires.Nil(err)
					requires.Equal(*scenario.mockReturnTask, *responseTask)
				}
			case "DeleteTaskByID":
				if scenario.mockReturnError != nil {
					se := new(errors.Error)
					requires.True(errors.As(err, &se))
					requires.Equal(*scenario.expectedError, *se)
				} else {
					requires.Nil(err)
				}
			case "ListTasks":
				if scenario.mockReturnError != nil {
					se := new(errors.Error)
					requires.True(errors.As(err, &se))
					requires.Equal(*scenario.expectedError, *se)
				} else {
					requires.Nil(err)
					requires.Equal(scenario.mockReturnTasks, responseTasks)
				}
			}
		})
	}
}
