package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrorsCounter struct {
	sync.RWMutex
	count     int
	maxErrors int
}

func (e *ErrorsCounter) inc() {
	e.Lock()
	defer e.Unlock()

	e.count++
}

func (e *ErrorsCounter) isLimitExceeded() bool {
	e.RLock()
	defer e.RUnlock()

	return e.count >= e.maxErrors
}

var wg sync.WaitGroup

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	done := make(chan struct{}, 1)
	tasksChan := fillTasksChan(tasks)
	errorsCounter := &ErrorsCounter{maxErrors: m}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go worker(done, tasksChan, errorsCounter)
	}

	wg.Wait()

	select {
	case <-done:
		if m > 0 {
			return ErrErrorsLimitExceeded
		}

		return nil

	default:
		return nil
	}
}

func worker(done chan struct{}, tasksChan <-chan Task, e *ErrorsCounter) {
	defer wg.Done()

	for task := range tasksChan {
		if e.isLimitExceeded() {
			done <- struct{}{}
			close(done)
		}

		err := task()
		if err != nil {
			e.inc()
		}

		select {
		case <-done:
			return
		default:
			continue
		}
	}
}

func fillTasksChan(tasks []Task) <-chan Task {
	tasksChan := make(chan Task, len(tasks))
	defer close(tasksChan)

	for _, task := range tasks {
		tasksChan <- task
	}

	return tasksChan
}
