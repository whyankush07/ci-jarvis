package orchestrator

import (
	"context"
	"log"
	"time"

	db "ci-jarvis/internal/store/postgres"
	"ci-jarvis/internal/store/queue"
)

type Orchestrator struct {
	queue *queue.Queue
	db    *db.DB
}

func NewOrchestrator(q *queue.Queue, database *db.DB) *Orchestrator {
	return &Orchestrator{
		queue: q,
		db:    database,
	}
}

func (o *Orchestrator) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("Orchestrator started, polling queue every 5 seconds")

	for {
		select {
		case <-ctx.Done():
			log.Println("Orchestrator shutting down")
			return
		case <-ticker.C:
			o.processQueue()
		}
	}
}

// processQueue dequeues a job and processes it.
func (o *Orchestrator) processQueue() {
	job, err := o.queue.Dequeue()
	if err != nil {
		log.Printf("error dequeuing job: %v", err)
		return
	}

	// If queue is empty, job is nil
	if job == nil {
		return
	}

	log.Printf("dequeued job: id=%s, repo=%s, pr=%s", job.ID, job.RepoURL, job.PullRequestURL)

	// TODO: Store job in Postgres
	// TODO: Dispatch to planner agent
}
