package api

import (
	"ci-jarvis/internal/store/queue"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func RegisterRoutes(app *fiber.App, q *queue.Queue) error {
	app.Get("/", HealthCheck)
	app.Post("/webhook", func(c *fiber.Ctx) error {
		return WebhookHandler(c, q)
	})
	return nil
}
