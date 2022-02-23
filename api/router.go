package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"worker_pool/internal/dispatcher"
	"worker_pool/internal/task"
)

type SchedulerController struct {
	scheduler dispatcher.Scheduler
}

func NewSchedulerController(scheduler dispatcher.Scheduler) SchedulerController {
	return SchedulerController{scheduler: scheduler}
}

func (s *SchedulerController) CreateTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request map[string]int
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, fmt.Errorf("failed to decode body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	tasks := make([]task.Service, 0, len(request))
	for name, sec := range request {
		tasks = append(tasks, task.NewSleepTask(name, time.Second*time.Duration(sec)))
	}

	s.scheduler.CreateTasks(tasks)

	w.WriteHeader(http.StatusCreated)
}

func (s *SchedulerController) GetSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_, _ = fmt.Fprintf(w, "%+v\n", s.scheduler.GetSnapshot())
}
