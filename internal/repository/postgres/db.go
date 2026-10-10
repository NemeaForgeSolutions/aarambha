package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

// OpenDB initializes and verifies a PostgreSQL connection pool using the pgx driver.
func OpenDB(dbURL string) (*sql.DB, error) {
	if dbURL == "" {
		return nil, fmt.Errorf("database connection URL is empty")
	}

	connConfig, err := pgx.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database configuration: %w", err)
	}

	// Using QueryExecModeSimpleProtocol avoids prepared statement caching conflicts:
	// "ERROR: prepared statement already exists (SQLSTATE 42P05)"
	connConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	db := stdlib.OpenDB(*connConfig)

	// Connection pool tuning for cloud and db connections
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(15 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
