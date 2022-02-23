package dispatcher

import (
	"context"
	"fmt"
	"sync"

	"worker_pool/internal/task"
)

type worker struct {
	ID      int
	work    chan task.Service
	counter *counter

	workerQueue chan<- workTicket
	wg          *sync.WaitGroup
}

type workTicket struct {
	id   int
	work chan<- task.Service
}

func newWorker(id int, workerQueue chan<- workTicket, wg *sync.WaitGroup, counter *counter) *worker {
	return &worker{
		ID:      id,
		work:    make(chan task.Service),
		counter: counter,

		workerQueue: workerQueue,
		wg:          wg,
	}
}

func (w *worker) Start(ctx context.Context) {
	go func() {
		for {
			w.workerQueue <- workTicket{
				id:   w.ID,
				work: w.work,
			}

			select {
			case work := <-w.work:
				err := work.Process(ctx)
				if err != nil {
					fmt.Printf("Failed to process task: name = %s", work.Name())
				}
				w.counter.finishWorkingTask(work)

			case <-ctx.Done():
				w.wg.Done()
				return
			}
		}
	}()
}
