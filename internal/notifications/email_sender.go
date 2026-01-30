package notifications

import (
	"fmt"
)

// EmailSender реализует интерфейс Notifier для отправки e-mail
type EmailSender struct{}

// Send отправляет уведомление по электронной почте
func (es *EmailSender) Send(message string) error {
	fmt.Println("Отправлен e-mail:", message)
	return nil
}
