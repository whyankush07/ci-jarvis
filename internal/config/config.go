package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port         string
	DatabaseURL  string
	RedisURL     string
	QdrantURL    string
	GeminiApiKey string
	GithubToken  string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		RedisURL:     os.Getenv("REDIS_URL"),
		QdrantURL:    os.Getenv("QDRANT_URL"),
		GeminiApiKey: os.Getenv("GEMINI_API_KEY"),
		GithubToken:  os.Getenv("GITHUB_TOKEN"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.RedisURL == "" {
		return nil, fmt.Errorf("REDIS_URL is required")
	}
	// if cfg.QdrantURL == "" {
	// 	return nil, fmt.Errorf("QDRANT_URL is required")
	// }
	if cfg.GeminiApiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
