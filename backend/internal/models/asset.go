package models

import (
	"time"

	"github.com/google/uuid"
)

// Asset represents a file asset (image, document, etc.)
type Asset struct {
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	ThumbnailURL *string    `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	Width        *int       `json:"width,omitempty" db:"width"`
	Height       *int       `json:"height,omitempty" db:"height"`
	Filename     string     `json:"filename" db:"filename"`
	ContentType  string     `json:"content_type" db:"content_type"`
	URL          string     `json:"url" db:"url"`
	Size         int64      `json:"size" db:"size"`
	ID           uuid.UUID  `json:"id" db:"id"`
	WorkspaceID  uuid.UUID  `json:"workspace_id" db:"workspace_id"`
	UploadedBy   uuid.UUID  `json:"uploaded_by" db:"uploaded_by"`
}

// UploadAssetRequest represents a file upload request
type UploadAssetRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

// AssetResponse represents an asset in API responses
type AssetResponse struct {
	CreatedAt    time.Time `json:"created_at"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	Width        *int      `json:"width,omitempty"`
	Height       *int      `json:"height,omitempty"`
	Filename     string    `json:"filename"`
	ContentType  string    `json:"content_type"`
	URL          string    `json:"url"`
	Size         int64     `json:"size"`
	ID           uuid.UUID `json:"id"`
	WorkspaceID  uuid.UUID `json:"workspace_id"`
}

// ToResponse converts Asset to AssetResponse
func (a *Asset) ToResponse() AssetResponse {
	return AssetResponse{
		ID:           a.ID,
		WorkspaceID:  a.WorkspaceID,
		Filename:     a.Filename,
		ContentType:  a.ContentType,
		Size:         a.Size,
		URL:          a.URL,
		ThumbnailURL: a.ThumbnailURL,
		Width:        a.Width,
		Height:       a.Height,
		CreatedAt:    a.CreatedAt,
	}
}
