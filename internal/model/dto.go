package model

import (
	"errors"
	"regexp"
	"strings"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	if r.Email == "" {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email format")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

type CreatePostRequest struct {
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	PublishAt *string `json:"publish_at,omitempty"` // RFC3339 or null
}

func (r *CreatePostRequest) Validate() error {
	r.Title = strings.TrimSpace(r.Title)
	if r.Title == "" {
		return errors.New("title is required")
	}
	if len(r.Content) == 0 {
		return errors.New("content is required")
	}
	return nil
}

type UpdatePostRequest struct {
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	PublishAt *string `json:"publish_at,omitempty"`
	Status    *string `json:"status,omitempty"`
}

func (r *UpdatePostRequest) Validate() error {
	r.Title = strings.TrimSpace(r.Title)
	if r.Title == "" {
		return errors.New("title is required")
	}
	if len(r.Content) == 0 {
		return errors.New("content is required")
	}
	if r.Status != nil {
		status := PostStatus(*r.Status)
		if status != StatusDraft && status != StatusPublished {
			return errors.New("status must be 'draft' or 'published'")
		}
	}
	return nil
}

type CreateCommentRequest struct {
	Content string `json:"content"`
}

func (r *CreateCommentRequest) Validate() error {
	r.Content = strings.TrimSpace(r.Content)
	if r.Content == "" {
		return errors.New("content is required")
	}
	return nil
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
