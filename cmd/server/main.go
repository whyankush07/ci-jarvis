package main

import (
	"ci-jarvis/api"
	"ci-jarvis/internal/config"
	db "ci-jarvis/internal/store/postgres"
	"ci-jarvis/internal/store/queue"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	q, err := queue.NewQueue(cfg.RedisURL)
	if err != nil {
		log.Fatalf("failed to initialize queue: %v", err)
	}

	pg, err := db.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to initialize db connection: %v", err)
	}

	pg.RunMigrations()

	app := fiber.New()

	app.Get("/", api.HealthCheck)
	app.Post("/webhook", func(c *fiber.Ctx) error {
		return api.WebhookHandler(c, q)
	})

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s\n", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
