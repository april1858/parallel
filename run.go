package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	//"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Place your code here.

	var wg sync.WaitGroup
	done := make(chan interface{})
	forMain := make(chan error)
	c := genTasks(done, tasks)
	var count int

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case v := <-c:
					forMain <- v()
				}
			}
		}()
	}

	for e := range forMain {
		if e != nil {
			count++
			fmt.Println(count)
		}
	}

	wg.Wait()

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
