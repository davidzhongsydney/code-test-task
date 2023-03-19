package data_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/require"
	"qantas.com/task/internal/conf"
	"qantas.com/task/internal/data"
)

func Test_UNIT_NewData_Success(t *testing.T) {
	requires := require.New(t)
	logger := log.With(log.NewStdLogger(os.Stdout))
	conf := conf.Data{}

	data, _, err := data.NewData(&conf, logger)
	if err != nil {
		t.Logf("unable to connect to database. Error %s ", err.Error())
		t.FailNow()
	}

	requires.Equal("*data.Data", fmt.Sprint(reflect.TypeOf(data)))
}
