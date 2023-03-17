package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	conf "qantas.com/task/internal/conf"
)

var (
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func main() {

	flag.Parse()

	fmt.Println("flagconf: ", flagconf)

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)

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
