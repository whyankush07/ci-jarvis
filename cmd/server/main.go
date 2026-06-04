package main

import (
	"ci-jarvis/api"
	"ci-jarvis/internal/agents"
	"ci-jarvis/internal/config"
	"ci-jarvis/internal/llm"
	"ci-jarvis/internal/orchestrator"
	"ci-jarvis/internal/store/postgres"
	"ci-jarvis/internal/store/queue"
	"ci-jarvis/internal/tools"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file loaded: %v", err)
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	q, err := queue.NewQueue(cfg.RedisURL)
	if err != nil {
		log.Fatalf("failed to initialize queue: %v", err)
	}

	pg, err := postgres.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to initialize db connection: %v", err)
	}

	if err := pg.RunMigrations(); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	llmClient, err := llm.NewClient(context.Background(), cfg.GeminiApiKey)
	if err != nil {
		log.Fatalf("failed to initialize llm client: %v", err)
	}
	defer llmClient.Close()

	// Orchestrator will be cancelled when we receive a shutdown signal.
	orchCtx, orchCancel := context.WithCancel(context.Background())
	defer orchCancel()

	planner := agents.NewPlannerAgent(llmClient)
	coder := agents.NewCoderAgent(llmClient)
	ghTool := tools.NewGitHubTool(cfg.GithubToken)
	orch := orchestrator.NewOrchestrator(q, pg, llmClient, planner, coder, ghTool)
	go orch.Start(orchCtx)

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

	// Cancel the orchestrator context to signal it to shut down.
	orchCancel()

	// Allow 10s for graceful shutdown of all components.
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
