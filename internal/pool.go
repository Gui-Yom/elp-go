package internal

import "sync"

// TaskPool a thread-safe task pool.
// Works by reslicing the inner slice under a lock.
type TaskPool struct {
	tasks []Task
	lock  sync.Mutex
}

func NewTaskPool(tasks []Task) *TaskPool {
	return &TaskPool{tasks: tasks, lock: sync.Mutex{}}
}

func (pool *TaskPool) Get() Task {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	if len(pool.tasks) > 0 {
		old := pool.tasks[0]
		pool.tasks = pool.tasks[1:]
		return old
	} else {
		return nil
	}
}
