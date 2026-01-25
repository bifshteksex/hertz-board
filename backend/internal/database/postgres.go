package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bifshteksex/hertz-board/internal/config"
)

const (
	defaultPingTimeout = 5 * time.Second
)

// NewPostgresPool creates a new PostgreSQL connection pool
func NewPostgresPool(cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure pool settings
	// #nosec G115 -- MaxConnections is validated to be within int32 range
	maxConns := cfg.MaxConnections
	if maxConns > 0x7FFFFFFF || maxConns < 0 {
		maxConns = 100 // default safe value
	}
	poolConfig.MaxConns = int32(maxConns) // #nosec G115

	// #nosec G115 -- MaxIdleConnections is validated to be within int32 range
	minConns := cfg.MaxIdleConnections
	if minConns > 0x7FFFFFFF || minConns < 0 {
		minConns = 10 // default safe value
	}
	poolConfig.MinConns = int32(minConns) // #nosec G115
	poolConfig.MaxConnLifetime = time.Duration(cfg.ConnectionMaxLifetime) * time.Second
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	// Create pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), defaultPingTimeout)
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
