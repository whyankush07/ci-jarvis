package main

import (
	"ci-jarvis/api"
	"ci-jarvis/internal/config"
	db "ci-jarvis/internal/store/postgres"
	"ci-jarvis/internal/store/queue"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	api.RegisterRoutes(app, q)

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s\n", addr)

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- app.Listen(addr)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Printf("shutdown signal received: %v", sig)
	case err := <-srvErr:
		log.Printf("server error: %v", err)
	}

	// allow 10s for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		log.Printf("fiber shutdown error: %v", err)
	}

	if err := q.Close(); err != nil {
		log.Printf("error closing queue: %v", err)
	}

	if err := pg.Close(); err != nil {
		log.Printf("error closing db: %v", err)
	}

	log.Println("shutdown complete")
}
