package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Migrate runs all pending migrations
func Migrate(pool *pgxpool.Pool, migrationsPath string) error {
	ctx := context.Background()

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(ctx, pool); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	applied, err := getAppliedMigrations(ctx, pool)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Read migration files
	migrations, err := readMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if _, ok := applied[migration.Name]; ok {
			continue
		}

		fmt.Printf("Applying migration: %s\n", migration.Name)

		// Start transaction
		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		// Execute migration
		if _, err := tx.Exec(ctx, migration.SQL); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to execute migration %s: %w", migration.Name, err)
		}

		// Record migration
		if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
			migration.Version, migration.Name); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to record migration %s: %w", migration.Name, err)
		}

		// Commit transaction
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.Name, err)
		}

		fmt.Printf("Migration applied: %s\n", migration.Name)
	}

	return nil
}

type Migration struct {
	Version int
	Name    string
	SQL     string
}

func createMigrationsTable(ctx context.Context, pool *pgxpool.Pool) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`
	_, err := pool.Exec(ctx, query)
	return err
}

func getAppliedMigrations(ctx context.Context, pool *pgxpool.Pool) (map[string]bool, error) {
	applied := make(map[string]bool)

	rows, err := pool.Query(ctx, "SELECT name FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied[name] = true
	}

	return applied, rows.Err()
}

func readMigrationFiles(path string) ([]Migration, error) {
	var migrations []Migration

	files, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return migrations, nil
		}
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Parse version from filename (e.g., 001_create_users.sql)
		var version int
		var name string
		if _, err := fmt.Sscanf(file.Name(), "%d_%s", &version, &name); err != nil {
			continue
		}

		// Read SQL content
		content, err := os.ReadFile(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    strings.TrimSuffix(file.Name(), ".sql"),
			SQL:     string(content),
		})
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}
