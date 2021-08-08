package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type TaskCh chan Task

func runTask(taskCh TaskCh, wg *sync.WaitGroup, maxErrCount *int64) {
	for task := range taskCh {
		if atomic.LoadInt64(maxErrCount) > 0 {
			err := task()
			if err != nil {
				atomic.AddInt64(maxErrCount, -1)
			}
		}
		wg.Done()
	}
}

func addTasksToTaskCh(taskCh TaskCh, tasks []Task) {
	for _, task := range tasks {
		taskCh <- task
	}
	close(taskCh)
}

func Run(tasks []Task, n, m int) error {
	taskCh := make(TaskCh, n)

	var wg sync.WaitGroup
	wg.Add(len(tasks))

	var maxErrCount int64
	atomic.StoreInt64(&maxErrCount, int64(m))

	for i := 0; i < n; i++ {
		go runTask(taskCh, &wg, &maxErrCount)
	}

	go addTasksToTaskCh(taskCh, tasks)

	wg.Wait()
	if maxErrCount < 1 {
		return ErrErrorsLimitExceeded
	}
	return nil
}
