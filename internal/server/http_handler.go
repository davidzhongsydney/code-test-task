package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/log"
	"qantas.com/task/internal/encoder"
	"qantas.com/task/internal/service"
	"qantas.com/task/model"
)

type TasksHTTPHandler struct {
	taskSvc *service.TaskService
	ctx     context.Context
	log     *log.Helper
}

func (h TasksHTTPHandler) ListTasksHTTPHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		result, err := h.taskSvc.ListTasks(h.ctx)

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
		result, err := h.taskSvc.CreateTask(h.ctx, &task)

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			json.NewEncoder(w).Encode(encoder.FromError(err))
			return
		}

		json.NewEncoder(w).Encode(encoder.FromResponse(result))
	}
	return fn
}

func (h TasksHTTPHandler) GetTaskByIdHTTPHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)

		result, err := h.taskSvc.GetTaskByID(h.ctx, id)

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
		result, err := h.taskSvc.UpdateTaskByID(h.ctx, &task)

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

		err := h.taskSvc.DeleteTaskByID(h.ctx, id)

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			json.NewEncoder(w).Encode(encoder.FromError(err))
			return
		}

		json.NewEncoder(w).Encode(encoder.FromResponse(nil))
	}

	return fn
}

func (h TasksHTTPHandler) GetTaskService() *service.TaskService {
	return h.taskSvc
}
