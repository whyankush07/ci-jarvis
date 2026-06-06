package main

import (
	"ci-jarvis/api"
	"ci-jarvis/internal/agents"
	"ci-jarvis/internal/config"
	"ci-jarvis/internal/llm"
	"ci-jarvis/internal/orchestrator"
	"ci-jarvis/internal/store/postgres"
	"ci-jarvis/internal/store/queue"
	"ci-jarvis/internal/store/vector"
	"ci-jarvis/internal/tools/github"
	"ci-jarvis/internal/tools/rag"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	// if err := pg.CreateMoc(context.Background()); err != nil {
	// 	log.Fatalf("Error in seeding database: %v", err)

	// }

	llmClient, err := llm.NewClient(context.Background(), cfg.GeminiApiKey)
	if err != nil {
		log.Fatalf("failed to initialize llm client: %v", err)
	}
	defer llmClient.Close()

	vStore, err := vector.NewQdrantStore(cfg.QdrantURL, cfg.QdrantApiKey, "code_snippets")
	if err != nil {
		log.Printf("warning: failed to initialize vector store: %v", err)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := vStore.CreateCollection(ctx, 768); err != nil {
			log.Printf("info: collection might already exist or failed to create: %v", err)
		}
		cancel()

		go func() {
			crawler := rag.NewCrawler(llmClient, vStore)
			log.Println("Starting initial codebase indexing...")
			if err := crawler.IndexRepository(context.Background(), "."); err != nil {
				log.Printf("error indexing repository: %v", err)
			} else {
				log.Println("Codebase indexing complete")
			}
		}()
	}

	orchCtx, orchCancel := context.WithCancel(context.Background())
	defer orchCancel()

	planner := agents.NewPlannerAgent(llmClient, vStore)
	coder := agents.NewCoderAgent(llmClient)
	reviewer := agents.NewReviewerAgent(llmClient)
	ghTool := github.NewGitHubTool(cfg.GithubToken)
	orch := orchestrator.NewOrchestrator(q, pg, llmClient, planner, coder, reviewer, ghTool)
	go orch.Start(orchCtx)

	app := fiber.New()

	origins := "http://localhost:8080"
	if cfg.AllowedOrigins != "" {
		origins += "," + cfg.AllowedOrigins
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: origins,
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	api.RegisterRoutes(app, q, pg)

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
