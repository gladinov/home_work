package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if n <= 0 || len(tasks) == 0 {
		return nil
	}

	if m < 0 {
		m = 0
	}

	var (
		mu            sync.Mutex
		wg            sync.WaitGroup
		nextTaskIndex int
		errCount      int
	)

	limitExceeded := func() bool {
		return errCount > 0 && errCount >= m
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				mu.Lock()
				if nextTaskIndex >= len(tasks) || limitExceeded() {
					mu.Unlock()
					return
				}
				task := tasks[nextTaskIndex]
				nextTaskIndex++
				mu.Unlock()

				if task() != nil {
					mu.Lock()
					errCount++
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	if limitExceeded() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
