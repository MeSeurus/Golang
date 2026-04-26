package service

import (
	"errors"
	"log"
	"time"

	"golang/internal/model"
	"golang/internal/repository"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrForbidden    = errors.New("forbidden")
)

type PostService interface {
	Create(userID int, req *model.CreatePostRequest) (*model.Post, error)
	GetByID(id int) (*model.Post, error)
	Update(userID, postID int, req *model.UpdatePostRequest) (*model.Post, error)
	Delete(userID, postID int) error
	GetAll(limit, offset int) ([]*model.Post, error)
	PublishScheduledPosts() (int, error)
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{postRepo: postRepo}
}

func (s *postService) Create(userID int, req *model.CreatePostRequest) (*model.Post, error) {
	post := &model.Post{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
		Status:  model.StatusPublished, // default
	}

	if req.PublishAt != nil {
		t, err := time.Parse(time.RFC3339, *req.PublishAt)
		if err != nil {
			return nil, errors.New("invalid publish_at format, use RFC3339")
		}
		if t.After(time.Now()) {
			post.Status = model.StatusDraft
			post.PublishAt = &t
		} else {
			// if time is in the past, still set publish_at but keep it published
			post.PublishAt = &t
		}
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *postService) GetByID(id int) (*model.Post, error) {
	return s.postRepo.GetByID(id)
}

func (s *postService) Update(userID, postID int, req *model.UpdatePostRequest) (*model.Post, error) {
	post, err := s.postRepo.GetByID(postID)
	if err != nil || post == nil {
		return nil, ErrPostNotFound
	}
	if post.UserID != userID {
		return nil, ErrForbidden
	}

	post.Title = req.Title
	post.Content = req.Content
	if req.Status != nil {
		post.Status = model.PostStatus(*req.Status)
	}
	if req.PublishAt != nil {
		t, err := time.Parse(time.RFC3339, *req.PublishAt)
		if err != nil {
			return nil, errors.New("invalid publish_at format")
		}
		post.PublishAt = &t
	}

	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *postService) Delete(userID, postID int) error {
	post, err := s.postRepo.GetByID(postID)
	if err != nil || post == nil {
		return ErrPostNotFound
	}
	if post.UserID != userID {
		return ErrForbidden
	}
	return s.postRepo.Delete(postID)
}

func (s *postService) GetAll(limit, offset int) ([]*model.Post, error) {
	return s.postRepo.GetAll(limit, offset)
}

func (s *postService) PublishScheduledPosts() (int, error) {
	now := time.Now().UTC()
	posts, err := s.postRepo.GetPostsToPublish(now)
	if err != nil {
		return 0, err
	}

	published := 0
	for _, p := range posts {
		if err := s.postRepo.PublishPost(p.ID); err != nil {
			log.Printf("Failed to publish post ID %d: %v", p.ID, err)
		} else {
			published++
			log.Printf("Post ID %d published at %s", p.ID, now.Format(time.RFC3339))
		}
	}
	return published, nil
}
