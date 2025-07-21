package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"item-pdp-service/internal/infrastructure/config"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// DB wraps sql.DB to provide additional functionality
type DB struct {
	*sql.DB
}

// NewConnection creates a new database connection
func NewConnection(config *config.Config) (*DB, error) {
	db, err := sql.Open("postgres", config.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.Database.MaxOpenConns)
	db.SetMaxIdleConns(config.Database.MaxIdleConns)
	db.SetConnMaxLifetime(config.Database.ConnMaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().
		Str("host", config.Database.Host).
		Int("port", config.Database.Port).
		Str("database", config.Database.DBName).
		Msg("Connected to database")

	return &DB{DB: db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.DB != nil {
		log.Info().Msg("Closing database connection")
		return db.DB.Close()
	}
	return nil
}

// Health checks the database connection health
func (db *DB) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// WithTransaction executes a function within a database transaction
func (db *DB) WithTransaction(fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error().Err(rollbackErr).Msg("Failed to rollback transaction")
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
} 