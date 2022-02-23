package task

import (
	"context"
	"time"
)

type sleepTask struct {
	TaskName string
	Duration time.Duration
}

func NewSleepTask(name string, duration time.Duration) Service {
	return sleepTask{
		TaskName: name,
		Duration: duration,
	}
}

func (t sleepTask) Name() string {
	return t.TaskName
}

func (t sleepTask) Process(_ context.Context) error {
	time.Sleep(t.Duration)

	return nil
}
