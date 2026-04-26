package service

import (
	"context"
	"log"
	"sync"
	"time"

	"golang/internal/config"
)

type Scheduler struct {
	postService PostService
	interval    time.Duration
	workers     int
}

func NewScheduler(postService PostService, cfg *config.Config) *Scheduler {
	return &Scheduler{
		postService: postService,
		interval:    cfg.SchedulerInterval,
		workers:     cfg.WorkerPoolSize,
	}
}

func (s *Scheduler) Start(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	jobs := make(chan struct{}, s.workers)

	// worker pool
	for i := 0; i < s.workers; i++ {
		go func(id int) {
			for range jobs {
				count, err := s.postService.PublishScheduledPosts()
				if err != nil {
					log.Printf("[Scheduler worker %d] Error: %v", id, err)
				} else {
					log.Printf("[Scheduler worker %d] Published %d posts", id, count)
				}
			}
		}(i)
	}

	log.Println("Scheduler started with worker pool, interval:", s.interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler shutting down...")
			close(jobs)
			wg.Done()
			return
		case <-ticker.C:
			jobs <- struct{}{}
		}
	}
}
