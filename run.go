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
	// Place your code here.

	var wg sync.WaitGroup
	done := make(chan interface{})
	c := genTasks(done, tasks)
	var count int32

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			var err error
			defer wg.Done()
			for v := range c {
				err = v()
				if int(count) > m-1 {
					return
				}
				if err != nil {
					atomic.AddInt32(&count, 1)
				}
			}
		}()
	}
	wg.Wait()
	if int(count) > m-1 {
		close(done)
		return ErrErrorsLimitExceeded
	}
	return nil
}

func genTasks(done chan interface{}, tasks []Task) <-chan Task {
	out := make(chan Task)

	go func() {
		defer close(out)
		for _, task := range tasks {
			select {
			case <-done:
				return
			case out <- task:
			}
		}
	}()

	return out
}
