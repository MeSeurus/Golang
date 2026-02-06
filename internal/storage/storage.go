package storage

import "golang/internal/models"

// Интерфейс хранилища
type Storage interface {
	List() []models.Task
	Create(models.Task) (models.Task, error)
	Get(id int) (*models.Task, bool)
	Update(id int, _ models.Task) (*models.Task, error)
	Delete(id int) error
}
