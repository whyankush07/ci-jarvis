package orchestrator

import (
	"time"
)

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
	AgentName string `json:"agent_name"`
	Success   bool   `json:"success"`
	Output    string `json:"output"`
}
