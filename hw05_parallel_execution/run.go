package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	tasksCh := make(chan Task)
	var currentErrors int64

	defer func() {
		close(tasksCh)
		wg.Wait()
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range tasksCh {
				if task() != nil {
					atomic.AddInt64(&currentErrors, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		tasksCh <- task

		if atomic.LoadInt64(&currentErrors) >= int64(m) {
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}
