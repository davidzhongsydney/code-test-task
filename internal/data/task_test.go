package data_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	errors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
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

type DataSourceTestSuite struct {
	suite.Suite
	taskRepo biz.ITaskRepo
	context  context.Context
}

func (s *DataSourceTestSuite) SetupSuite() {
	logger := log.With(log.NewStdLogger(os.Stdout))
	conf := conf.Data{}
	dataRepo, _, err := data.NewData(&conf, logger)

	if err != nil {
		s.FailNow("unable to connect to database.", err.Error())
	}

	s.taskRepo = data.NewTaskRepo(dataRepo, logger)
	s.context = context.Background()
}

func (s *DataSourceTestSuite) TearDownTest() {
	s.taskRepo.Empty(s.context)
}

func TestTaskSuite(t *testing.T) {
	suite.Run(t, new(DataSourceTestSuite))
}

func (s *DataSourceTestSuite) Test_AddTask() {
	// Create first task
	t1 := model.Task{Name: "user name 1", Content: "content text 1"}

	startTime := time.Now()
	st, err := s.taskRepo.Create(s.context, &t1)
	endTime := time.Now()

	s.Require().Nil(err)

	// Verify the first tasked task return value
	creationTime := st.CreatedAt

	s.Require().True((creationTime.After(startTime) && creationTime.Before(endTime) || creationTime.Equal(startTime) || creationTime.Equal(endTime)))
	s.Require().Nil(st.UpdatedAt)
	s.Require().Nil(st.DeletedAt)

	s.Require().Equal(uint64(1), st.TaskID)
	s.Require().Equal("user name 1", st.Name)
	s.Require().Equal("content text 1", st.Content)

	// Create second task
	t2 := model.Task{TaskID: 5, Name: "user name 2", Content: "content text 2"}
	st, err = s.taskRepo.Create(s.context, &t2)
	s.Require().Nil(err)

	// Verify the second tasked task return value
	s.Require().Equal(uint64(2), st.TaskID)
	s.Require().Equal("user name 2", st.Name)
	s.Require().Equal("content text 2", st.Content)
}

func (s *DataSourceTestSuite) Test_GetTask_Success() {
	// Create task
	t1 := model.Task{Name: "user name 1", Content: "content text 1"}
	startTime := time.Now()
	_, err := s.taskRepo.Create(s.context, &t1)
	endTime := time.Now()
	s.Require().Nil(err)

	// Get task
	st, err := s.taskRepo.Get(s.context, 1)
	s.Require().Nil(err)

	// Verify the Get return value
	s.Require().Equal(uint64(1), st.TaskID)
	s.Require().Equal("user name 1", st.Name)
	s.Require().Equal("content text 1", st.Content)

	creationTime := st.CreatedAt
	s.Require().True((creationTime.After(startTime) && creationTime.Before(endTime) || creationTime.Equal(startTime) || creationTime.Equal(endTime)))
	s.Require().Nil(st.UpdatedAt)
	s.Require().Nil(st.DeletedAt)
}

func (s *DataSourceTestSuite) Test_GetTask_TaskNotFound() {
	se := new(errors.Error)

	// Query nonexistent task
	st, err := s.taskRepo.Get(s.context, 1)
	s.Require().True(errors.As(err, &se))
	s.Require().True(model.IsTaskNotFound(se))
	s.Require().Nil(st)

	// Query a deleted task
	t1 := model.Task{Name: "user name 1", Content: "content text 1"}
	st, err = s.taskRepo.Create(s.context, &t1)
	s.Require().Nil(err)

	err = s.taskRepo.Delete(s.context, st.TaskID)
	s.Require().Nil(err)

	st, err = s.taskRepo.Get(s.context, 1)
	s.Require().True(errors.As(err, &se))
	s.Require().True(model.IsTaskNotFound(se))
	s.Require().Nil(st)
}

