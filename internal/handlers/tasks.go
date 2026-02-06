package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang/internal/models"
	"golang/internal/storage"
)

type Handler struct {
	Store storage.Storage
}

func New(s storage.Storage) *Handler {
	return &Handler{Store: s}
}

// /tasks (GET, POST)
func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	// TODO: реализуйте разбор метода, JSON, коды статусов, валидацию
	log.Printf("%s %s", r.Method, r.URL)
	switch r.Method {
	case http.MethodGet:
		h.handleListTasks(w, r)
	case http.MethodPost:
		h.handleCreateTask(w, r)
	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// /tasks/{id} (GET, PUT, DELETE)
func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	// TODO: извлечение id, маршрутизация по методу, ошибки - добавлено
	// Ошибки взяты из библиотеки http
	log.Printf("%s %s", r.Method, r.URL)
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) <= 2 || parts[len(parts)-1] == "" {
		http.NotFound(w, r)
		return
	}
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Invalid id format '%v'"}`, idStr), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetTaskByID(w, r, id)
	case http.MethodPut:
		h.handleUpdateTask(w, r, id)
	case http.MethodDelete:
		h.handleDeleteTask(w, r, id)
	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// Обработка списка задач
func (h *Handler) handleListTasks(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tasks := h.Store.List()
	json.NewEncoder(w).Encode(tasks)
}

// Создание задачи
func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, `{"error":"Failed to decode request body"}`, http.StatusBadRequest)
		return
	}
	if newTask.Title == "" {
		http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
		return
	}
	createdTask, err := h.Store.Create(newTask)
	if err != nil {
		log.Printf("Error creating task: %v\n", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

// Получение задачи по ID
func (h *Handler) handleGetTaskByID(w http.ResponseWriter, _ *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")
	task, found := h.Store.Get(id)
	if !found {
		http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(task)
}

// Обновление задачи
func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request, id int) {
	var updatedTask models.Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, `{"error":"Failed to decode request body"}`, http.StatusBadRequest)
		return
	}
	if updatedTask.Title == "" {
		http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
		return
	}
	task, err := h.Store.Update(id, updatedTask)
	if err != nil {
		http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Удаление задачи
func (h *Handler) handleDeleteTask(w http.ResponseWriter, _ *http.Request, id int) {
	err := h.Store.Delete(id)
	if err != nil {
		http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
