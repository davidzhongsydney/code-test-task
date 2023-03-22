package main

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/suite"
	"qantas.com/task/internal/biz"
	conf "qantas.com/task/internal/conf"
	"qantas.com/task/internal/encoder"
	"qantas.com/task/internal/server"
	"qantas.com/task/model"
	"qantas.com/task/utils"
)

type _HTTPSuccess_Task struct {
	Code int          `json:"code,omitempty"`
	Data model.T_Task `json:"data,omitempty"`
}

type _HTTPSuccess_Tasks struct {
	Code int            `json:"code,omitempty"`
	Data []model.T_Task `json:"data,omitempty"`
}

// type _HTTPSuccess struct {
// 	Code int `json:"code,omitempty"`
// }

func (s *IntegrationTestSuite) SetupSuite() {
	c := config.New(
		config.WithSource(
			file.NewSource("../../configs"),
		),
	)
	if err := c.Load(); err != nil {
		s.T().Fatalf("failed to load config file. Error: %s", err.Error())
	}

	var bc conf.Bootstrap

	if err := c.Scan(&bc); err != nil {
		s.T().Fatalf("failed to scan config file to conf.Bootstrap. Error: %s", err.Error())
	}
	logger := log.With(log.NewStdLogger(os.Stdout))

	s.context = context.Background()
	app, cleanup, err := wireApp(bc.Server, bc.Data, logger, s.context)
	if err != nil {
		s.T().Fatalf("failed to wire app. Error: %s", err.Error())
	}

	defer cleanup()

	httpServer := app.(*server.HTTPServer)
	s.testServer = httptest.NewServer(httpServer.GetRouter())
	httpHandler := httpServer.GetHttpHandler()
	taskHttpHandler := httpHandler.(*server.TasksHTTPHandler)
	taskService := taskHttpHandler.GetTaskService()
	s.uc = taskService.GetTaskUsecase()

	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("failed to run app. Error: %s", err.Error())
		}
	}()
}

func (s *IntegrationTestSuite) TearDownTest() {
	s.uc.ClearTasks(s.context)
}

type IntegrationTestSuite struct {
	suite.Suite
	context    context.Context
	testServer *httptest.Server
	uc         *biz.TaskUsecase
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) Test_CreateTask_Success() {
	// Create a task "POST", "/task"
	t := model.Task{Name: "user1", Content: "content1"}
	taskJson, err := json.Marshal(t)
	s.Require().Nil(err)

	startTime := time.Now()
	_, resp := utils.TestRequest(s.T(), s.testServer, "POST", "/task", strings.NewReader(string(taskJson)))
	endTime := time.Now()

	rt := _HTTPSuccess_Task{}

	err = json.Unmarshal([]byte(resp), &rt)
	s.Require().Nil(err)

	// Verify the output
	creationTime := rt.Data.CreatedAt

	s.Require().Equal(200, rt.Code)
	s.Require().True((creationTime.After(startTime) && creationTime.Before(endTime) || creationTime.Equal(startTime) || creationTime.Equal(endTime)))
	s.Require().Nil(rt.Data.UpdatedAt)
	s.Require().Nil(rt.Data.DeletedAt)
	s.Require().Equal(uint64(1), rt.Data.TaskID)
	s.Require().Equal("user1", rt.Data.Name)
	s.Require().Equal("content1", rt.Data.Content)
}

func (s *IntegrationTestSuite) Test_GetTask_Success() {
	// Create a task "POST", "/task"
	t := model.Task{Name: "user1", Content: "content2"}
	taskJson, err := json.Marshal(t)
	s.Require().Nil(err)

	_, resp := utils.TestRequest(s.T(), s.testServer, "POST", "/task", strings.NewReader(string(taskJson)))

	rt := _HTTPSuccess_Task{}

	err = json.Unmarshal([]byte(resp), &rt)
	s.Require().Nil(err)

	// Get task by id "GET", "/task/1"
	_, resp = utils.TestRequest(s.T(), s.testServer, "GET", "/task/1", nil)

	gt := _HTTPSuccess_Task{}
	err = json.Unmarshal([]byte(resp), &gt)
	s.Require().Nil(err)

	// Verify the output
	s.Require().Equal(200, rt.Code)
	s.Require().Equal(rt, gt)
}

func (s *IntegrationTestSuite) Test_GetTask_TaskNotFound() {
	// Get task by id "GET", "/task/1"
	_, resp := utils.TestRequest(s.T(), s.testServer, "GET", "/task/1", nil)

	// Error returned from "GET", "/task/1"
	actualError := encoder.HTTPError{}
	err := json.Unmarshal([]byte(resp), &actualError)
	s.Require().Nil(err)

	// Expected error
	er := model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST))
	expectedError := *encoder.FromError(er)

	// Verify the error
	s.Require().Equal(expectedError, actualError)
}

func (s *IntegrationTestSuite) Test_UpdateTask_Success() {
	// Create a task
	t1 := model.Task{Name: "user1", Content: "content1"}
	taskJson, err := json.Marshal(t1)
	s.Require().Nil(err)
	_, resp := utils.TestRequest(s.T(), s.testServer, "POST", "/task", strings.NewReader(string(taskJson)))
	ct := _HTTPSuccess_Task{}

	err = json.Unmarshal([]byte(resp), &ct)
	s.Require().Equal(200, ct.Code)
	s.Require().Nil(err)

	// Update a task "PUT", "/task"
	t2 := model.Task{TaskID: ct.Data.TaskID, Name: "user2", Content: "content2"}
	taskJson, err = json.Marshal(t2)
	s.Require().Nil(err)

	startTime := time.Now()
	_, resp = utils.TestRequest(s.T(), s.testServer, "PUT", "/task", strings.NewReader(string(taskJson)))
	endTime := time.Now()

	rt := _HTTPSuccess_Task{}
	err = json.Unmarshal([]byte(resp), &rt)
	s.Require().Nil(err)

	// Verify the updated task
	updateTime := rt.Data.UpdatedAt

	s.Require().Equal(200, rt.Code)
	s.Require().Equal(ct.Data.TaskID, rt.Data.TaskID)
	s.Require().Equal("user2", rt.Data.Name)
	s.Require().Equal("content2", rt.Data.Content)
	s.Require().Equal(ct.Data.CreatedAt, rt.Data.CreatedAt)
	s.Require().True((updateTime.After(startTime) && updateTime.Before(endTime) || updateTime.Equal(startTime) || updateTime.Equal(endTime)))
	s.Require().Nil(rt.Data.DeletedAt)
}

