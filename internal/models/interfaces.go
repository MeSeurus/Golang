package models

// RepositoryWriter определяет контракт для сохранения заказов
type RepositoryWriter interface {
	Save(order *Order) error
}

// Notifier определяет контракт для отправки уведомлений
type Notifier interface {
	Send(message string) error
}
