package orchestrator

import (
	"context"
	"log"
	"time"

	"ci-jarvis/internal/agents"
	"ci-jarvis/internal/llm"
	"ci-jarvis/internal/store/postgres"
	"ci-jarvis/internal/store/queue"
	"ci-jarvis/internal/tools"
	"ci-jarvis/internal/types"
)

type Orchestrator struct {
	queue     *queue.Queue
	db        *postgres.DB
	llmClient *llm.Client
	planner   agents.Agent
	coder     agents.Agent
	github    *tools.GitHubTool
}

func NewOrchestrator(q *queue.Queue, database *postgres.DB, llmClient *llm.Client, planner agents.Agent, coder agents.Agent, github *tools.GitHubTool) *Orchestrator {
	return &Orchestrator{
		queue:     q,
		db:        database,
		llmClient: llmClient,
		planner:   planner,
		coder:     coder,
		github:    github,
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

	// Fetch PR diff for context
	log.Printf("fetching diff for PR: %s", job.PullRequestURL)
	diff, err := o.github.FetchPRDiff(job.PullRequestURL)
	if err != nil {
		log.Printf("warning: failed to fetch PR diff: %v", err)
	}

	// Dispatch to planner agent
	run := &types.Run{
		ID:             job.ID,
		RepoURL:        job.RepoURL,
		PullRequestURL: job.PullRequestURL,
		Diff:           diff,
		Status:         types.RunPlanning,
		CurrentStep:    "Analyzing PR",
	}

	// Update DB to planning
	o.db.UpdateRun(ctx, run.ID, string(types.RunPlanning), "Analyzing PR")

	plannerResult, err := o.planner.Execute(ctx, run)
	if err != nil {
		log.Printf("planner execution failed: %v", err)
		o.db.UpdateRun(ctx, run.ID, string(types.RunFailed), "Planning failed")
		return
	}

	// Log planner output fields
	if plan, ok := plannerResult.Data.(types.PlannerOutput); ok {
		log.Printf("planner finished: success=%v, summary=%s, priority=%s, steps=%d, risks=%d",
			plannerResult.Success, plan.Summary, plan.Priority, len(plan.Steps), len(plan.Risks))
		o.db.UpdateRunMetadata(ctx, run.ID, plan)
		run.Plan = &plan
	}

	// Move to Coding Phase
	run.Status = types.RunCoding
	o.db.UpdateRun(ctx, run.ID, string(types.RunCoding), "Generating code changes")

	coderResult, err := o.coder.Execute(ctx, run)
	if err != nil {
		log.Printf("coder execution failed: %v", err)
		o.db.UpdateRun(ctx, run.ID, string(types.RunFailed), "Coding failed")
		return
	}

	if coderOutput, ok := coderResult.Data.(types.CoderOutput); ok {
		log.Printf("coder finished: success=%v, explanation=%s, changes=%d",
			coderResult.Success, coderOutput.Explanation, len(coderOutput.Changes))
		o.db.UpdateRunMetadata(ctx, run.ID, coderOutput)
	}

	// For now, mark as completed after coding
	o.db.UpdateRun(ctx, run.ID, string(types.RunCompleted), "Finished processing")
}
