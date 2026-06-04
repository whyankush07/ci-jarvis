package postgres

import (
	"ci-jarvis/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (db *DB) CreateRun(ctx context.Context, j *types.Job) error {
	query := `
		INSERT INTO runs(id, repo_url, pr_url, status, current_step, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	_, err := db.client.ExecContext(ctx, query, j.ID, j.RepoURL, j.PullRequestURL, "pending", "", now, now)
	if err != nil {
		return fmt.Errorf("Failed to insert run! %v", err)
	}
	return nil
}

func (db *DB) UpdateRun(ctx context.Context, id string, status string, currentStep string) error {
	query := `
		UPDATE runs
		SET status = $1, current_step = $2, updated_at = $3
		WHERE id = $4
	`

	now := time.Now()

	_, err := db.client.ExecContext(ctx, query, status, currentStep, now, id)
	if err != nil {
		return fmt.Errorf("Failed to update run! %v", err)
	}
	return nil
}

func (db *DB) UpdateRunMetadata(ctx context.Context, id string, metadata interface{}) error {
	metaBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("Failed to marshel run metadata!! %v", err)
	}
	query := `
		UPDATE runs
		SET metadata = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err = db.client.ExecContext(ctx, query, metaBytes, now, id)
	if err != nil {
		return fmt.Errorf("Failed to add metadata to run! %v", err)
	}
	return nil
}
