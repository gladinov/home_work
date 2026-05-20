package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23

		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t,
			atomic.LoadInt32(&runTasksCount),
			int32(workersCount+maxErrorsCount),
			"extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)

		require.NoError(t, err)
		require.Equal(t, int32(tasksCount), atomic.LoadInt32(&runTasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestRunExtra(t *testing.T) {
	t.Run("maxErr equal to task with errors", func(t *testing.T) {
		tasksCount := 2
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 2
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, err)
	})

	t.Run("zero max errors with first task with expected err", func(t *testing.T) {
		tasks := []Task{func() error {
			err := errors.New("error from task")
			return err
		}}

		workersCount := 10
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, err)
	})

	t.Run("less zero max errors with first task with expected err", func(t *testing.T) {
		tasks := []Task{func() error {
			err := errors.New("error from task")
			return err
		}}

		workersCount := 10
		maxErrorsCount := -1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, err)
	})

	t.Run("zero max errors with first task without err", func(t *testing.T) {
		tasks := []Task{func() error {
			return nil
		}}

		workersCount := 10
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
	})

	t.Run("less zero max errors with first task without err", func(t *testing.T) {
		tasks := []Task{func() error {
			return nil
		}}

		workersCount := 10
		maxErrorsCount := -1
		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
	})

	t.Run("runs all tasks when tasks count is less than workers count", func(t *testing.T) {
		var runTasksCount int32
		tasks := []Task{
			func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			},
			func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			},
		}

		err := Run(tasks, 10, 1)

		require.NoError(t, err)
		require.Equal(t, int32(len(tasks)), atomic.LoadInt32(&runTasksCount))
	})

	t.Run("returns nil when workers count is zero", func(t *testing.T) {
		tasks := []Task{}
		err := Run(tasks, 0, 1)
		require.NoError(t, err)
	})

	t.Run("returns nil when tasks are empty", func(t *testing.T) {
		tasks := []Task{}
		err := Run(tasks, 10, 1)
		require.NoError(t, err)
	})
}
