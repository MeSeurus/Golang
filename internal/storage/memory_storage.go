package storage

import (
	"errors"
	"sync"
	"time"

	"golang/internal/models"
)

// Структура in-memory хранилища
type MemoryStorage struct {
	tasks map[int]*models.Task
	mu    sync.RWMutex
	idGen int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks: make(map[int]*models.Task),
	}
}

// Реализация методов интерфейса с учетом специфики
// конкурентного взаимодействия и необходимости синхронизации
// участков кода
func (s *MemoryStorage) List() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []models.Task
	for _, task := range s.tasks {
		result = append(result, *task)
	}
	return result
}

func (s *MemoryStorage) Create(task models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task.ID = s.idGen
	task.CreatedAt = time.Now().UTC().String()
	s.tasks[s.idGen] = &task
	s.idGen++
	return task, nil
}

func (s *MemoryStorage) Get(id int) (*models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[id]
	return task, exists
}

func (s *MemoryStorage) Update(id int, updatedTask models.Task) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existingTask, ok := s.tasks[id]; !ok {
		return nil, ErrNotFound
	} else {
		updatedTask.ID = id
		updatedTask.CreatedAt = existingTask.CreatedAt
		s.tasks[id] = &updatedTask
		return &updatedTask, nil
	}
}

func (s *MemoryStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, found := s.tasks[id]; found {
		delete(s.tasks, id)
		return nil
	}
	return ErrNotFound
}

var ErrNotFound = errors.New("Not found")
