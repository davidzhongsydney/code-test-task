package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/log"
	"qantas.com/task/internal/conf"
	"qantas.com/task/internal/service"
	"qantas.com/task/model"
)

type HTTPServer struct {
	router *chi.Mux
}

func NewHTTPServer(c *conf.Server, taskSvc *service.TaskService, logger log.Logger) Server {

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// r.Use(middleware.Timeout(c.Http.Timeout.AsDuration()))

	ctx := context.Background()

	r.Get("/ListingTasks", func(w http.ResponseWriter, r *http.Request) {
		result, err := taskSvc.ListTasks(ctx)

		if err != nil {

		}

		json.NewEncoder(w).Encode(result)
	})

	r.Post("/CreateTask", func(w http.ResponseWriter, r *http.Request) {
		var task model.Task
		json.NewDecoder(r.Body).Decode(&task)
		result, err := taskSvc.CreateTask(ctx, &task)

		if err != nil {

		}

		json.NewEncoder(w).Encode(result)
	})

	return &HTTPServer{router: r}
}

func (s *HTTPServer) Run() error {
	err := http.ListenAndServe(":8000", s.router)

	if err != nil {
		return err
	}
	return nil
}
