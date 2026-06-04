package api

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"ci-jarvis/internal/store/queue"
	"ci-jarvis/internal/tools/github"
	"ci-jarvis/internal/types"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func WebhookHandler(c *fiber.Ctx, q *queue.Queue) error {
	body := c.Body()

	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if secret != "" {
		sig := c.Get("X-Hub-Signature-256")
		if sig == "" {
			sig = c.Get("X-Hub-Signature")
		}
		if sig == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing signature"})
		}

		// support sha256=... and sha1=... prefixes
		var ok bool
		if strings.HasPrefix(sig, "sha256=") {
			ok = github.VerifyHMAC(body, secret, strings.TrimPrefix(sig, "sha256="), sha256.New)
		} else if strings.HasPrefix(sig, "sha1=") {
			ok = github.VerifyHMAC(body, secret, strings.TrimPrefix(sig, "sha1="), sha1.New)
		} else {
			// try to treat header as raw hex (assume sha256)
			ok = github.VerifyHMAC(body, secret, sig, sha256.New)
		}

		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid signature"})
		}
	} else {
		log.Println("warning: GITHUB_WEBHOOK_SECRET not set; skipping signature verification")
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json payload"})
	}

	job := &types.Job{
		ID:             uuid.New().String(),
		RepoURL:        github.ExtractRepoURL(payload),
		PullRequestURL: github.ExtractPRURL(payload),
		CreatedAt:      time.Now(),
	}

	if err := q.Enqueue(job); err != nil {
		log.Printf("failed to enqueue job: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to enqueue"})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "accepted", "job_id": job.ID})
}
