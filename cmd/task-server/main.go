package main

import (
	"fmt"
	"log"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	conf "qantas.com/task/internal/conf"
	model "qantas.com/task/model"
)

var path = "../../configs/config.yaml"

func main() {
	var task model.Task
	fmt.Println(task.Name)

	c := config.New(
		config.WithSource( // 初始化配置源
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
	// Get the corresponding value
	// name, err := c.Value("server.http.timeout").String()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(name)
}
