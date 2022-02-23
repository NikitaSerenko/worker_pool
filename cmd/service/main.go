package main

import (
	"context"
	"fmt"
	"net/http"

	"worker_pool/api"
	"worker_pool/internal/dispatcher"
)

func main() {
	ctx := context.Background()

	dispatcherEntity := dispatcher.NewDispatcher(10)
	dispatcherEntity.Start(ctx)

	schedulerAPI := api.NewSchedulerController(dispatcherEntity)
	http.HandleFunc("/tasks", schedulerAPI.CreateTasks)
	http.HandleFunc("/get_tasks", schedulerAPI.GetSnapshot)

	if err := http.ListenAndServe("127.0.0.1:8000", nil); err != nil {
		fmt.Println(err.Error())
	}
}
