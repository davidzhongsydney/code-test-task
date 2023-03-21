package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/log"
	"qantas.com/task/internal/conf"
	"qantas.com/task/internal/service"
)

type HTTPServer struct {
	router          *chi.Mux
	conf            *conf.Server
	taskHttpHandler ITaskHTTPHandler
}

func NewHTTPServer(c *conf.Server, logger log.Logger, httpHandler ITaskHTTPHandler) Server {

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(c.Http.Timeout.AsDuration()))

	// ctx := context.Background()
	// httpHandler := TasksHTTPHandler{TaskSvc: taskSvc, Ctx: ctx}

	r.Get("/tasks", httpHandler.ListTasksHTTPHandler()) // GET /tasks - Get a list of tasks.
	r.Route("/task", func(r chi.Router) {
		r.Get("/{id:[0-9]+}", httpHandler.GetTaskByIdHTTPHandler())       // GET      /task/{id} - Get a task by id.
		r.Post("/", httpHandler.CreateTaskHTTPHandler())                  // POST     /task      - Create a new task.
		r.Put("/", httpHandler.UpdateTaskByIdHTTPHandler())               // PUT      /task      - Update a new task by id.
		r.Delete("/{id:[0-9]+}", httpHandler.DeleteTaskByIdHTTPHandler()) // DELETE   /task/{id} - Delete a task by id.
	})

	return &HTTPServer{router: r, conf: c, taskHttpHandler: httpHandler}
}

func NewTaskHTTPHandler(taskSvc *service.TaskService, logger log.Logger, ctx context.Context) ITaskHTTPHandler {
	return &TasksHTTPHandler{taskSvc: taskSvc, ctx: ctx, log: log.NewHelper(logger)}
}

func (s *HTTPServer) Run() error {
	err := http.ListenAndServe(s.conf.Http.Addr, s.router)

	if err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) GetRouter() *chi.Mux {
	return s.router
}

func (s *HTTPServer) GetHttpHandler() ITaskHTTPHandler {
	return s.taskHttpHandler
}

type ITaskHTTPHandler interface {
	ListTasksHTTPHandler() http.HandlerFunc
	CreateTaskHTTPHandler() http.HandlerFunc
	GetTaskByIdHTTPHandler() http.HandlerFunc
	UpdateTaskByIdHTTPHandler() http.HandlerFunc
	DeleteTaskByIdHTTPHandler() http.HandlerFunc
}
