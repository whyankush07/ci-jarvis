package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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

	// // Run migrations (stub for now)
	// if err := p.RunMigrations(); err != nil {
	// 	log.Printf("warning: migrations failed: %v", err)
	// }

	return p, nil
}

func (p *DB) RunMigrations() error {
	//!
	//  integrate a migration tool (golang-migrate) later.
	//!
	log.Println("RunMigrations: no migrations configured (stub)")
	return nil
}
