package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bifshteksex/hertzboard/internal/config"
)

// NewPostgresPool creates a new PostgreSQL connection pool
func NewPostgresPool(cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure pool settings
	maxConns := cfg.MaxConnections
	if maxConns > 0x7FFFFFFF {
		maxConns = 0x7FFFFFFF
	}
	poolConfig.MaxConns = int32(maxConns)

	minConns := cfg.MaxIdleConnections
	if minConns > 0x7FFFFFFF {
		minConns = 0x7FFFFFFF
	}
	poolConfig.MinConns = int32(minConns)
	poolConfig.MaxConnLifetime = time.Duration(cfg.ConnectionMaxLifetime) * time.Second
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	// Create pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// ClosePostgresPool closes the database connection pool
func ClosePostgresPool(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
