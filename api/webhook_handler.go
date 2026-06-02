package api

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"hash"
	"log"
	"os"
	"strings"
	"time"

	"ci-jarvis/internal/store/queue"
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
			ok = verifyHMAC(body, secret, strings.TrimPrefix(sig, "sha256="), sha256.New)
		} else if strings.HasPrefix(sig, "sha1=") {
			ok = verifyHMAC(body, secret, strings.TrimPrefix(sig, "sha1="), sha1.New)
		} else {
			// try to treat header as raw hex (assume sha256)
			ok = verifyHMAC(body, secret, sig, sha256.New)
		}

		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid signature"})
		}
	} else {
		log.Println("warning: GITHUB_WEBHOOK_SECRET not set; skipping signature verification")
		//! todo -------
		//  In production, we should require a secret and reject unsigned requests.
		// add a return statement here to enforce this once we have a way to set secrets in our deployment environment.
		//!
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json payload"})
	}

	job := &types.Job{
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

// * revisit this later
func verifyHMAC(body []byte, secret, providedHex string, hashFunc func() hash.Hash) bool {
	provided, err := hex.DecodeString(providedHex)
	if err != nil {
		return false
	}
	mac := hmac.New(hashFunc, []byte(secret))
	mac.Write(body)
	expected := mac.Sum(nil)
	return hmac.Equal(expected, provided)
}

func extractRepoURL(payload map[string]interface{}) string {
	if repo, ok := payload["repository"]; ok {
		if repoMap, ok := repo.(map[string]interface{}); ok {
			if url, ok := repoMap["clone_url"]; ok {
				if s, ok := url.(string); ok {
					return s
				}
			}
		}
	}
	return ""
}

func extractPRURL(payload map[string]interface{}) string {
	if pr, ok := payload["pull_request"]; ok {
		if prMap, ok := pr.(map[string]interface{}); ok {
			if url, ok := prMap["html_url"]; ok {
				if s, ok := url.(string); ok {
					return s
				}
			}
		}
	}
	return ""
}
