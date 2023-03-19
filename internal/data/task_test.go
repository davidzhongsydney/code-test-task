package data_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"qantas.com/task/internal/biz"
	"qantas.com/task/internal/conf"
	"qantas.com/task/internal/data"
	"qantas.com/task/model"
)

func Test_INT_NewTaskRepo_Success(t *testing.T) {
	requires := require.New(t)
	logger := log.With(log.NewStdLogger(os.Stdout))
	conf := conf.Data{}

	dataRepo, _, err := data.NewData(&conf, logger)

	if err != nil {
		t.Logf("unable to connect to database. Error %s ", err.Error())
		t.FailNow()
	}

	taskRepo := data.NewTaskRepo(dataRepo, logger)
	requires.Equal("*data.taskRepo", fmt.Sprint(reflect.TypeOf(taskRepo)))
}

type INT_TaskTestSuite struct {
	suite.Suite
	taskRepo biz.TaskRepo
	context  context.Context
}

func (s *INT_TaskTestSuite) SetupSuite() {
	logger := log.With(log.NewStdLogger(os.Stdout))
	conf := conf.Data{}
	dataRepo, _, err := data.NewData(&conf, logger)

	if err != nil {
		s.FailNow("unable to connect to database.", err.Error())
	}

	s.taskRepo = data.NewTaskRepo(dataRepo, logger)
	s.context = context.Background()
}

func (s *INT_TaskTestSuite) TearDownTest() {
	s.taskRepo.Empty(s.context)
}

func TestTaskSuite(t *testing.T) {
	suite.Run(t, new(INT_TaskTestSuite))
}

func (s *INT_TaskTestSuite) Test_AddTask() {
	// Create first task
	t1 := model.Task{Name: "user name 1", Content: "content text 1"}

	startTime := time.Now()
	st, err := s.taskRepo.Create(s.context, &t1)
	endTime := time.Now()

	if err != nil {
		s.FailNow("unable to create a task in database.", err.Error())
	}

	creationTime := st.CreatedAt.AsTime()
	var emptyTimeStamp *timestamp.Timestamp

	s.Require().True((creationTime.After(startTime) && creationTime.Before(endTime) || creationTime.Equal(startTime) || creationTime.Equal(endTime)))
	s.Require().Equal(emptyTimeStamp, st.UpdatedAt)
	s.Require().Equal(emptyTimeStamp, st.DeletedAt)
	s.Require().Equal(uint64(1), st.TaskID)
	s.Require().Equal("user name 1", st.Name)
	s.Require().Equal("content text 1", st.Content)

	// Create second task
	t2 := model.Task{TaskID: 5, Name: "user name 2", Content: "content text 2"}
	st, err = s.taskRepo.Create(s.context, &t2)
	if err != nil {
		s.FailNow("unable to create a task in database.", err.Error())
	}

	s.Require().Equal(uint64(2), st.TaskID)
	s.Require().Equal("user name 2", st.Name)
	s.Require().Equal("content text 2", st.Content)
}

func (s *INT_TaskTestSuite) Test_GetTask_Success() {
	// Create task
	t1 := model.Task{Name: "user name 1", Content: "content text 1"}
	startTime := time.Now()
	_, err := s.taskRepo.Create(s.context, &t1)
	endTime := time.Now()
	if err != nil {
		s.FailNow("unable to create a task in database.", err.Error())
	}

	// Get task
	st, err := s.taskRepo.Get(s.context, 1)
	s.Require().Empty(err)

	s.Require().Equal(uint64(1), st.TaskID)
	s.Require().Equal("user name 1", st.Name)
	s.Require().Equal("content text 1", st.Content)

	var emptyTimeStamp *timestamp.Timestamp
	creationTime := st.CreatedAt.AsTime()
	s.Require().True((creationTime.After(startTime) && creationTime.Before(endTime) || creationTime.Equal(startTime) || creationTime.Equal(endTime)))
	s.Require().Equal(emptyTimeStamp, st.UpdatedAt)
	s.Require().Equal(emptyTimeStamp, st.DeletedAt)
}

func (s *INT_TaskTestSuite) Test_GetTask_TaskNotFound() {
	// Query nonexistent task
	st, err := s.taskRepo.Get(s.context, 1)
	s.Require().EqualError(err, "error: code = 404 reason = TASK_NOT_FOUND message = task does not exist metadata = map[] cause = <nil>")
	s.Require().Empty(st)

	// Query a deleted task
	t1 := model.Task{Name: "user name 1", Content: "content text 1"}
	st, err = s.taskRepo.Create(s.context, &t1)
	s.Require().Empty(err)

	err = s.taskRepo.Delete(s.context, st.TaskID)
	s.Require().Empty(err)

	st, err = s.taskRepo.Get(s.context, 1)
	s.Require().EqualError(err, "error: code = 404 reason = TASK_NOT_FOUND message = task has been logically deleted metadata = map[] cause = <nil>")
	s.Require().Empty(st)
}