func (s *IntegrationTestSuite) Test_UpdateTask_TaskNotFound() {
	// Update a task "PUT", "/task"
	t := model.Task{TaskID: 1, Name: "user2", Content: "content2"}
	taskJson, err := json.Marshal(t)
	s.Require().Nil(err)

	_, resp := utils.TestRequest(s.T(), s.testServer, "PUT", "/task", strings.NewReader(string(taskJson)))

	// Error returned from "PUT", "/task"
	actualError := encoder.HTTPError{}
	err = json.Unmarshal([]byte(resp), &actualError)
	s.Require().Nil(err)

	// Expected error
	er := model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST))
	expectedError := *encoder.FromError(er)

	// Verify the error
	s.Require().Equal(expectedError, actualError)
}

func (s *IntegrationTestSuite) Test_UpdateTask_TaskIdNotSpecified() {

	// Update a task "PUT", "/task"
	t := model.Task{Name: "user2", Content: "content2"}
	taskJson, err := json.Marshal(t)
	s.Require().Nil(err)

	_, resp := utils.TestRequest(s.T(), s.testServer, "PUT", "/task", strings.NewReader(string(taskJson)))

	// Error returned from "PUT", "/task"
	actualError := encoder.HTTPError{}
	err = json.Unmarshal([]byte(resp), &actualError)
	s.Require().Nil(err)

	// Expected error
	er := model.ErrorTaskIdUnspecified(string(encoder.TASK_ID_NOT_SPECIFIED))
	expectedError := *encoder.FromError(er)

	// Verify the error
	s.Require().Equal(expectedError, actualError)
}

func (s *IntegrationTestSuite) Test_DeleteTask_Success() {
	// Create a task
	t1 := model.Task{Name: "user1", Content: "content1"}
	taskJson, err := json.Marshal(t1)
	s.Require().Nil(err)
	_, resp := utils.TestRequest(s.T(), s.testServer, "POST", "/task", strings.NewReader(string(taskJson)))
	ct := _HTTPSuccess_Task{}

	err = json.Unmarshal([]byte(resp), &ct)
	s.Require().Equal(200, ct.Code)
	s.Require().Nil(err)

	// Delete a task "DELETE", "/task/1"
	_, resp = utils.TestRequest(s.T(), s.testServer, "DELETE", "/task/1", nil)

	actual := encoder.HTTPSuccess{}
	err = json.Unmarshal([]byte(resp), &actual)
	s.Require().Nil(err)

	// Verify the updated task
	expected := *encoder.FromResponse(nil)
	s.Require().Equal(expected, actual)
}

func (s *IntegrationTestSuite) Test_DeleteTask_TaskNotFound() {
	// Delete a task "DELETE", "/task/1"
	_, resp := utils.TestRequest(s.T(), s.testServer, "DELETE", "/task/1", nil)

	// Error returned from "DELETE", "/task/1"
	actualError := encoder.HTTPError{}
	err := json.Unmarshal([]byte(resp), &actualError)
	s.Require().Nil(err)

	// Expected error
	er := model.ErrorTaskNotFound(string(encoder.TASK_NOT_EXIST))
	expectedError := *encoder.FromError(er)

	// Verify the updated task
	s.Require().Equal(expectedError, actualError)
}

func (s *IntegrationTestSuite) Test_ListTask_Success() {
	// Create first task "POST", "/task"
	t1 := model.Task{Name: "user1", Content: "content1"}
	taskJson, err := json.Marshal(t1)
	s.Require().Nil(err)

	_, resp := utils.TestRequest(s.T(), s.testServer, "POST", "/task", strings.NewReader(string(taskJson)))

	rt1 := _HTTPSuccess_Task{}
	err = json.Unmarshal([]byte(resp), &rt1)
	s.Require().Nil(err)

	// Create second task "POST", "/task"
	t2 := model.Task{Name: "user2", Content: "content2"}
	taskJson, err = json.Marshal(t2)
	s.Require().Nil(err)

	_, resp = utils.TestRequest(s.T(), s.testServer, "POST", "/task", strings.NewReader(string(taskJson)))

	rt2 := _HTTPSuccess_Task{}
	err = json.Unmarshal([]byte(resp), &rt2)
	s.Require().Nil(err)

	// Query task list
	_, resp = utils.TestRequest(s.T(), s.testServer, "GET", "/tasks", nil)
	rts := _HTTPSuccess_Tasks{}
	err = json.Unmarshal([]byte(resp), &rts)
	s.Require().Nil(err)

	// Verify the output
	s.Require().Equal(200, rts.Code)

	if rts.Data[0].TaskID == rt1.Data.TaskID {
		s.Require().Equal(rt1.Data, rts.Data[0])
		s.Require().Equal(rt2.Data, rts.Data[1])
	} else {
		s.Require().Equal(rt1.Data, rts.Data[1])
		s.Require().Equal(rt2.Data, rts.Data[0])
	}
}
