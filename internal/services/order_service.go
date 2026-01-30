package services

import (
	"fmt"

	"github.com/MeSeurus/Golang/internal/models"
)

// OrderService управляет обработкой заказов
type OrderService struct {
	repository models.RepositoryWriter
	notifier   models.Notifier
}

// NewOrderService возвращает новый экземпляр OrderService
func NewOrderService(repo models.RepositoryWriter, notifier models.Notifier) *OrderService {
	return &OrderService{
		repository: repo,
		notifier:   notifier,
	}
}

// CreateOrder создает новый заказ и отправляет уведомление
func (os *OrderService) CreateOrder(customer string, products []string, total float64) (*models.Order, error) {
	order := &models.Order{
		Customer: customer,
		Products: fmt.Sprintf("%v", products),
		Total:    total,
		Status:   "pending",
	}

	err := os.repository.Save(order)
	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf("Заказ клиента %s успешно сохранён.", customer)
	err = os.notifier.Send(message)
	if err != nil {
		return nil, err
	}

	return order, nil
}
