package worker

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"orchestrate/task"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (a *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	te := task.TaskEvent{}
	if err := d.Decode(&te); err != nil {
		slog.Error("Error unmarshalling task event", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		e := ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			Message:        "Error unmarshalling task event",
		}
		_ = json.NewEncoder(w).Encode(e)
		return
	}

	a.Worker.AddTask(te.Task)
	slog.Info("Added task", "task", te.Task.ID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(te.Task)
}

func (a *Api) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(a.Worker.GetTasks()); err != nil {
		slog.Error("Error marshalling task %v\n", err)
	}
}

func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		slog.Error("No taskID provided in request")
		w.WriteHeader(http.StatusBadRequest)
	}

	tID, _ := uuid.Parse(taskID)
	_, ok := a.Worker.Db[tID]
	if !ok {
		slog.Error("No task with ID found", "taskID", tID)
		w.WriteHeader(http.StatusBadRequest)
	}

	taskToStop := a.Worker.Db[tID]
	taskCopy := *taskToStop
	taskCopy.State = task.Completed
	a.Worker.AddTask(taskCopy)

	slog.Info("Added task to stop container", "task", taskToStop, "container id", taskToStop.ContainerID)
	w.WriteHeader(http.StatusNoContent)
}

func (a *Api) initRouter() {
	a.Router = chi.NewRouter()
	a.Router.Route("/tasks", func(r chi.Router) {
		r.Post("/", a.StartTaskHandler) // on /tasks POST
		r.Get("/", a.GetTaskHandler)    // on /tasks GET
		r.Route("/{taskID}", func(r chi.Router) {
			r.Delete("/", a.StopTaskHandler) // on /tasks/{taskID} DELETE
		})
	})
}

func (a *Api) Start() {
	a.initRouter()
	http.ListenAndServe(fmt.Sprintf("%s:%d", a.Address, a.Port), a.Router)
}
