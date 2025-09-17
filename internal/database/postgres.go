package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDb struct {
	Pool *pgxpool.Pool
}

func NewPostgresPool(ctx context.Context, url string) *PostgresDb {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		log.Fatalln("Error connecting to PostgresDB", err.Error())
	}
	// Execute schema.sql
	schemaPath := "sql/schema.sql"
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Failed to read schema.sql: %v", err)
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("Failed to acquire connection for schema migration: %v", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, string(schemaBytes))
	if err != nil {
		log.Fatalf("Failed to execute schema.sql: %v", err)
	}

	return &PostgresDb{Pool: pool}
}
