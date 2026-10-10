package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"os"
)

//go:embed schema.sql
var embeddedSchemaSQL string

//go:embed seed.sql
var embeddedSeedSQL string

const InitialSeedMigration = "initial_seed"

// ExecuteDDL executes table DDL migrations and ensures initial seed data
// is only populated once on a fresh/empty database, preventing data from being
// re-seeded on application restarts.
func ExecuteDDL(ctx context.Context, db *sql.DB) error {
	log.Printf("Verifying database schema...")
	if _, err := db.ExecContext(ctx, embeddedSchemaSQL); err != nil {
		return err
	}
	log.Printf("Database schema verified successfully.")

	// Allow manual override if explicit re-seeding is desired
	if os.Getenv("FORCE_SEED") == "true" {
		log.Printf("FORCE_SEED=true: applying seed data...")
		if _, err := db.ExecContext(ctx, embeddedSeedSQL); err != nil {
			return err
		}
		_, _ = db.ExecContext(ctx, `
			INSERT INTO schema_migrations (version) VALUES ($1)
			ON CONFLICT (version) DO NOTHING
		`, InitialSeedMigration)
		log.Printf("Seed data applied successfully.")
		return nil
	}

	seeded, err := isDatabaseSeeded(ctx, db)
	if err != nil {
		return err
	}

	if seeded {
		log.Printf("Database already initialized; skipping seed data.")
		return nil
	}

	log.Printf("Empty database detected; applying initial seed data...")
	if _, err := db.ExecContext(ctx, embeddedSeedSQL); err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO schema_migrations (version) VALUES ($1)
		ON CONFLICT (version) DO NOTHING
	`, InitialSeedMigration)
	if err != nil {
		log.Printf("Warning: failed to record seed migration: %v", err)
	}

	log.Printf("Initial seed data applied successfully.")
	return nil
}

// isDatabaseSeeded checks if the database has already been seeded.
func isDatabaseSeeded(ctx context.Context, db *sql.DB) (bool, error) {
	// 1. Check if seed migration has already been recorded
	var migrationExists bool
	err := db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM schema_migrations WHERE version = $1
		)`, InitialSeedMigration).Scan(&migrationExists)
	if err == nil && migrationExists {
		return true, nil
	}

	// 2. For databases created before schema_migrations existed,
	// check if core tables already have any records.
	var hasData bool
	err = db.QueryRowContext(ctx, `
		SELECT (
			EXISTS (SELECT 1 FROM properties LIMIT 1) OR
			EXISTS (SELECT 1 FROM agents LIMIT 1) OR
			EXISTS (SELECT 1 FROM locations LIMIT 1) OR
			EXISTS (SELECT 1 FROM users LIMIT 1)
		)`).Scan(&hasData)
	if err != nil {
		return false, err
	}

	if hasData {
		// Existing database: mark as seeded so future restarts don't re-seed
		_, _ = db.ExecContext(ctx, `
			INSERT INTO schema_migrations (version) VALUES ($1)
			ON CONFLICT (version) DO NOTHING
		`, InitialSeedMigration)
		return true, nil
	}

	return false, nil
}
