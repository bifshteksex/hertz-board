package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bifshteksex/hertz-board/internal/models"
)

type CanvasRepository struct {
	db *pgxpool.Pool
}

func NewCanvasRepository(db *pgxpool.Pool) *CanvasRepository {
	return &CanvasRepository{db: db}
}

// CreateElement creates a new canvas element
func (r *CanvasRepository) CreateElement(ctx context.Context, element *models.CanvasElement) error {
	query := `
		INSERT INTO canvas_elements (
			id, workspace_id, element_type, element_data, z_index, parent_id, created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`

	return r.db.QueryRow(ctx, query,
		element.ID,
		element.WorkspaceID,
		element.ElementType,
		element.ElementData,
		element.ZIndex,
		element.ParentID,
		element.CreatedBy,
		element.UpdatedBy,
	).Scan(&element.CreatedAt, &element.UpdatedAt)
}

// GetElementByID retrieves a canvas element by ID
func (r *CanvasRepository) GetElementByID(ctx context.Context, id uuid.UUID) (*models.CanvasElement, error) {
	query := `
		SELECT id, workspace_id, element_type, element_data, z_index, parent_id,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM canvas_elements
		WHERE id = $1 AND deleted_at IS NULL
	`

	var element models.CanvasElement
	err := r.db.QueryRow(ctx, query, id).Scan(
		&element.ID,
		&element.WorkspaceID,
		&element.ElementType,
		&element.ElementData,
		&element.ZIndex,
		&element.ParentID,
		&element.CreatedBy,
		&element.UpdatedBy,
		&element.CreatedAt,
		&element.UpdatedAt,
		&element.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("element not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get element: %w", err)
	}

	return &element, nil
}

// GetElementsByWorkspace retrieves all elements for a workspace
func (r *CanvasRepository) GetElementsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.CanvasElement, error) {
	query := `
		SELECT id, workspace_id, element_type, element_data, z_index, parent_id,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM canvas_elements
		WHERE workspace_id = $1 AND deleted_at IS NULL
		ORDER BY z_index ASC, created_at ASC
	`

	rows, err := r.db.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query elements: %w", err)
	}
	defer rows.Close()

	var elements []models.CanvasElement
	for rows.Next() {
		var element models.CanvasElement
		err := rows.Scan(
			&element.ID,
			&element.WorkspaceID,
			&element.ElementType,
			&element.ElementData,
			&element.ZIndex,
			&element.ParentID,
			&element.CreatedBy,
			&element.UpdatedBy,
			&element.CreatedAt,
			&element.UpdatedAt,
			&element.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan element: %w", err)
		}
		elements = append(elements, element)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating elements: %w", err)
	}

	return elements, nil
}

