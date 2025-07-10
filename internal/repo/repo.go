package repo

import (
	"errors"
	"sync"
)

type repository struct {
	inMemory map[int]Task
	mu       sync.RWMutex
	autoID   int
}

type Repository interface {
	CreateTask(task Task) (int, error)
	GetTaskById(id int) (*Task, error)
	GetTasks() ([]Task, error)
	UpdateTask(task Task)
	DeleteTask(id int)
}

func NewRepository(capacity ...int) Repository {
	cap := 0
	if len(capacity) > 0 {
		cap = capacity[0]
	}
	return &repository{
		inMemory: make(map[int]Task, cap),
		autoID:   1,
	}
}

func (r *repository) CreateTask(task Task) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.autoID
	r.autoID++

	r.inMemory[task.ID] = task
	return task.ID, nil
}

func (r *repository) GetTaskById(id int) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.inMemory[id]
	if !ok {
		return nil, errors.New("task not found")
	}
	return &task, nil
}

func (r *repository) GetTasks() ([]Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]Task, 0, len(r.inMemory))
	for _, task := range r.inMemory {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *repository) UpdateTask(task Task) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.inMemory[task.ID] = task
}

func (r *repository) DeleteTask(id int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.inMemory, id)
}
