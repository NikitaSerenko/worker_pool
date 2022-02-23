package dispatcher

import (
	"sync"

	"worker_pool/internal/task"
)

type counter struct {
	waitingTaskMap map[string]int  // taskName: count
	runningTaskMap map[string]bool // taskName: bool

	mtx sync.Mutex
}

type Snapshot struct {
	Running []string
	Waiting []string
}

func newCounter() *counter {
	return &counter{
		waitingTaskMap: make(map[string]int),
		runningTaskMap: make(map[string]bool),
	}
}

func (c *counter) addWaitingTask(works []task.Service) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	for _, work := range works {
		c.waitingTaskMap[work.Name()]++
	}
}

func (c *counter) finishWorkingTask(work task.Service) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.runningTaskMap[work.Name()] = false
}

func (c *counter) moveToRunningTask(work task.Service) (success bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.runningTaskMap[work.Name()] {
		return false
	}

	c.waitingTaskMap[work.Name()]--
	c.runningTaskMap[work.Name()] = true

	return true
}

func (c *counter) getSnapshot() Snapshot {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	s := Snapshot{}
	for taskName, count := range c.waitingTaskMap {
		for i := 0; i < count; i++ {
			s.Waiting = append(s.Waiting, taskName)
		}
	}

	for taskName, ok := range c.runningTaskMap {
		if ok {
			s.Running = append(s.Running, taskName)
		}
	}

	return s
}
