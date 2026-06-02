package api

import "github.com/gofiber/fiber/v2"

func WebhookHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "accepted"})
}
