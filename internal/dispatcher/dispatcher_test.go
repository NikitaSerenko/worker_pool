package dispatcher

import (
	"context"
	"testing"
	"time"

	"worker_pool/internal/task"
)

func TestGetSnapshot(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	d := NewDispatcher(5)
	d.Start(ctx)

	d.CreateTasks([]task.Service{
		task.NewSleepTask("cba", 5*time.Second),
		task.NewSleepTask("cba", 5*time.Second),
		task.NewSleepTask("cba", 5*time.Second),
	})

	time.Sleep(2 * time.Second)

	snapshot := d.GetSnapshot()

	if len(snapshot.Running) != 1 {
		t.Errorf("invalid number of running tasks")
	}

	if len(snapshot.Waiting) != 2 {
		t.Errorf("invalid number of waiting tasks")
	}

	d.Stop()
}
