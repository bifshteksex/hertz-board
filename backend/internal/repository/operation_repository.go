package repository

import (
	"context"
	"time"

	"github.com/bifshteksex/hertz-board/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OperationRepository struct {
	db *pgxpool.Pool
}

func NewOperationRepository(db *pgxpool.Pool) *OperationRepository {
	return &OperationRepository{db: db}
}

// Create stores a new operation
func (r *OperationRepository) Create(ctx context.Context, op *models.Operation) error {
	query := `
		INSERT INTO operations (
			id, workspace_id, element_id, user_id, op_type, data, timestamp, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	_, err := r.db.Exec(ctx, query,
		op.ID,
		op.WorkspaceID,
		op.ElementID,
		op.UserID,
		op.OpType,
		op.Data,
		op.Timestamp,
		op.CreatedAt,
	)

	return err
}

// GetByID retrieves an operation by ID
func (r *OperationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Operation, error) {
	query := `
		SELECT id, workspace_id, element_id, user_id, op_type, data, timestamp, created_at
		FROM operations
		WHERE id = $1
	`

	var op models.Operation
	err := r.db.QueryRow(ctx, query, id).Scan(
		&op.ID,
		&op.WorkspaceID,
		&op.ElementID,
		&op.UserID,
		&op.OpType,
		&op.Data,
		&op.Timestamp,
		&op.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &op, nil
}

// GetByWorkspaceID retrieves operations for a workspace
func (r *OperationRepository) GetByWorkspaceID(
	ctx context.Context,
	workspaceID uuid.UUID,
	limit int,
) ([]*models.Operation, error) {
	query := `
		SELECT id, workspace_id, element_id, user_id, op_type, data, timestamp, created_at
		FROM operations
		WHERE workspace_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, workspaceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	operations := make([]*models.Operation, 0)
	for rows.Next() {
		var op models.Operation
		err := rows.Scan(
			&op.ID,
			&op.WorkspaceID,
			&op.ElementID,
			&op.UserID,
			&op.OpType,
			&op.Data,
			&op.Timestamp,
			&op.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		operations = append(operations, &op)
	}

	return operations, nil
}

// GetByElementID retrieves operations for an element
func (r *OperationRepository) GetByElementID(
	ctx context.Context,
	elementID uuid.UUID,
) ([]*models.Operation, error) {
	query := `
		SELECT id, workspace_id, element_id, user_id, op_type, data, timestamp, created_at
		FROM operations
		WHERE element_id = $1
		ORDER BY timestamp ASC
	`

	rows, err := r.db.Query(ctx, query, elementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	operations := make([]*models.Operation, 0)
	for rows.Next() {
		var op models.Operation
		err := rows.Scan(
			&op.ID,
			&op.WorkspaceID,
			&op.ElementID,
			&op.UserID,
			&op.OpType,
			&op.Data,
			&op.Timestamp,
			&op.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		operations = append(operations, &op)
	}

	return operations, nil
}

// GetSince retrieves operations since a given timestamp
func (r *OperationRepository) GetSince(
	ctx context.Context,
	workspaceID uuid.UUID,
	sinceTimestamp int64,
	limit int,
) ([]*models.Operation, error) {
	query := `
		SELECT id, workspace_id, element_id, user_id, op_type, data, timestamp, created_at
		FROM operations
		WHERE workspace_id = $1 AND timestamp > $2
		ORDER BY timestamp ASC
		LIMIT $3
	`

	rows, err := r.db.Query(ctx, query, workspaceID, sinceTimestamp, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	operations := make([]*models.Operation, 0)
	for rows.Next() {
		var op models.Operation
		err := rows.Scan(
			&op.ID,
			&op.WorkspaceID,
			&op.ElementID,
			&op.UserID,
			&op.OpType,
			&op.Data,
			&op.Timestamp,
			&op.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		operations = append(operations, &op)
	}

	return operations, nil
}

// DeleteOldOperations deletes operations older than specified duration
func (r *OperationRepository) DeleteOldOperations(ctx context.Context, olderThan time.Duration) (int64, error) {
	query := `
		DELETE FROM operations
		WHERE created_at < $1
	`

	cutoffTime := time.Now().Add(-olderThan)
	result, err := r.db.Exec(ctx, query, cutoffTime)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

// GetOperationCount returns the count of operations for a workspace
func (r *OperationRepository) GetOperationCount(ctx context.Context, workspaceID uuid.UUID) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM operations
		WHERE workspace_id = $1
	`

	var count int64
	err := r.db.QueryRow(ctx, query, workspaceID).Scan(&count)
	return count, err
}
