package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang/internal/model"
	"golang/internal/service"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	postID, err := strconv.Atoi(r.PathValue("postId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}
	var req model.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	comment, err := h.commentService.Create(userID, postID, &req)
	if err != nil {
		if err == service.ErrPostNotFoundForComment {
			respondError(w, http.StatusNotFound, "Post not found")
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to create comment")
		}
		return
	}
	respondJSON(w, http.StatusCreated, comment)
}

func (h *CommentHandler) GetByPostID(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("postId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}
	comments, err := h.commentService.GetByPostID(postID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get comments")
		return
	}
	respondJSON(w, http.StatusOK, comments)
}
