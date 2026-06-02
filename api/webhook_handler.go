package api

import (
	"ci-jarvis/internal/orchestrator"
	"ci-jarvis/internal/store/queue"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func WebhookHandler(c *fiber.Ctx, q *queue.Queue) error {
	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	job := &orchestrator.Job{
		ID:             uuid.New().String(),
		RepoURL:        extractRepoURL(payload),
		PullRequestURL: extractPRURL(payload),
		CreatedAt:      time.Now(),
	}

	if err := q.Enqueue(job); err != nil {
		log.Printf("failed to enqueue job: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to enqueue"})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "accepted", "job_id": job.ID})
}

func extractRepoURL(payload map[string]interface{}) string {
	if repo, ok := payload["repository"]; ok {
		if repoMap, ok := repo.(map[string]interface{}); ok {
			if url, ok := repoMap["clone_url"]; ok {
				return url.(string)
			}
		}
	}
	return ""
}

func extractPRURL(payload map[string]interface{}) string {
	if pr, ok := payload["pull_request"]; ok {
		if prMap, ok := pr.(map[string]interface{}); ok {
			if url, ok := prMap["html_url"]; ok {
				return url.(string)
			}
		}
	}
	return ""
}
