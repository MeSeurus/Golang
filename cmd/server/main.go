package main

import (
	"log"
	"net/http"

	"golang/internal/handlers"
	"golang/internal/storage"
)

func main() {
	// TODO: подключите конкретную реализацию (in‑memory) интерфейса Storage - подключено
	store := storage.NewMemoryStorage()

	h := handlers.New(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", h.TasksCollection) // GET, POST
	mux.HandleFunc("/tasks/", h.TaskItem)       // GET, PUT, DELETE by index

	log.Println("server listening on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}
}
