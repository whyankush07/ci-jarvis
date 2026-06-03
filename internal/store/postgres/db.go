package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	client *sql.DB
}

func NewDB(dbUrl string) (*DB, error) {

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	p := &DB{client: db}

	return p, nil
}

func (p *DB) RunMigrations() error {
	log.Println("Starting database migrations...")

	files, err := os.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		filePath := filepath.Join("migrations", file.Name())
		log.Printf("Running migration file: %s", file.Name())

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		_, err = p.client.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration query from %s: %w", file.Name(), err)
		}
	}

	log.Println("Database migrations applied successfully!")
	return nil
}

// Close closes the underlying database connection pool.
func (p *DB) Close() error {
	if p == nil || p.client == nil {
		return nil
	}
	return p.client.Close()
}
