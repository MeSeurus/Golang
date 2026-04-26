package service

import (
	"testing"
	"time"

	"golang/internal/model"
)

// Mock repository
type mockPostRepo struct {
	posts    map[int]*model.Post
	sequence int
}

func (m *mockPostRepo) Create(post *model.Post) error {
	m.sequence++
	post.ID = m.sequence
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepo) GetByID(id int) (*model.Post, error) {
	p, ok := m.posts[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

// mockPostRepo уже определён, но для PublishScheduledPosts нужно расширить его методы
func (m *mockPostRepo) GetPostsToPublish(now time.Time) ([]*model.Post, error) {
	var res []*model.Post
	for _, p := range m.posts {
		if p.Status == model.StatusDraft && p.PublishAt != nil && !p.PublishAt.After(now) {
			res = append(res, p)
		}
	}
	return res, nil
}

func (m *mockPostRepo) PublishPost(id int) error {
	if p, ok := m.posts[id]; ok {
		p.Status = model.StatusPublished
		return nil
	}
	return nil
}

func (m *mockPostRepo) Update(post *model.Post) error {
	if _, ok := m.posts[post.ID]; ok {
		m.posts[post.ID] = post
		return nil
	}
	return nil
}

func (m *mockPostRepo) Delete(id int) error {
	delete(m.posts, id)
	return nil
}

func (m *mockPostRepo) GetAll(limit, offset int) ([]*model.Post, error) {
	return nil, nil
}

func TestCreatePostImmediate(t *testing.T) {
	repo := &mockPostRepo{posts: make(map[int]*model.Post)}
	svc := NewPostService(repo)

	req := &model.CreatePostRequest{
		Title:   "Test",
		Content: "Body",
	}
	post, err := svc.Create(1, req)
	if err != nil {
		t.Fatal(err)
	}
	if post.Status != model.StatusPublished {
		t.Errorf("expected status published, got %s", post.Status)
	}
}

func TestCreateScheduledPost(t *testing.T) {
	repo := &mockPostRepo{posts: make(map[int]*model.Post)}
	svc := NewPostService(repo)

	future := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	req := &model.CreatePostRequest{
		Title:     "Future",
		Content:   "Will be published later",
		PublishAt: &future,
	}
	post, err := svc.Create(1, req)
	if err != nil {
		t.Fatal(err)
	}
	if post.Status != model.StatusDraft {
		t.Errorf("expected status draft for future publish, got %s", post.Status)
	}
}

func TestPublishScheduledPosts(t *testing.T) {
	repo := &mockPostRepo{posts: make(map[int]*model.Post)}
	svc := NewPostService(repo)

	// Создаём отложенный пост, который должен быть опубликован (время в прошлом)
	pastTime := time.Now().Add(-1 * time.Hour)
	post := &model.Post{
		ID:        1,
		UserID:    1,
		Title:     "Scheduled",
		Content:   "Content",
		Status:    model.StatusDraft,
		PublishAt: &pastTime,
	}
	repo.posts[1] = post

	published, err := svc.PublishScheduledPosts()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if published != 1 {
		t.Errorf("expected 1 published post, got %d", published)
	}
	if repo.posts[1].Status != model.StatusPublished {
		t.Error("post status was not updated to published")
	}
}

func TestPublishScheduledPostsNoFuture(t *testing.T) {
	repo := &mockPostRepo{posts: make(map[int]*model.Post)}
	svc := NewPostService(repo)

	future := time.Now().Add(2 * time.Hour)
	post := &model.Post{
		ID:        1,
		Status:    model.StatusDraft,
		PublishAt: &future,
	}
	repo.posts[1] = post

	published, err := svc.PublishScheduledPosts()
	if err != nil {
		t.Fatal(err)
	}
	if published != 0 {
		t.Errorf("expected 0 published, got %d", published)
	}
}
