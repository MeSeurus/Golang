package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/MeSeurus/Golang/internal/database"
	"github.com/MeSeurus/Golang/internal/notifications"
	"github.com/MeSeurus/Golang/internal/services"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "orders.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS orders (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        customer TEXT NOT NULL,
        products TEXT NOT NULL,
        total REAL NOT NULL,
        status TEXT NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}

	sqliteRepo := database.NewSQLiteRepository(db)
	emailNotif := new(notifications.EmailSender)
	smsNotif := new(notifications.SMSSender)

	serviceWithEmail := services.NewOrderService(sqliteRepo, emailNotif)
	serviceWithSMS := services.NewOrderService(sqliteRepo, smsNotif)

	order1, err := serviceWithEmail.CreateOrder("Иван", []string{"яблоки", "бананы"}, 10.5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Создан заказ №%d\n", order1.ID)

	order2, err := serviceWithSMS.CreateOrder("Елена", []string{"груши", "виноград"}, 8.0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Создан заказ №%d\n", order2.ID)
}
