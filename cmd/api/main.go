package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"golang/internal/config"
	"golang/internal/handler"
	"golang/internal/middleware"
	"golang/internal/repository"
	"golang/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := sqlx.Connect("postgres", buildDSN(cfg))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	postService := service.NewPostService(postRepo)
	commentService := service.NewCommentService(commentRepo, postRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)

	// Router (Go 1.22+ pattern)
	mux := http.NewServeMux()

	// Публичные маршруты
	mux.HandleFunc("GET /api/health", handler.HealthCheck)
	mux.HandleFunc("POST /api/register", authHandler.Register)
	mux.HandleFunc("POST /api/login", authHandler.Login)
	mux.HandleFunc("GET /api/posts", postHandler.GetAll)
	mux.HandleFunc("GET /api/posts/{id}", postHandler.GetByID)
	mux.HandleFunc("GET /api/posts/{postId}/comments", commentHandler.GetByPostID)

	// Защищённые маршруты (оборачиваем каждый JWT‑middleware)
	protected := func(pattern string, h http.HandlerFunc) {
		mux.Handle(pattern, middleware.JWTAuth(authService)(http.HandlerFunc(h)))
	}
	protected("POST /api/posts", postHandler.Create)
	protected("PUT /api/posts/{id}", postHandler.Update)
	protected("DELETE /api/posts/{id}", postHandler.Delete)
	protected("POST /api/posts/{postId}/comments", commentHandler.Create)

	// Глобальные middleware (логирование, восстановление)
	handler := middleware.Logger(
		middleware.Recovery(mux),
	)

	// Планировщик
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	scheduler := service.NewScheduler(postService, cfg)
	wg.Add(1)
	go scheduler.Start(ctx, &wg)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down server...")
		cancel()
		wg.Wait()

		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownRelease()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("HTTP shutdown error: %v", err)
		}
		log.Println("Server stopped")
	}()

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}

func buildDSN(cfg *config.Config) string {
	return "host=" + cfg.DBHost +
		" port=" + cfg.DBPort +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" sslmode=" + cfg.DBSSLMode
}

func runMigrations(db *sqlx.DB) error {
	migrationSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'draft',
		publish_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
	CREATE INDEX IF NOT EXISTS idx_posts_publish_at ON posts(publish_at) WHERE status = 'draft';

	CREATE TABLE IF NOT EXISTS comments (
		id SERIAL PRIMARY KEY,
		post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
	`
	_, err := db.Exec(migrationSQL)
	if err != nil {
		return err
	}
	log.Println("Migrations applied successfully")
	return nil
}
