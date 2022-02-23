package dispatcher

import (
	"context"
	"sync"

	"worker_pool/internal/task"
)

type Dispatcher interface {
	Start(ctx context.Context)
	Stop()
}

type Scheduler interface {
	CreateTasks(tasks []task.Service)
	GetSnapshot() Snapshot
}

type dispatcher struct {
	numberOfWorkers int
	workerQueue     chan workTicket
	waitingQueue    chan []task.Service

	counter *counter

	cancelFunc context.CancelFunc
	wg         *sync.WaitGroup
}

var _ Dispatcher = &dispatcher{}
var _ Scheduler = &dispatcher{}

func NewDispatcher(numberOfWorkers int) *dispatcher {
	return &dispatcher{
		numberOfWorkers: numberOfWorkers,

		workerQueue:  make(chan workTicket, numberOfWorkers),
		waitingQueue: make(chan []task.Service),

		counter: newCounter(),
		wg:      &sync.WaitGroup{},
	}
}

func (d *dispatcher) Start(ctx context.Context) {
	ctx, cancelFunc := context.WithCancel(ctx)
	d.cancelFunc = cancelFunc

	d.wg.Add(d.numberOfWorkers)
	for i := 0; i < d.numberOfWorkers; i++ {
		newWorker(i+1, d.workerQueue, d.wg, d.counter).Start(ctx)
	}

	go d.runDispatcherWorker(ctx)
}

func (d *dispatcher) Stop() {
	d.cancelFunc()
	d.wg.Wait()
}

func (d *dispatcher) CreateTasks(tasks []task.Service) {
	d.waitingQueue <- tasks
}

func (d *dispatcher) GetSnapshot() Snapshot {
	return d.counter.getSnapshot()
}

func (d *dispatcher) runDispatcherWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case works := <-d.waitingQueue:
			d.counter.addWaitingTask(works)

			for _, work := range works {
				go func(work task.Service) {
					for {
						workerTicket := <-d.workerQueue

						if d.counter.moveToRunningTask(work) {
							workerTicket.work <- work
							return
						}

						// task with this name is processing already, move workerTicket back to the worker pool
						d.workerQueue <- workerTicket
					}
				}(work)
			}
		}
	}
}