func (s *DataSourceTestSuite) Test_UpdateTask_Success() {
	// Create task
	t := model.Task{Name: "user name 1", Content: "content text 1"}
	ct, err := s.taskRepo.Create(s.context, &t)
	s.Require().Nil(err)

	// Update a task
	tt := model.Task{TaskID: ct.TaskID, Name: "user name 2", Content: "content text 2"}
	startTime := time.Now()
	ut, err := s.taskRepo.Update(s.context, &tt)
	endTime := time.Now()

	s.Require().Nil(err)

	// Verify the updated return value
	s.Require().Equal(tt.TaskID, ut.TaskID)
	s.Require().Equal(tt.Name, ut.Name)
	s.Require().Equal(tt.Content, ut.Content)
	s.Require().Equal(ct.CreatedAt, ut.CreatedAt)
	updatedTime := ut.UpdatedAt
	s.Require().True((updatedTime.After(startTime) && updatedTime.Before(endTime) || updatedTime.Equal(startTime) || updatedTime.Equal(endTime)))
	s.Require().Nil(ut.DeletedAt)
}

func (s *DataSourceTestSuite) Test_UpdateTask_TaskNotFound() {
	se := new(errors.Error)

	// Update nonexistent task
	tt := model.Task{TaskID: 0, Name: "user name 2", Content: "content text 2"}
	ut, err := s.taskRepo.Update(s.context, &tt)
	s.Require().Nil(ut)
	s.Require().True(errors.As(err, &se))
	s.Require().True(model.IsTaskNotFound(se))

	// Update a deleted task
	t1 := model.Task{Name: "user name 1", Content: "content text 1"}
	_, err = s.taskRepo.Create(s.context, &t1)
	s.Require().Nil(err)

	err = s.taskRepo.Delete(s.context, 1)
	s.Require().Nil(err)

	tt = model.Task{TaskID: 1, Name: "user name 2", Content: "content text 2"}
	ut, err = s.taskRepo.Update(s.context, &tt)
	s.Require().Nil(ut)
	s.Require().True(errors.As(err, &se))
	s.Require().True(model.IsTaskNotFound(se))
}

func (s *DataSourceTestSuite) Test_DeleteTask_Success() {
	// Create task
	t1 := model.Task{Name: "user name 1", Content: "content text 1"}
	_, err := s.taskRepo.Create(s.context, &t1)
	s.Require().Nil(err)

	// Delete a task
	err = s.taskRepo.Delete(s.context, 1)
	s.Require().Nil(err)
}

func (s *DataSourceTestSuite) Test_DeleteTask_TaskNotFound() {
	se := new(errors.Error)

	// Delete nonexistent task
	err := s.taskRepo.Delete(s.context, 1)
	s.Require().True(errors.As(err, &se))
	s.Require().True(model.IsTaskNotFound(se))

	// Delete a deleted task
	t1 := model.Task{Name: "user name 1", Content: "content text 1"}
	_, err = s.taskRepo.Create(s.context, &t1)
	s.Require().Nil(err)

	err = s.taskRepo.Delete(s.context, 1)
	s.Require().Nil(err)

	err = s.taskRepo.Delete(s.context, 1)
	s.Require().True(errors.As(err, &se))
	s.Require().True(model.IsTaskNotFound(se))
}

func (s *DataSourceTestSuite) Test_ListTask() {

	// Add first task
	at := model.Task{Name: "user 1", Content: "content text 1"}
	ctask1, err := s.taskRepo.Create(s.context, &at)
	s.Require().Nil(err)

	// Add second task, and update it value
	at = model.Task{Name: "user 2", Content: "content text 2"}
	_, err = s.taskRepo.Create(s.context, &at)
	s.Require().Nil(err)

	ut := model.Task{TaskID: 2, Content: "content text 2 updated"}
	ctask2, err := s.taskRepo.Update(s.context, &ut)
	s.Require().Nil(err)

	// Add third task, and delete
	at = model.Task{Name: "user 3", Content: "content text 3"}
	_, err = s.taskRepo.Create(s.context, &at)
	s.Require().Nil(err)

	err = s.taskRepo.Delete(s.context, 3)
	s.Require().Nil(err)

	// Get the list of tasks
	listTask, err := s.taskRepo.List(s.context)
	s.Require().Nil(err)

	// Verify the returned task list
	s.Require().Equal(2, len(listTask))

	if listTask[0].TaskID == 1 {
		s.Require().Equal(*ctask1, listTask[0])
		s.Require().Equal(*ctask2, listTask[1])
	} else {
		s.Require().Equal(*ctask1, listTask[1])
		s.Require().Equal(*ctask2, listTask[0])
	}
}
