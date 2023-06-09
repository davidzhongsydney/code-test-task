package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"qantas.com/task/internal/biz"
	"qantas.com/task/internal/encoder"
	"qantas.com/task/internal/server"
	"qantas.com/task/internal/service"
	"qantas.com/task/mocks"
	"qantas.com/task/model"
	"qantas.com/task/utils"
)

var timeLayout = "2006-01-02T15:04:05.000Z"
var time1, _ = time.Parse(timeLayout, "2014-11-12T11:45:26.371Z")
var time2, _ = time.Parse(timeLayout, "2015-11-12T11:45:26.371Z")

func TestHTTPHandler(t *testing.T) {

	type testCases struct {
		description     string
		mockMethod      string
		callMethod      string
		mockInputTask   *model.Task
		mockReturnTask  *model.T_Task
		mockReturnTasks []model.T_Task
		mockReturnError *errors.Error
		url             string
		httpMethod      string
		expectedOutput  string
	}

	for _, scenario := range []testCases{
		{
			description: "get task by id success",
			mockMethod:  "Get",
			url:         "/task/1",
			mockReturnTask: &model.T_Task{Task: model.Task{TaskID: 2, Name: "david", Content: "content text"},
				T_Internal: model.T_Internal{CreatedAt: &time1}},
			callMethod:     "GetTaskByIdHTTPHandler",
			httpMethod:     "GET",
			expectedOutput: "{\"code\":200,\"data\":{\"taskID\":2,\"name\":\"david\",\"content\":\"content text\",\"createdAt\":\"2014-11-12T11:45:26.371Z\"}}\n",
		},
		{
			description:     "get task by id failed - task not found",
			mockMethod:      "Get",
			url:             "/task/1",
			mockReturnError: model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
			callMethod:      "GetTaskByIdHTTPHandler",
			httpMethod:      "GET",
			expectedOutput:  "{\"code\":404,\"errors\":{\"TASK_NOT_FOUND\":\"task does not exist\"}}\n",
		},
		{
			description:   "create task success",
			mockMethod:    "Create",
			url:           "/task",
			mockInputTask: &model.Task{Name: "david", Content: "content text"},
			mockReturnTask: &model.T_Task{Task: model.Task{TaskID: 2, Name: "david", Content: "content text"},
				T_Internal: model.T_Internal{CreatedAt: &time1}},
			callMethod:     "CreateTaskHTTPHandler",
			httpMethod:     "POST",
			expectedOutput: "{\"code\":200,\"data\":{\"taskID\":2,\"name\":\"david\",\"content\":\"content text\",\"createdAt\":\"2014-11-12T11:45:26.371Z\"}}\n",
		},
		{
			description:     "create task failed - creation error",
			mockMethod:      "Create",
			url:             "/task",
			mockInputTask:   &model.Task{Name: "user", Content: "content"},
			mockReturnError: model.ErrorTaskCreationError(string(encoder.TASK_CREATION_ERROR)),
			callMethod:      "CreateTaskHTTPHandler",
			httpMethod:      "POST",
			expectedOutput:  "{\"code\":500,\"errors\":{\"TASK_CREATION_ERROR\":\"task is failed to be created\"}}\n",
		},
		{
			description:   "update task by id success",
			mockMethod:    "Update",
			url:           "/task",
			mockInputTask: &model.Task{TaskID: 2, Name: "david", Content: "content text"},
			mockReturnTask: &model.T_Task{Task: model.Task{TaskID: 2, Name: "john", Content: "content text updated"},
				T_Internal: model.T_Internal{CreatedAt: &time1, UpdatedAt: &time2}},
			callMethod:     "UpdateTaskByIdHTTPHandler",
			httpMethod:     "PUT",
			expectedOutput: "{\"code\":200,\"data\":{\"taskID\":2,\"name\":\"john\",\"content\":\"content text updated\",\"createdAt\":\"2014-11-12T11:45:26.371Z\",\"updatedAt\":\"2015-11-12T11:45:26.371Z\"}}\n",
		},
		{
			description:     "update task by id failed - task not found",
			mockMethod:      "Update",
			url:             "/task",
			mockInputTask:   &model.Task{TaskID: 2, Name: "david", Content: "content text"},
			mockReturnError: model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
			callMethod:      "UpdateTaskByIdHTTPHandler",
			httpMethod:      "PUT",
			expectedOutput:  "{\"code\":404,\"errors\":{\"TASK_NOT_FOUND\":\"task does not exist\"}}\n",
		},
		{
			description:     "update task by id failed - task id not specified",
			mockMethod:      "Update",
			url:             "/task",
			mockInputTask:   &model.Task{Name: "david", Content: "content text"},
			mockReturnError: model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
			callMethod:      "UpdateTaskByIdHTTPHandler",
			httpMethod:      "PUT",
			expectedOutput:  "{\"code\":400,\"errors\":{\"TASK_ID_UNSPECIFIED\":\"task id not specified\"}}\n",
		},
		{
			description:    "delete task by id success",
			mockMethod:     "Delete",
			url:            "/task/2",
			callMethod:     "DeleteTaskByIdHTTPHandler",
			httpMethod:     "DELETE",
			expectedOutput: "{\"code\":200}\n",
		},
		{
			description:     "delete task by id failed - task not found",
			mockMethod:      "Delete",
			url:             "/task/2",
			mockReturnError: model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST)),
			callMethod:      "DeleteTaskByIdHTTPHandler",
			httpMethod:      "DELETE",
			expectedOutput:  "{\"code\":404,\"errors\":{\"TASK_NOT_FOUND\":\"task does not exist\"}}\n",
		},
		{
			description: "list tasks success",
			mockMethod:  "List",
			url:         "/tasks",
			mockReturnTasks: []model.T_Task{
				{Task: model.Task{TaskID: 1, Name: "user1", Content: "content1"},
					T_Internal: model.T_Internal{CreatedAt: &time1}},
				{Task: model.Task{TaskID: 2, Name: "user2", Content: "content2"},
					T_Internal: model.T_Internal{CreatedAt: &time2}}},
			callMethod:     "ListTasksHTTPHandler",
			httpMethod:     "GET",
			expectedOutput: "{\"code\":200,\"data\":[{\"taskID\":1,\"name\":\"user1\",\"content\":\"content1\",\"createdAt\":\"2014-11-12T11:45:26.371Z\"},{\"taskID\":2,\"name\":\"user2\",\"content\":\"content2\",\"createdAt\":\"2015-11-12T11:45:26.371Z\"}]}\n",
		},
		{
			description:     "list tasks failed - database timeout",
			mockMethod:      "List",
			url:             "/tasks",
			mockReturnError: model.ErrorTaskDbTimeout(string(encoder.TASK_DATABASE_TIMEOUT)),
			callMethod:      "ListTasksHTTPHandler",
			httpMethod:      "GET",
			expectedOutput:  "{\"code\":500,\"errors\":{\"TASK_DB_TIMEOUT\":\"task database timeout\"}}\n",
		},
	} {
		requires := require.New(t)
		context := context.Background()
		logger := log.With(log.NewStdLogger(os.Stdout))
		taskRepoMock := mocks.TaskRepo{}

		// Set up database method mock
		switch scenario.callMethod {
		case "CreateTaskHTTPHandler", "UpdateTaskByIdHTTPHandler", "GetTaskByIdHTTPHandler":
			if scenario.mockReturnError != nil {
				taskRepoMock.On(scenario.mockMethod, mock.Anything, mock.Anything).Return(nil, scenario.mockReturnError)
			} else {
				taskRepoMock.On(scenario.mockMethod, mock.Anything, mock.Anything).Return(scenario.mockReturnTask, nil)
			}
		case "DeleteTaskByIdHTTPHandler":
			if scenario.mockReturnError != nil {
				taskRepoMock.On(scenario.mockMethod, mock.Anything, mock.Anything).Return(scenario.mockReturnError)
			} else {
				taskRepoMock.On(scenario.mockMethod, mock.Anything, mock.Anything).Return(nil)
			}
		case "ListTasksHTTPHandler":
			if scenario.mockReturnError != nil {
				taskRepoMock.On(scenario.mockMethod, mock.Anything).Return(nil, scenario.mockReturnError)
			} else {
				taskRepoMock.On(scenario.mockMethod, mock.Anything).Return(scenario.mockReturnTasks, nil)
			}
		}

		taskUseCase := biz.NewTaskUsecase(&taskRepoMock, logger)
		taskService := service.NewTaskService(taskUseCase, logger)

		// Set up router
		r := chi.NewRouter()

		httpHandler := server.NewTaskHTTPHandler(taskService, logger, context)

		r.Get("/tasks", httpHandler.ListTasksHTTPHandler()) // GET /tasks - Get a list of tasks.
		r.Route("/task", func(r chi.Router) {
			r.Get("/{id:[0-9]+}", httpHandler.GetTaskByIdHTTPHandler())       // GET      /task/{id} - Get a task by id.
			r.Post("/", httpHandler.CreateTaskHTTPHandler())                  // POST     /task      - Create a new task.
			r.Put("/", httpHandler.UpdateTaskByIdHTTPHandler())               // PUT      /task      - Update a new task by id.
			r.Delete("/{id:[0-9]+}", httpHandler.DeleteTaskByIdHTTPHandler()) // DELETE   /task/{id} - Delete a task by id.
		})
		ts := httptest.NewServer(r)
		defer ts.Close()

		// Send request
		var resp string

		switch scenario.callMethod {
		case "GetTaskByIdHTTPHandler", "DeleteTaskByIdHTTPHandler", "ListTasksHTTPHandler":
			_, resp = utils.TestRequest(t, ts, scenario.httpMethod, scenario.url, nil)
		case "CreateTaskHTTPHandler", "UpdateTaskByIdHTTPHandler":
			byteArray, err := json.Marshal(*scenario.mockInputTask)
			if err != nil {
				t.Fatal(err)
			}
			_, resp = utils.TestRequest(t, ts, scenario.httpMethod, scenario.url, bytes.NewBuffer(byteArray))
		}

		// Verify response
		requires.Equal(scenario.expectedOutput, resp)
	}
}
