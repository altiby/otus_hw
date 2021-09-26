package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(ctx context.Context, wg *sync.WaitGroup, tasks <-chan Task, errCounter *int32) {
	defer wg.Done()
	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				return
			}
			if err := task(); err != nil {
				atomic.AddInt32(errCounter, 1)
			}
		case <-ctx.Done():
			return
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	tasksChain := make(chan Task)
	defer close(tasksChain)
	wg := &sync.WaitGroup{}

	var errCounter int32 = 0
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(ctx, wg, tasksChain, &errCounter)
	}

	waitForEnd := func() { cancelFunc(); wg.Wait() }

	for _, task := range tasks {
		if m != 0 && int(atomic.LoadInt32(&errCounter)) >= m {
			waitForEnd()
			return ErrErrorsLimitExceeded
		}
		tasksChain <- task
	}
	waitForEnd()
	return nil
}
