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
	t.Run("returns error when errors count equals limit", func(t *testing.T) {
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

	t.Run("treats non-positive max errors as zero allowed errors", func(t *testing.T) {
		testCases := []struct {
			name     string
			maxErr   int
			taskErr  error
			expected error
		}{
			{
				name:     "zero max errors returns error after first failed task",
				maxErr:   0,
				taskErr:  errors.New("error from task"),
				expected: ErrErrorsLimitExceeded,
			},
			{
				name:     "negative max errors returns error after first failed task",
				maxErr:   -1,
				taskErr:  errors.New("error from task"),
				expected: ErrErrorsLimitExceeded,
			},
			{
				name:   "zero max errors returns nil when task succeeds",
				maxErr: 0,
			},
			{
				name:   "negative max errors returns nil when task succeeds",
				maxErr: -1,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tasks := []Task{func() error {
					return tc.taskErr
				}}

				err := Run(tasks, 10, tc.maxErr)
				require.Equal(t, tc.expected, err)
			})
		}
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

	t.Run("returns nil when there is nothing to run", func(t *testing.T) {
		testCases := []struct {
			name  string
			tasks []Task
			n     int
			m     int
		}{
			{
				name: "zero workers",
				tasks: []Task{func() error {
					return nil
				}},
				n: 0,
				m: 1,
			},
			{
				name:  "empty tasks",
				tasks: []Task{},
				n:     10,
				m:     1,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := Run(tc.tasks, tc.n, tc.m)
				require.NoError(t, err)
			})
		}
	})
}
