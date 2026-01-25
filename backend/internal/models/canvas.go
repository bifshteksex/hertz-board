package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ElementType represents the type of canvas element
type ElementType string

const (
	ElementTypeText      ElementType = "text"
	ElementTypeShape     ElementType = "shape"
	ElementTypeImage     ElementType = "image"
	ElementTypeDrawing   ElementType = "drawing"
	ElementTypeSticky    ElementType = "sticky"
	ElementTypeList      ElementType = "list"
	ElementTypeConnector ElementType = "connector"
	ElementTypeGroup     ElementType = "group"
)

// Valid returns true if the element type is valid
func (t ElementType) Valid() bool {
	switch t {
	case ElementTypeText, ElementTypeShape, ElementTypeImage, ElementTypeDrawing,
		ElementTypeSticky, ElementTypeList, ElementTypeConnector, ElementTypeGroup:
		return true
	}
	return false
}

// ElementData represents the flexible JSONB structure for element properties
type ElementData map[string]interface{}

// Scan implements the sql.Scanner interface for JSONB
func (e *ElementData) Scan(value interface{}) error {
	if value == nil {
		*e = make(ElementData)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("failed to scan ElementData: unexpected type %T", value)
		}
		return json.Unmarshal([]byte(str), e)
	}
	return json.Unmarshal(bytes, e)
}

// Value implements the driver.Valuer interface for JSONB
func (e ElementData) Value() (driver.Value, error) {
	if e == nil {
		return "{}", nil
	}
	return json.Marshal(e)
}

// CanvasElement represents a canvas element in the database
type CanvasElement struct {
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
	ParentID    *uuid.UUID  `json:"parent_id,omitempty" db:"parent_id"`
	UpdatedBy   *uuid.UUID  `json:"updated_by,omitempty" db:"updated_by"`
	DeletedAt   *time.Time  `json:"deleted_at,omitempty" db:"deleted_at"`
	ElementData ElementData `json:"element_data" db:"element_data"`
	ElementType ElementType `json:"element_type" db:"element_type"`
	ZIndex      int         `json:"z_index" db:"z_index"`
	ID          uuid.UUID   `json:"id" db:"id"`
	WorkspaceID uuid.UUID   `json:"workspace_id" db:"workspace_id"`
	CreatedBy   uuid.UUID   `json:"created_by" db:"created_by"`
}

// Common element properties (for type-safe access to element_data)
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Size struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Style struct {
	Fill        string  `json:"fill,omitempty"`
	Stroke      string  `json:"stroke,omitempty"`
	FontFamily  string  `json:"font_family,omitempty"`
	FontWeight  string  `json:"font_weight,omitempty"`
	TextAlign   string  `json:"text_align,omitempty"`
	StrokeWidth float64 `json:"stroke_width,omitempty"`
	Opacity     float64 `json:"opacity,omitempty"`
	FontSize    float64 `json:"font_size,omitempty"`
}

// BaseElementData contains common fields for all element types
type BaseElementData struct {
	Style    Style    `json:"style"`
	Position Position `json:"position"`
	Size     Size     `json:"size"`
	Rotation float64  `json:"rotation"`
}

// TextElementData represents a text block
type TextElementData struct {
	Content   string `json:"content"`
	PlainText string `json:"plain_text"`
	BaseElementData
}

// ShapeElementData represents a geometric shape
type ShapeElementData struct {
	ShapeType string `json:"shape_type"`
	BaseElementData
}

// ImageElementData represents an image
type ImageElementData struct {
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	URL          string `json:"url"`
	BaseElementData
	AssetID uuid.UUID `json:"asset_id"`
}

// DrawingElementData represents freehand drawing
type DrawingElementData struct {
	Points []Point `json:"points"`
	BaseElementData
	Smooth bool `json:"smooth"`
}

type Point struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Pressure float64 `json:"pressure,omitempty"`
}

// StickyNoteData represents a sticky note
type StickyNoteData struct {
	Content string `json:"content"`
	Color   string `json:"color"`
	BaseElementData
}

// ListElementData represents a list (checklist, bullet list)
type ListElementData struct {
	Items    []ListItem `json:"items"`
	ListType string     `json:"list_type"` // bullet, numbered, checkbox
	BaseElementData
}

type ListItem struct {
	Content string    `json:"content"`
	ID      uuid.UUID `json:"id"`
	Checked bool      `json:"checked,omitempty"`
}

// ConnectorElementData represents a line connecting two elements
type ConnectorElementData struct {
	StartElementID *uuid.UUID `json:"start_element_id,omitempty"`
	EndElementID   *uuid.UUID `json:"end_element_id,omitempty"`
	StartPoint     *Position  `json:"start_point,omitempty"`
	EndPoint       *Position  `json:"end_point,omitempty"`
	LineType       string     `json:"line_type"`
	BaseElementData
	ArrowStart bool `json:"arrow_start"`
	ArrowEnd   bool `json:"arrow_end"`
}

// GroupElementData represents a group of elements
type GroupElementData struct {
	ChildIDs []uuid.UUID `json:"child_ids"`
	BaseElementData
}

