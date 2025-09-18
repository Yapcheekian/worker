package worker

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool_Basic(t *testing.T) {
	w := New(context.Background(), 3)

	var counter int32

	for i := 0; i < 10; i++ {
		w.Add(func() {
			atomic.AddInt32(&counter, 1)
		})
	}

	w.Wait()

	if counter != 10 {
		t.Errorf("expected counter 10, got %d", counter)
	}
}

func TestWorkerPool_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	w := New(ctx, 2)

	var counter int32

	for i := 0; i < 5; i++ {
		w.Add(func() {
			time.Sleep(100 * time.Millisecond)
			atomic.AddInt32(&counter, 1)
		})
	}

	cancel()

	w.Wait()

	if counter >= 5 {
		t.Errorf("expected fewer than 5 executed tasks due to cancel, got %d", counter)
	}
}
