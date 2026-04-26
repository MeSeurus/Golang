package repository

import (
	"golang/internal/model"

	"github.com/jmoiron/sqlx"
)

type CommentRepository interface {
	Create(comment *model.Comment) error
	GetByPostID(postID int) ([]*model.Comment, error)
}

type commentRepo struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) CommentRepository {
	return &commentRepo{db: db}
}

func (r *commentRepo) Create(comment *model.Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(query, comment.PostID, comment.UserID, comment.Content).Scan(&comment.ID, &comment.CreatedAt)
}

func (r *commentRepo) GetByPostID(postID int) ([]*model.Comment, error) {
	comments := []*model.Comment{}
	err := r.db.Select(&comments, `SELECT * FROM comments WHERE post_id=$1 ORDER BY created_at ASC`, postID)
	return comments, err
}
