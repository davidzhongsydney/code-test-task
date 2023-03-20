package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/log"
	"qantas.com/task/internal/conf"
	"qantas.com/task/internal/encoder"
	"qantas.com/task/internal/service"
	"qantas.com/task/model"
)

type HTTPServer struct {
	router *chi.Mux
	conf   *conf.Server
}

func NewHTTPServer(c *conf.Server, taskSvc *service.TaskService, logger log.Logger) Server {

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(c.Http.Timeout.AsDuration()))

	ctx := context.Background()

	httpHandler := TasksHTTPHandler{TaskSvc: taskSvc, Ctx: ctx}

	r.Get("/tasks", httpHandler.ListTasksHTTPHandler()) // GET /tasks - Get a list of tasks.
	r.Route("/task", func(r chi.Router) {
		r.Get("/{id:[0-9]+}", httpHandler.GetTaskByIdHTTPHandler())       // GET      /task/{id} - Get a task by id.
		r.Post("/", httpHandler.CreateTaskHTTPHandler())                  // POST     /task      - Create a new task.
		r.Put("/", httpHandler.UpdateTaskByIdHTTPHandler())               // PUT      /task      - Update a new task by id.
		r.Delete("/{id:[0-9]+}", httpHandler.DeleteTaskByIdHTTPHandler()) // DELETE   /task/{id} - Delete a task by id.
	})

	return &HTTPServer{router: r, conf: c}
}

func (s *HTTPServer) Run() error {
	err := http.ListenAndServe(s.conf.Http.Addr, s.router)

	if err != nil {
		return err
	}
	return nil
}

type TasksHTTPHandler struct {
	TaskSvc *service.TaskService
	Ctx     context.Context
}

func (h TasksHTTPHandler) ListTasksHTTPHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		result, err := h.TaskSvc.ListTasks(h.Ctx)

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			json.NewEncoder(w).Encode(encoder.FromError(err))
			return
		}

		json.NewEncoder(w).Encode(encoder.FromResponse(result))
	}
	return fn
}

func (h TasksHTTPHandler) CreateTaskHTTPHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var task model.Task
		json.NewDecoder(r.Body).Decode(&task)
		result, err := h.TaskSvc.CreateTask(h.Ctx, &task)

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			json.NewEncoder(w).Encode(encoder.FromError(err))
			return
		}

		json.NewEncoder(w).Encode(result)
	}
	return fn
}

func (h TasksHTTPHandler) GetTaskByIdHTTPHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)

		result, err := h.TaskSvc.GetTaskByID(h.Ctx, id)

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			json.NewEncoder(w).Encode(encoder.FromError(err))
			return
		}

		json.NewEncoder(w).Encode(encoder.FromResponse(result))
	}
	return fn
}

func (h TasksHTTPHandler) UpdateTaskByIdHTTPHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var task model.Task
		json.NewDecoder(r.Body).Decode(&task)
		result, err := h.TaskSvc.UpdateTaskByID(h.Ctx, &task)

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			json.NewEncoder(w).Encode(encoder.FromError(err))
			return
		}

		json.NewEncoder(w).Encode(encoder.FromResponse(result))
	}

	return fn
}

func (h TasksHTTPHandler) DeleteTaskByIdHTTPHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {

		id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)

		err := h.TaskSvc.DeleteTaskByID(h.Ctx, id)

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			json.NewEncoder(w).Encode(encoder.FromError(err))
			return
		}

		json.NewEncoder(w).Encode(encoder.FromResponse(nil))
	}

	return fn
}
