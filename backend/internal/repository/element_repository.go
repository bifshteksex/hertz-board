package repository

import (
	"context"
	"time"

	"github.com/bifshteksex/hertz-board/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ElementRepository struct {
	db *pgxpool.Pool
}

func NewElementRepository(db *pgxpool.Pool) *ElementRepository {
	return &ElementRepository{db: db}
}

// Create creates a new element
func (r *ElementRepository) Create(ctx context.Context, element *models.Element) error {
	query := `
		INSERT INTO elements (
			id, workspace_id, type, content, pos_x, pos_y, width, height,
			z_index, rotation, style, version, created_by, updated_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`

	now := time.Now()
	if element.CreatedAt.IsZero() {
		element.CreatedAt = now
	}
	if element.UpdatedAt.IsZero() {
		element.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		element.ID,
		element.WorkspaceID,
		element.Type,
		element.Content,
		element.PosX,
		element.PosY,
		element.Width,
		element.Height,
		element.ZIndex,
		element.Rotation,
		element.Style,
		element.Version,
		element.CreatedBy,
		element.UpdatedBy,
		element.CreatedAt,
		element.UpdatedAt,
	)

	return err
}

// GetByID retrieves an element by ID
func (r *ElementRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Element, error) {
	query := `
		SELECT id, workspace_id, type, content, pos_x, pos_y, width, height,
			z_index, rotation, style, version, created_by, updated_by, created_at, updated_at, deleted_at
		FROM elements
		WHERE id = $1 AND deleted_at IS NULL
	`

	var element models.Element
	err := r.db.QueryRow(ctx, query, id).Scan(
		&element.ID,
		&element.WorkspaceID,
		&element.Type,
		&element.Content,
		&element.PosX,
		&element.PosY,
		&element.Width,
		&element.Height,
		&element.ZIndex,
		&element.Rotation,
		&element.Style,
		&element.Version,
		&element.CreatedBy,
		&element.UpdatedBy,
		&element.CreatedAt,
		&element.UpdatedAt,
		&element.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &element, nil
}

// Update updates an element
func (r *ElementRepository) Update(ctx context.Context, element *models.Element) error {
	query := `
		UPDATE elements
		SET content = $1, pos_x = $2, pos_y = $3, width = $4, height = $5,
			z_index = $6, rotation = $7, style = $8, version = $9, updated_by = $10, updated_at = $11
		WHERE id = $12 AND deleted_at IS NULL
	`

	element.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		element.Content,
		element.PosX,
		element.PosY,
		element.Width,
		element.Height,
		element.ZIndex,
		element.Rotation,
		element.Style,
		element.Version,
		element.UpdatedBy,
		element.UpdatedAt,
		element.ID,
	)

	return err
}

// Delete soft deletes an element
func (r *ElementRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE elements
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}

// GetByWorkspaceID retrieves all elements for a workspace
func (r *ElementRepository) GetByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]*models.Element, error) {
	query := `
		SELECT id, workspace_id, type, content, pos_x, pos_y, width, height,
			z_index, rotation, style, version, created_by, updated_by, created_at, updated_at
		FROM elements
		WHERE workspace_id = $1 AND deleted_at IS NULL
		ORDER BY z_index ASC, created_at ASC
	`

	rows, err := r.db.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	elements := make([]*models.Element, 0)
	for rows.Next() {
		var element models.Element
		err := rows.Scan(
			&element.ID,
			&element.WorkspaceID,
			&element.Type,
			&element.Content,
			&element.PosX,
			&element.PosY,
			&element.Width,
			&element.Height,
			&element.ZIndex,
			&element.Rotation,
			&element.Style,
			&element.Version,
			&element.CreatedBy,
			&element.UpdatedBy,
			&element.CreatedAt,
			&element.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		elements = append(elements, &element)
	}

	return elements, nil
}
