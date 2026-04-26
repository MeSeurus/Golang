package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang/internal/model"
	"golang/internal/service"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	var req model.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	post, err := h.postService.Create(userID, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}
	respondJSON(w, http.StatusCreated, post)
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}
	post, err := h.postService.GetByID(id)
	if err != nil || post == nil {
		respondError(w, http.StatusNotFound, "Post not found")
		return
	}
	respondJSON(w, http.StatusOK, post)
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}
	var req model.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	post, err := h.postService.Update(userID, postID, &req)
	if err != nil {
		if err == service.ErrPostNotFound {
			respondError(w, http.StatusNotFound, "Post not found")
		} else if err == service.ErrForbidden {
			respondError(w, http.StatusForbidden, "You can only edit your own posts")
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to update post")
		}
		return
	}
	respondJSON(w, http.StatusOK, post)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}
	err = h.postService.Delete(userID, postID)
	if err != nil {
		if err == service.ErrPostNotFound {
			respondError(w, http.StatusNotFound, "Post not found")
		} else if err == service.ErrForbidden {
			respondError(w, http.StatusForbidden, "You can only delete your own posts")
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to delete post")
		}
		return
	}
	respondJSON(w, http.StatusOK, model.MessageResponse{Message: "Post deleted"})
}

func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	posts, err := h.postService.GetAll(limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch posts")
		return
	}
	respondJSON(w, http.StatusOK, posts)
}
