package types

import "time"

type Job struct {
	ID             string    `json:"id"`
	RepoURL        string    `json:"repo_url"`
	PullRequestURL string    `json:"pr_url"`
	CreatedAt      time.Time `json:"created_at"`
}

type Plan struct {
	ID     string   `json:"id"`
	JobID  string   `json:"job_id"`
	Steps  []string `json:"steps"`
	Status string   `json:"status"`
}

type AgentResult struct {
	AgentName string      `json:"agent_name"`
	Success   bool        `json:"success"`
	Output    string      `json:"output"`
	Data      interface{} `json:"data"`
}

type RunStatus string

const (
	RunPending          RunStatus = "pending"
	RunPlanning         RunStatus = "planning"
	RunExecuting        RunStatus = "executing"
	RunAwaitingApproval RunStatus = "awaiting_approval"
	RunCompleted        RunStatus = "completed"
	RunFailed           RunStatus = "failed"
)

type Run struct {
	ID             string      `json:"id"`
	JobID          string      `json:"job_id"`
	RepoURL        string      `json:"repo_url"`
	PullRequestURL string      `json:"pr_url"`
	Status         RunStatus   `json:"status"`
	CurrentStep    string      `json:"current_step"`
	Metadata       interface{} `json:"metadata"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type PlannerOutput struct {
	Summary  string   `json:"summary"`
	Priority string   `json:"priority"`
	Steps    []string `json:"steps"`
	Risks    []string `json:"risks"`
}
