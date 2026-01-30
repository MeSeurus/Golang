package database

import (
	"database/sql"

	"github.com/MeSeurus/Golang/internal/models"
)

// SQLiteRepository хранит логику работы с SQLite
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository возвращает новый экземпляр SQLiteRepository
func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

// Save сохраняет заказ в SQLite
func (repo *SQLiteRepository) Save(order *models.Order) error {
	stmt, err := repo.db.Prepare(`INSERT INTO orders (customer, products, total, status) VALUES (?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(order.Customer, order.Products, order.Total, order.Status)
	return err
}
