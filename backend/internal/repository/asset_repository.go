package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bifshteksex/hertz-board/internal/models"
)

type AssetRepository struct {
	db *pgxpool.Pool
}

func NewAssetRepository(db *pgxpool.Pool) *AssetRepository {
	return &AssetRepository{db: db}
}

// CreateAsset creates a new asset record
func (r *AssetRepository) CreateAsset(ctx context.Context, asset *models.Asset) error {
	query := `
		INSERT INTO assets (
			id, workspace_id, uploaded_by, filename, content_type, size, url, thumbnail_url, width, height
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at
	`

	return r.db.QueryRow(ctx, query,
		asset.ID,
		asset.WorkspaceID,
		asset.UploadedBy,
		asset.Filename,
		asset.ContentType,
		asset.Size,
		asset.URL,
		asset.ThumbnailURL,
		asset.Width,
		asset.Height,
	).Scan(&asset.CreatedAt)
}

// GetAssetByID retrieves an asset by ID
func (r *AssetRepository) GetAssetByID(ctx context.Context, id uuid.UUID) (*models.Asset, error) {
	query := `
		SELECT id, workspace_id, uploaded_by, filename, content_type, size, url, thumbnail_url, width, height, created_at, deleted_at
		FROM assets
		WHERE id = $1 AND deleted_at IS NULL
	`

	return r.scanAsset(r.db.QueryRow(ctx, query, id))
}

// scanAsset scans a single row into an Asset
func (r *AssetRepository) scanAsset(row pgx.Row) (*models.Asset, error) {
	var asset models.Asset
	err := row.Scan(
		&asset.ID,
		&asset.WorkspaceID,
		&asset.UploadedBy,
		&asset.Filename,
		&asset.ContentType,
		&asset.Size,
		&asset.URL,
		&asset.ThumbnailURL,
		&asset.Width,
		&asset.Height,
		&asset.CreatedAt,
		&asset.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("asset not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan asset: %w", err)
	}

	return &asset, nil
}

// scanAssets scans multiple rows into Assets
func (r *AssetRepository) scanAssets(rows pgx.Rows) ([]models.Asset, error) {
	var assets []models.Asset
	for rows.Next() {
		var asset models.Asset
		err := rows.Scan(
			&asset.ID,
			&asset.WorkspaceID,
			&asset.UploadedBy,
			&asset.Filename,
			&asset.ContentType,
			&asset.Size,
			&asset.URL,
			&asset.ThumbnailURL,
			&asset.Width,
			&asset.Height,
			&asset.CreatedAt,
			&asset.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}
		assets = append(assets, asset)
	}

	return assets, rows.Err()
}

// GetAssetsByWorkspace retrieves all assets for a workspace
func (r *AssetRepository) GetAssetsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.Asset, error) {
	query := `
		SELECT id, workspace_id, uploaded_by, filename, content_type, size, url, thumbnail_url, width, height, created_at, deleted_at
		FROM assets
		WHERE workspace_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query assets: %w", err)
	}
	defer rows.Close()

	return r.scanAssets(rows)
}

// DeleteAsset soft deletes an asset
func (r *AssetRepository) DeleteAsset(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE assets
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("asset not found or already deleted")
	}

	return nil
}

// GetOrphanedAssets retrieves assets that are not referenced by any canvas element
func (r *AssetRepository) GetOrphanedAssets(ctx context.Context, workspaceID uuid.UUID) ([]models.Asset, error) {
	query := `
		SELECT a.id, a.workspace_id, a.uploaded_by, a.filename, a.content_type,
		       a.size, a.url, a.thumbnail_url, a.width, a.height,
		       a.created_at, a.deleted_at
		FROM assets a
		WHERE a.workspace_id = $1
		  AND a.deleted_at IS NULL
		  AND NOT EXISTS (
		      SELECT 1 FROM canvas_elements ce
		      WHERE ce.workspace_id = a.workspace_id
		        AND ce.deleted_at IS NULL
		        AND ce.element_type = 'image'
		        AND ce.element_data->>'asset_id' = a.id::text
		  )
		  AND a.created_at < NOW() - INTERVAL '1 hour' -- Grace period for upload
	`

	rows, err := r.db.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query orphaned assets: %w", err)
	}
	defer rows.Close()

	return r.scanAssets(rows)
}
