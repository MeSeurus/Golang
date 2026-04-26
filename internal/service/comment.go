package service

import (
	"errors"
	"golang/internal/model"
	"golang/internal/repository"
)

var (
	ErrPostNotFoundForComment = errors.New("post not found")
)

type CommentService interface {
	Create(userID, postID int, req *model.CreateCommentRequest) (*model.Comment, error)
	GetByPostID(postID int) ([]*model.Comment, error)
}

type commentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
}

func NewCommentService(commentRepo repository.CommentRepository, postRepo repository.PostRepository) CommentService {
	return &commentService{commentRepo: commentRepo, postRepo: postRepo}
}

func (s *commentService) Create(userID, postID int, req *model.CreateCommentRequest) (*model.Comment, error) {
	post, err := s.postRepo.GetByID(postID)
	if err != nil || post == nil {
		return nil, ErrPostNotFoundForComment
	}

	comment := &model.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: req.Content,
	}
	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *commentService) GetByPostID(postID int) ([]*model.Comment, error) {
	return s.commentRepo.GetByPostID(postID)
}
