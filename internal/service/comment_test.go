package service

import (
	"testing"

	"golang/internal/model"
)

type mockCommentRepo struct {
	comments []*model.Comment
}

func (m *mockCommentRepo) Create(comment *model.Comment) error {
	comment.ID = len(m.comments) + 1
	m.comments = append(m.comments, comment)
	return nil
}

func (m *mockCommentRepo) GetByPostID(postID int) ([]*model.Comment, error) {
	var res []*model.Comment
	for _, c := range m.comments {
		if c.PostID == postID {
			res = append(res, c)
		}
	}
	return res, nil
}

func TestCommentCreate_Success(t *testing.T) {
	postRepo := &mockPostRepo{posts: make(map[int]*model.Post)}
	postRepo.posts[10] = &model.Post{ID: 10} // пост существует
	commentRepo := &mockCommentRepo{}
	svc := NewCommentService(commentRepo, postRepo)

	req := &model.CreateCommentRequest{Content: "Nice post"}
	comment, err := svc.Create(1, 10, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID != 1 {
		t.Errorf("expected ID 1, got %d", comment.ID)
	}
	if comment.Content != req.Content {
		t.Error("content mismatch")
	}
}

func TestCommentCreate_PostNotFound(t *testing.T) {
	postRepo := &mockPostRepo{posts: make(map[int]*model.Post)}
	commentRepo := &mockCommentRepo{}
	svc := NewCommentService(commentRepo, postRepo)

	_, err := svc.Create(1, 999, &model.CreateCommentRequest{Content: "test"})
	if err != ErrPostNotFoundForComment {
		t.Errorf("expected ErrPostNotFoundForComment, got %v", err)
	}
}

func TestGetCommentsByPostID(t *testing.T) {
	commentRepo := &mockCommentRepo{}
	commentRepo.comments = []*model.Comment{
		{ID: 1, PostID: 1, Content: "First"},
		{ID: 2, PostID: 1, Content: "Second"},
		{ID: 3, PostID: 2, Content: "Other"},
	}
	svc := NewCommentService(commentRepo, nil) // postRepo не понадобится

	comments, err := svc.GetByPostID(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}
}
