package worker

import (
	"context"
	"errors"
	"sync"
)

const defaultTaskBufferCount = 100

var (
	ErrWorkerPoolClosed = errors.New("worker pool is closed")
)

type Worker struct {
	tasks  chan func()
	wg     sync.WaitGroup
	ctx    context.Context
	closed bool
}

func New(ctx context.Context, workerCount int) *Worker {
	w := &Worker{
		tasks: make(chan func(), defaultTaskBufferCount),
		ctx:   ctx,
	}

	for i := 0; i < workerCount; i++ {
		w.wg.Add(1)
		go w.execute()
	}

	return w
}

func (w *Worker) Add(task func()) error {
	if w.closed {
		return ErrWorkerPoolClosed
	}

	select {
	case <-w.ctx.Done():
		return w.ctx.Err()
	case w.tasks <- task:
	}

	return nil
}

func (w *Worker) Wait() {
	w.stop()
	w.wg.Wait()
}

func (w *Worker) execute() {
	defer w.wg.Done()
	for {
		select {
		case <-w.ctx.Done():
			return
		case task, ok := <-w.tasks:
			if !ok {
				return
			}
			task()
		}
	}
}

func (w *Worker) stop() {
	if w.closed {
		return
	}

	close(w.tasks)
	w.closed = true
}