// UpdateElement updates a canvas element
func (r *CanvasRepository) UpdateElement(ctx context.Context, element *models.CanvasElement) error {
	query := `
		UPDATE canvas_elements
		SET element_data = $1, z_index = $2, parent_id = $3, updated_by = $4, updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		element.ElementData,
		element.ZIndex,
		element.ParentID,
		element.UpdatedBy,
		element.ID,
	).Scan(&element.UpdatedAt)

	if err == pgx.ErrNoRows {
		return fmt.Errorf("element not found or already deleted")
	}
	if err != nil {
		return fmt.Errorf("failed to update element: %w", err)
	}

	return nil
}

// DeleteElement soft deletes a canvas element
func (r *CanvasRepository) DeleteElement(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE canvas_elements
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete element: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("element not found or already deleted")
	}

	return nil
}

// HardDeleteElement permanently deletes a canvas element
func (r *CanvasRepository) HardDeleteElement(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM canvas_elements WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete element: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("element not found")
	}

	return nil
}

// Batch operations

// BatchCreateElements creates multiple canvas elements in a transaction
func (r *CanvasRepository) BatchCreateElements(ctx context.Context, elements []models.CanvasElement) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := `
		INSERT INTO canvas_elements (
			id, workspace_id, element_type, element_data, z_index, parent_id, created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`

	for i := range elements {
		err := tx.QueryRow(ctx, query,
			elements[i].ID,
			elements[i].WorkspaceID,
			elements[i].ElementType,
			elements[i].ElementData,
			elements[i].ZIndex,
			elements[i].ParentID,
			elements[i].CreatedBy,
			elements[i].UpdatedBy,
		).Scan(&elements[i].CreatedAt, &elements[i].UpdatedAt)

		if err != nil {
			return fmt.Errorf("failed to create element %d: %w", i, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// BatchUpdateElements updates multiple canvas elements in a transaction
func (r *CanvasRepository) BatchUpdateElements(ctx context.Context, elements []models.CanvasElement) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := `
		UPDATE canvas_elements
		SET element_data = $1, z_index = $2, parent_id = $3, updated_by = $4, updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING updated_at
	`

	for i := range elements {
		err := tx.QueryRow(ctx, query,
			elements[i].ElementData,
			elements[i].ZIndex,
			elements[i].ParentID,
			elements[i].UpdatedBy,
			elements[i].ID,
		).Scan(&elements[i].UpdatedAt)

		if err == pgx.ErrNoRows {
			return fmt.Errorf("element %s not found or already deleted", elements[i].ID)
		}
		if err != nil {
			return fmt.Errorf("failed to update element %d: %w", i, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// BatchDeleteElements soft deletes multiple canvas elements in a transaction
func (r *CanvasRepository) BatchDeleteElements(ctx context.Context, ids []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := `
		UPDATE canvas_elements
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	for _, id := range ids {
		result, err := tx.Exec(ctx, query, id)
		if err != nil {
			return fmt.Errorf("failed to delete element %s: %w", id, err)
		}
		if result.RowsAffected() == 0 {
			return fmt.Errorf("element %s not found or already deleted", id)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetElementCount returns the total number of elements in a workspace
func (r *CanvasRepository) GetElementCount(ctx context.Context, workspaceID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM canvas_elements
		WHERE workspace_id = $1 AND deleted_at IS NULL
	`

	var count int
	err := r.db.QueryRow(ctx, query, workspaceID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count elements: %w", err)
	}

	return count, nil
}

// GetElementsByType retrieves all elements of a specific type in a workspace
func (r *CanvasRepository) GetElementsByType(
	ctx context.Context,
	workspaceID uuid.UUID,
	elementType models.ElementType,
) ([]models.CanvasElement, error) {
	query := `
		SELECT id, workspace_id, element_type, element_data, z_index, parent_id,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM canvas_elements
		WHERE workspace_id = $1 AND element_type = $2 AND deleted_at IS NULL
		ORDER BY z_index ASC, created_at ASC
	`

	rows, err := r.db.Query(ctx, query, workspaceID, elementType)
	if err != nil {
		return nil, fmt.Errorf("failed to query elements by type: %w", err)
	}
	defer rows.Close()

	var elements []models.CanvasElement
	for rows.Next() {
		var element models.CanvasElement
		err := rows.Scan(
			&element.ID,
			&element.WorkspaceID,
			&element.ElementType,
			&element.ElementData,
			&element.ZIndex,
			&element.ParentID,
			&element.CreatedBy,
			&element.UpdatedBy,
			&element.CreatedAt,
			&element.UpdatedAt,
			&element.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan element: %w", err)
		}
		elements = append(elements, element)
	}

	return elements, rows.Err()
}

// GetChildElements retrieves all child elements of a parent (for groups)
func (r *CanvasRepository) GetChildElements(ctx context.Context, parentID uuid.UUID) ([]models.CanvasElement, error) {
	query := `
		SELECT id, workspace_id, element_type, element_data, z_index, parent_id,
		       created_by, updated_by, created_at, updated_at, deleted_at
		FROM canvas_elements
		WHERE parent_id = $1 AND deleted_at IS NULL
		ORDER BY z_index ASC
	`

	rows, err := r.db.Query(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query child elements: %w", err)
	}
	defer rows.Close()

	var elements []models.CanvasElement
	for rows.Next() {
		var element models.CanvasElement
		err := rows.Scan(
			&element.ID,
			&element.WorkspaceID,
			&element.ElementType,
			&element.ElementData,
			&element.ZIndex,
			&element.ParentID,
			&element.CreatedBy,
			&element.UpdatedBy,
			&element.CreatedAt,
			&element.UpdatedAt,
			&element.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan child element: %w", err)
		}
		elements = append(elements, element)
	}

	return elements, rows.Err()
}

// DeleteWorkspaceElements deletes all elements in a workspace (for workspace deletion)
func (r *CanvasRepository) DeleteWorkspaceElements(ctx context.Context, workspaceID uuid.UUID) error {
	query := `
		UPDATE canvas_elements
		SET deleted_at = NOW()
		WHERE workspace_id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to delete workspace elements: %w", err)
	}

	return nil
}
