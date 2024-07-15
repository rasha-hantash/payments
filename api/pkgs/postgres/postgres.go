package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

// NewDBClient creates a new database client
func NewDBClient(psqlConnStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	// db.SetMaxOpenConns(25)
	// db.SetMaxIdleConns(25)
	// db.SetConnMaxLifetime(5 * time.Minute)

	// Verify the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("postgres connection success")
	return db, nil
}
