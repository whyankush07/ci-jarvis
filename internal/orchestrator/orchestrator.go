package orchestrator

import (
	"context"
	"log"
	"time"

	"ci-jarvis/internal/agents"
	"ci-jarvis/internal/llm"
	"ci-jarvis/internal/store/postgres"
	"ci-jarvis/internal/store/queue"
	"ci-jarvis/internal/types"
)

type Orchestrator struct {
	queue     *queue.Queue
	db        *postgres.DB
	llmClient *llm.Client
	planner   agents.Agent
}

func NewOrchestrator(q *queue.Queue, database *postgres.DB, llmClient *llm.Client, planner agents.Agent) *Orchestrator {
	return &Orchestrator{
		queue:     q,
		db:        database,
		llmClient: llmClient,
		planner:   planner,
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
			o.processQueue(ctx)
		}
	}
}

// processQueue dequeues a job and processes it.
func (o *Orchestrator) processQueue(ctx context.Context) {
	job, err := o.queue.Dequeue()
	if err != nil {
		log.Printf("error dequeuing job: %v", err)
		return
	}

	if job == nil {
		return
	}

	log.Printf("dequeued job: id=%s, repo=%s, pr=%s", job.ID, job.RepoURL, job.PullRequestURL)

	if err := o.db.CreateRun(ctx, job); err != nil {
		log.Printf("failed to store run in database: %v", err)
		return
	}

	// Dispatch to planner agent
	run := &types.Run{
		ID:             job.ID,
		RepoURL:        job.RepoURL,
		PullRequestURL: job.PullRequestURL,
		Status:         types.RunPlanning,
	}

	result, err := o.planner.Execute(ctx, run)
	if err != nil {
		log.Printf("planner execution failed: %v", err)
		return
	}

	// Log planner output fields
	if plan, ok := result.Data.(types.PlannerOutput); ok {
		log.Printf("planner finished: success=%v, summary=%s, priority=%s, steps=%d, risks=%d",
			result.Success, plan.Summary, plan.Priority, len(plan.Steps), len(plan.Risks))
	} else {
		log.Printf("planner finished: success=%v, output_length=%d", result.Success, len(result.Output))
	}
}
