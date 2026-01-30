package notifications

import (
	"fmt"
)

// SMSSender реализует интерфейс Notifier для отправки SMS
type SMSSender struct{}

// Send отправляет уведомление через SMS
func (ss *SMSSender) Send(message string) error {
	fmt.Println("Отправлена SMS:", message)
	return nil
}
