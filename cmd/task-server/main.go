package main

import (
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	conf "qantas.com/task/internal/conf"
)

var path = "../../configs/config.yaml"

func main() {

	c := config.New(
		config.WithSource(
			file.NewSource(path),
		),
	)
	// load config source
	if err := c.Load(); err != nil {
		log.Fatal(err)
	}

	var bc conf.Bootstrap

	if err := c.Scan(&bc); err != nil {
		log.Fatal(err)
	}

	fmt.Println(bc.Server.Http.Timeout)

	logger := log.With(log.NewStdLogger(os.Stdout))

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}

	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
