package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitDB initializes and returns a connection pool for PostgreSQL.
func InitDB(ctx context.Context, connStr string) *pgxpool.Pool {
	dbpool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL")

	return dbpool
}
