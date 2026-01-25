package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bifshteksex/hertz-board/internal/models"
)

type SnapshotRepository struct {
	db *pgxpool.Pool
}

func NewSnapshotRepository(db *pgxpool.Pool) *SnapshotRepository {
	return &SnapshotRepository{db: db}
}

// CreateSnapshot creates a new canvas snapshot
func (r *SnapshotRepository) CreateSnapshot(ctx context.Context, snapshot *models.CanvasSnapshot) error {
	query := `
		INSERT INTO canvas_snapshots (
			id, workspace_id, version, description, snapshot_data, element_count, created_by
		) VALUES ($1, $2, get_next_snapshot_version($2), $3, $4, $5, $6)
		RETURNING version, created_at
	`

	return r.db.QueryRow(ctx, query,
		snapshot.ID,
		snapshot.WorkspaceID,
		snapshot.Description,
		snapshot.SnapshotData,
		snapshot.ElementCount,
		snapshot.CreatedBy,
	).Scan(&snapshot.Version, &snapshot.CreatedAt)
}

// scanSnapshot scans a row into a CanvasSnapshot
func (r *SnapshotRepository) scanSnapshot(row pgx.Row) (*models.CanvasSnapshot, error) {
	var snapshot models.CanvasSnapshot
	err := row.Scan(
		&snapshot.ID,
		&snapshot.WorkspaceID,
		&snapshot.Version,
		&snapshot.Description,
		&snapshot.SnapshotData,
		&snapshot.ElementCount,
		&snapshot.CreatedBy,
		&snapshot.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("snapshot not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan snapshot: %w", err)
	}

	return &snapshot, nil
}

// GetSnapshotByID retrieves a snapshot by ID
func (r *SnapshotRepository) GetSnapshotByID(ctx context.Context, id uuid.UUID) (*models.CanvasSnapshot, error) {
	query := `
		SELECT id, workspace_id, version, description, snapshot_data, element_count, created_by, created_at
		FROM canvas_snapshots
		WHERE id = $1
	`

	return r.scanSnapshot(r.db.QueryRow(ctx, query, id))
}

// GetSnapshotByVersion retrieves a snapshot by workspace and version number
func (r *SnapshotRepository) GetSnapshotByVersion(ctx context.Context, workspaceID uuid.UUID, version int) (*models.CanvasSnapshot, error) {
	query := `
		SELECT id, workspace_id, version, description, snapshot_data, element_count, created_by, created_at
		FROM canvas_snapshots
		WHERE workspace_id = $1 AND version = $2
	`

	return r.scanSnapshot(r.db.QueryRow(ctx, query, workspaceID, version))
}

// GetLatestSnapshot retrieves the latest snapshot for a workspace
func (r *SnapshotRepository) GetLatestSnapshot(ctx context.Context, workspaceID uuid.UUID) (*models.CanvasSnapshot, error) {
	query := `
		SELECT id, workspace_id, version, description, snapshot_data, element_count, created_by, created_at
		FROM canvas_snapshots
		WHERE workspace_id = $1
		ORDER BY version DESC
		LIMIT 1
	`

	return r.scanSnapshot(r.db.QueryRow(ctx, query, workspaceID))
}

// ListSnapshots retrieves all snapshots for a workspace with pagination
func (r *SnapshotRepository) ListSnapshots(
	ctx context.Context,
	workspaceID uuid.UUID,
	limit, offset int,
) ([]models.CanvasSnapshot, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM canvas_snapshots WHERE workspace_id = $1`
	if err := r.db.QueryRow(ctx, countQuery, workspaceID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count snapshots: %w", err)
	}

	// Get snapshots
	query := `
		SELECT id, workspace_id, version, description, snapshot_data, element_count, created_by, created_at
		FROM canvas_snapshots
		WHERE workspace_id = $1
		ORDER BY version DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, workspaceID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list snapshots: %w", err)
	}
	defer rows.Close()

	var snapshots []models.CanvasSnapshot
	for rows.Next() {
		var snapshot models.CanvasSnapshot
		err := rows.Scan(
			&snapshot.ID,
			&snapshot.WorkspaceID,
			&snapshot.Version,
			&snapshot.Description,
			&snapshot.SnapshotData,
			&snapshot.ElementCount,
			&snapshot.CreatedBy,
			&snapshot.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan snapshot: %w", err)
		}
		snapshots = append(snapshots, snapshot)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating snapshots: %w", err)
	}

	return snapshots, total, nil
}

// DeleteOldSnapshots deletes old snapshots keeping only the latest N versions
func (r *SnapshotRepository) DeleteOldSnapshots(ctx context.Context, workspaceID uuid.UUID, keepCount int) error {
	query := `
		DELETE FROM canvas_snapshots
		WHERE workspace_id = $1
		  AND version < (
		      SELECT MAX(version) - $2
		      FROM canvas_snapshots
		      WHERE workspace_id = $1
		  )
	`

	_, err := r.db.Exec(ctx, query, workspaceID, keepCount)
	if err != nil {
		return fmt.Errorf("failed to delete old snapshots: %w", err)
	}

	return nil
}

// GetSnapshotCount returns the total number of snapshots for a workspace
func (r *SnapshotRepository) GetSnapshotCount(ctx context.Context, workspaceID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM canvas_snapshots WHERE workspace_id = $1`

	err := r.db.QueryRow(ctx, query, workspaceID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count snapshots: %w", err)
	}

	return count, nil
}

// DeleteSnapshot deletes a specific snapshot
func (r *SnapshotRepository) DeleteSnapshot(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM canvas_snapshots WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("snapshot not found")
	}

	return nil
}
