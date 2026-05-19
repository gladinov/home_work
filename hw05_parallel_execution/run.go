package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	workersCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workersNumer := min(n, len(tasks))
	if m <= 0 {
		m = 0
	}
	errCh := make(chan error, m)
	var wg sync.WaitGroup

	jobCh := produceJob(workersCtx, tasks, workersNumer)

	for i := 0; i < workersNumer; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			doTask(workersCtx, jobCh, errCh)
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var errCount int

	for range errCh {
		errCount++
		if errCount > m {
			cancel()
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}

func produceJob(ctx context.Context, tasks []Task, workersNumber int) <-chan Task {
	jobCh := make(chan Task, workersNumber*2)

	go func() {
		defer close(jobCh)
		for _, task := range tasks {
			select {
			case <-ctx.Done():
				return
			case jobCh <- task:
			}
		}
	}()

	return jobCh
}

func doTask(ctx context.Context, jobCh <-chan Task, errCh chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-jobCh:
			if !ok {
				return
			}
			err := task()
			if err != nil {
				select {
				case errCh <- err:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}
