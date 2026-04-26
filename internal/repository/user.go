package repository

import (
	"database/sql"
	"golang/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetByID(id int) (*model.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *model.User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, created_at`
	return r.db.QueryRow(query, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
}

func (r *userRepo) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := r.db.Get(user, "SELECT * FROM users WHERE email=$1", email)
	if err == sql.ErrNoRows {
		return nil, nil // not found, no error
	}
	return user, err
}

func (r *userRepo) GetByID(id int) (*model.User, error) {
	user := &model.User{}
	err := r.db.Get(user, "SELECT * FROM users WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}
