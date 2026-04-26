package repository

import (
	"database/sql"
	"time"

	"golang/internal/model"

	"github.com/jmoiron/sqlx"
)

type PostRepository interface {
	Create(post *model.Post) error
	GetByID(id int) (*model.Post, error)
	Update(post *model.Post) error
	Delete(id int) error
	GetAll(limit, offset int) ([]*model.Post, error)
	GetPostsToPublish(now time.Time) ([]*model.Post, error)
	PublishPost(id int) error
}

type postRepo struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) PostRepository {
	return &postRepo{db: db}
}

func (r *postRepo) Create(post *model.Post) error {
	query := `
		INSERT INTO posts (user_id, title, content, status, publish_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query,
		post.UserID, post.Title, post.Content, post.Status, post.PublishAt,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (r *postRepo) GetByID(id int) (*model.Post, error) {
	post := &model.Post{}
	err := r.db.Get(post, "SELECT * FROM posts WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return post, err
}

func (r *postRepo) Update(post *model.Post) error {
	query := `
		UPDATE posts SET title=$1, content=$2, status=$3, publish_at=$4, updated_at=NOW()
		WHERE id=$5
		RETURNING updated_at`
	return r.db.QueryRow(query,
		post.Title, post.Content, post.Status, post.PublishAt, post.ID,
	).Scan(&post.UpdatedAt)
}

func (r *postRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE id=$1", id)
	return err
}

func (r *postRepo) GetAll(limit, offset int) ([]*model.Post, error) {
	posts := []*model.Post{}
	query := `SELECT * FROM posts WHERE status='published' ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.Select(&posts, query, limit, offset)
	return posts, err
}

func (r *postRepo) GetPostsToPublish(now time.Time) ([]*model.Post, error) {
	posts := []*model.Post{}
	err := r.db.Select(&posts,
		`SELECT * FROM posts WHERE status='draft' AND publish_at IS NOT NULL AND publish_at <= $1`,
		now)
	return posts, err
}

func (r *postRepo) PublishPost(id int) error {
	_, err := r.db.Exec(`UPDATE posts SET status='published', updated_at=NOW() WHERE id=$1`, id)
	return err
}