// DTOs for API requests/responses

// CreateElementRequest represents a request to create a canvas element
type CreateElementRequest struct {
	ParentID    *uuid.UUID  `json:"parent_id,omitempty"`
	ElementData ElementData `json:"element_data" binding:"required"`
	ElementType ElementType `json:"element_type" binding:"required"`
	ZIndex      int         `json:"z_index"`
}

// UpdateElementRequest represents a request to update a canvas element
type UpdateElementRequest struct {
	ElementData *ElementData `json:"element_data,omitempty"`
	ZIndex      *int         `json:"z_index,omitempty"`
	ParentID    *uuid.UUID   `json:"parent_id,omitempty"`
}

// BatchCreateRequest represents a request to create multiple elements
type BatchCreateRequest struct {
	Elements []CreateElementRequest `json:"elements" binding:"required"`
}

// BatchUpdateRequest represents a request to update multiple elements
type BatchUpdateRequest struct {
	Updates []BatchUpdateItem `json:"updates" binding:"required"`
}

type BatchUpdateItem struct {
	ParentID    *uuid.UUID   `json:"parent_id,omitempty"`
	ElementData *ElementData `json:"element_data,omitempty"`
	ZIndex      *int         `json:"z_index,omitempty"`
	ID          uuid.UUID    `json:"id" binding:"required"`
}

// BatchDeleteRequest represents a request to delete multiple elements
type BatchDeleteRequest struct {
	IDs []uuid.UUID `json:"ids" binding:"required"`
}

// ElementResponse represents a canvas element in API responses
type ElementResponse struct {
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	ParentID    *uuid.UUID  `json:"parent_id,omitempty"`
	UpdatedBy   *uuid.UUID  `json:"updated_by,omitempty"`
	ElementData ElementData `json:"element_data"`
	ElementType ElementType `json:"element_type"`
	ZIndex      int         `json:"z_index"`
	ID          uuid.UUID   `json:"id"`
	WorkspaceID uuid.UUID   `json:"workspace_id"`
	CreatedBy   uuid.UUID   `json:"created_by"`
}

// ElementListResponse represents a list of canvas elements
type ElementListResponse struct {
	Elements []ElementResponse `json:"elements"`
	Total    int               `json:"total"`
}

// ToResponse converts CanvasElement to ElementResponse
func (e *CanvasElement) ToResponse() ElementResponse {
	return ElementResponse{
		ID:          e.ID,
		WorkspaceID: e.WorkspaceID,
		ElementType: e.ElementType,
		ElementData: e.ElementData,
		ZIndex:      e.ZIndex,
		ParentID:    e.ParentID,
		CreatedBy:   e.CreatedBy,
		UpdatedBy:   e.UpdatedBy,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// Canvas Snapshot Models

// CanvasSnapshot represents a version snapshot of the canvas
type CanvasSnapshot struct {
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	Description  *string     `json:"description,omitempty" db:"description"`
	SnapshotData ElementData `json:"snapshot_data" db:"snapshot_data"`
	Version      int         `json:"version" db:"version"`
	ElementCount int         `json:"element_count" db:"element_count"`
	ID           uuid.UUID   `json:"id" db:"id"`
	WorkspaceID  uuid.UUID   `json:"workspace_id" db:"workspace_id"`
	CreatedBy    uuid.UUID   `json:"created_by" db:"created_by"`
}

// CreateSnapshotRequest represents a request to create a snapshot
type CreateSnapshotRequest struct {
	Description *string `json:"description,omitempty"`
}

// SnapshotResponse represents a snapshot in API responses
type SnapshotResponse struct {
	CreatedAt    time.Time `json:"created_at"`
	Description  *string   `json:"description,omitempty"`
	Version      int       `json:"version"`
	ElementCount int       `json:"element_count"`
	ID           uuid.UUID `json:"id"`
	WorkspaceID  uuid.UUID `json:"workspace_id"`
	CreatedBy    uuid.UUID `json:"created_by"`
}

// SnapshotDetailResponse includes the full snapshot data
type SnapshotDetailResponse struct {
	SnapshotData ElementData `json:"snapshot_data"`
	SnapshotResponse
}

// SnapshotListResponse represents a list of snapshots
type SnapshotListResponse struct {
	Snapshots []SnapshotResponse `json:"snapshots"`
	Total     int                `json:"total"`
}

// ToResponse converts CanvasSnapshot to SnapshotResponse
func (s *CanvasSnapshot) ToResponse() SnapshotResponse {
	return SnapshotResponse{
		ID:           s.ID,
		WorkspaceID:  s.WorkspaceID,
		Version:      s.Version,
		Description:  s.Description,
		ElementCount: s.ElementCount,
		CreatedBy:    s.CreatedBy,
		CreatedAt:    s.CreatedAt,
	}
}

// ToDetailResponse converts CanvasSnapshot to SnapshotDetailResponse
func (s *CanvasSnapshot) ToDetailResponse() SnapshotDetailResponse {
	return SnapshotDetailResponse{
		SnapshotResponse: s.ToResponse(),
		SnapshotData:     s.SnapshotData,
	}
}