func (s *INT_TaskTestSuite) Test_UpdateTask_Success() {
	// Create task
	t := model.Task{Name: "user name 1", Content: "content text 1"}
	st, err := s.taskRepo.Create(s.context, &t)
	if err != nil {
		s.FailNow("unable to create a task in database.", err.Error())
	}

	// Update a task
	tt := model.Task{TaskID: 1, Name: "user name 2", Content: "content text 2"}
	startTime := time.Now()
	ut, err := s.taskRepo.Update(s.context, &tt)
	endTime := time.Now()

	s.Require().Empty(err)
	s.Require().Equal(tt.TaskID, ut.TaskID)
	s.Require().Equal(tt.Name, ut.Name)
	s.Require().Equal(tt.Content, ut.Content)
	s.Require().Equal(st.CreatedAt, ut.CreatedAt)
	updatedTime := ut.UpdatedAt.AsTime()
	s.Require().True((updatedTime.After(startTime) && updatedTime.Before(endTime) || updatedTime.Equal(startTime) || updatedTime.Equal(endTime)))

	var emptyTimeStamp *timestamp.Timestamp
	s.Require().Equal(emptyTimeStamp, ut.DeletedAt)
}

func (s *INT_TaskTestSuite) Test_UpdateTask_TaskNotFound() {
	// Update nonexistent task
	tt := model.Task{TaskID: 0, Name: "user name 2", Content: "content text 2"}
	ut, err := s.taskRepo.Update(s.context, &tt)
	s.Require().Empty(ut)
	s.Require().EqualError(err, "error: code = 404 reason = TASK_NOT_FOUND message = task does not exist metadata = map[] cause = <nil>")

	// Update a deleted task
	t1 := model.Task{Name: "user name 1", Content: "content text 1"}
	_, err = s.taskRepo.Create(s.context, &t1)
	s.Require().Empty(err)

	err = s.taskRepo.Delete(s.context, 1)
	s.Require().Empty(err)

	tt = model.Task{TaskID: 1, Name: "user name 2", Content: "content text 2"}
	ut, err = s.taskRepo.Update(s.context, &tt)
	s.Require().Empty(ut)
	s.Require().EqualError(err, "error: code = 404 reason = TASK_NOT_FOUND message = task has been logically deleted metadata = map[] cause = <nil>")
}

func (s *INT_TaskTestSuite) Test_DeleteTask_Success() {
	// Create task
	t1 := model.Task{Name: "user name 1", Content: "content text 1"}
	_, err := s.taskRepo.Create(s.context, &t1)
	if err != nil {
		s.FailNow("unable to create a task in database.", err.Error())
	}

	// Delete a task
	err = s.taskRepo.Delete(s.context, 1)
	s.Require().Empty(err)
}

func (s *INT_TaskTestSuite) Test_DeleteTask_TaskNotFound() {
	// Delete nonexistent task
	err := s.taskRepo.Delete(s.context, 1)
	s.Require().EqualError(err, "error: code = 404 reason = TASK_NOT_FOUND message = task does not exist metadata = map[] cause = <nil>")

	// Delete a deleted task
	t1 := model.Task{Name: "user name 1", Content: "content text 1"}
	_, err = s.taskRepo.Create(s.context, &t1)
	s.Require().Empty(err)

	err = s.taskRepo.Delete(s.context, 1)
	s.Require().Empty(err)

	err = s.taskRepo.Delete(s.context, 1)
	s.Require().EqualError(err, "error: code = 404 reason = TASK_NOT_FOUND message = task has been logically deleted metadata = map[] cause = <nil>")
}

func (s *INT_TaskTestSuite) Test_ListTask() {

	// Add first task
	at := model.Task{Name: "user 1", Content: "content text 1"}
	stask1, err := s.taskRepo.Create(s.context, &at)
	s.Require().Empty(err)

	// Add second task, and update
	at = model.Task{Name: "user 2", Content: "content text 2"}
	_, err = s.taskRepo.Create(s.context, &at)
	s.Require().Empty(err)
	ut := model.Task{TaskID: 2, Content: "content text 2 updated"}
	stask2, err := s.taskRepo.Update(s.context, &ut)
	s.Require().Empty(err)

	// Add third task, and delete
	at = model.Task{Name: "user 3", Content: "content text 3"}
	_, err = s.taskRepo.Create(s.context, &at)
	s.Require().Empty(err)
	err = s.taskRepo.Delete(s.context, 3)
	s.Require().Empty(err)

	listTask, err := s.taskRepo.List(s.context)
	s.Require().Empty(err)
	s.Require().Equal(2, len(listTask))

	if listTask[0].TaskID == 1 {
		s.Require().Equal(*stask1, listTask[0])
		s.Require().Equal(*stask2, listTask[1])
	} else {
		s.Require().Equal(*stask1, listTask[1])
		s.Require().Equal(*stask2, listTask[0])
	}
}
