package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/bifshteksex/hertz-board/internal/models"
	"github.com/bifshteksex/hertz-board/internal/repository"
)

const (
	MaxSnapshotsPerWorkspace = 100 // Keep only the latest 100 snapshots
)

type SnapshotService struct {
	snapshotRepo  *repository.SnapshotRepository
	canvasRepo    *repository.CanvasRepository
	workspaceRepo *repository.WorkspaceRepository
}

func NewSnapshotService(
	snapshotRepo *repository.SnapshotRepository,
	canvasRepo *repository.CanvasRepository,
	workspaceRepo *repository.WorkspaceRepository,
) *SnapshotService {
	return &SnapshotService{
		snapshotRepo:  snapshotRepo,
		canvasRepo:    canvasRepo,
		workspaceRepo: workspaceRepo,
	}
}

// CreateSnapshot creates a new snapshot of the current canvas state
func (s *SnapshotService) CreateSnapshot(
	ctx context.Context,
	workspaceID, userID uuid.UUID,
	description *string,
) (*models.CanvasSnapshot, error) {
	// Get all current elements
	elements, err := s.canvasRepo.GetElementsByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace elements: %w", err)
	}

	// Serialize elements to snapshot data
	snapshotData := make(models.ElementData)
	elementsJSON := make([]interface{}, len(elements))

	for i := range elements {
		elementsJSON[i] = map[string]interface{}{
			"id":           elements[i].ID,
			"element_type": elements[i].ElementType,
			"element_data": elements[i].ElementData,
			"z_index":      elements[i].ZIndex,
			"parent_id":    elements[i].ParentID,
			"created_by":   elements[i].CreatedBy,
			"updated_by":   elements[i].UpdatedBy,
			"created_at":   elements[i].CreatedAt,
			"updated_at":   elements[i].UpdatedAt,
		}
	}

	snapshotData["elements"] = elementsJSON
	snapshotData["metadata"] = map[string]interface{}{
		"element_count": len(elements),
		"created_by":    userID,
	}

	// Create snapshot
	snapshot := &models.CanvasSnapshot{
		ID:           uuid.New(),
		WorkspaceID:  workspaceID,
		Description:  description,
		SnapshotData: snapshotData,
		ElementCount: len(elements),
		CreatedBy:    userID,
	}

	if err := s.snapshotRepo.CreateSnapshot(ctx, snapshot); err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}

	// Cleanup old snapshots
	go s.cleanupOldSnapshots(context.Background(), workspaceID)

	return snapshot, nil
}

// GetSnapshot retrieves a snapshot by ID
func (s *SnapshotService) GetSnapshot(ctx context.Context, id uuid.UUID) (*models.CanvasSnapshot, error) {
	snapshot, err := s.snapshotRepo.GetSnapshotByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	return snapshot, nil
}

// GetSnapshotByVersion retrieves a snapshot by version number
func (s *SnapshotService) GetSnapshotByVersion(ctx context.Context, workspaceID uuid.UUID, version int) (*models.CanvasSnapshot, error) {
	snapshot, err := s.snapshotRepo.GetSnapshotByVersion(ctx, workspaceID, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	return snapshot, nil
}

const (
	defaultSnapshotLimit = 20
	maxSnapshotLimit     = 100
)

// ListSnapshots retrieves all snapshots for a workspace
func (s *SnapshotService) ListSnapshots(
	ctx context.Context,
	workspaceID uuid.UUID,
	limit, offset int,
) ([]models.CanvasSnapshot, int, error) {
	// Set default limit
	if limit <= 0 {
		limit = defaultSnapshotLimit
	}
	if limit > maxSnapshotLimit {
		limit = maxSnapshotLimit
	}

	snapshots, total, err := s.snapshotRepo.ListSnapshots(ctx, workspaceID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list snapshots: %w", err)
	}

	return snapshots, total, nil
}

// RestoreSnapshot restores the canvas to a specific snapshot version
func (s *SnapshotService) RestoreSnapshot(
	ctx context.Context,
	workspaceID, userID, snapshotID uuid.UUID,
) error {
	// Get the snapshot
	snapshot, err := s.snapshotRepo.GetSnapshotByID(ctx, snapshotID)
	if err != nil {
		return fmt.Errorf("snapshot not found: %w", err)
	}

	// Verify workspace
	if snapshot.WorkspaceID != workspaceID {
		return fmt.Errorf("snapshot does not belong to workspace")
	}

	// Create backup before restoring
	if err := s.createBackupSnapshot(ctx, workspaceID, userID, snapshot.Version); err != nil {
		return err
	}

	// Delete current elements
	if err := s.deleteCurrentElements(ctx, workspaceID); err != nil {
		return err
	}

	// Restore elements from snapshot
	return s.restoreElementsFromSnapshot(ctx, workspaceID, userID, snapshot)
}

func (s *SnapshotService) createBackupSnapshot(ctx context.Context, workspaceID, userID uuid.UUID, version int) error {
	desc := fmt.Sprintf("Auto-backup before restoring to version %d", version)
	if _, err := s.CreateSnapshot(ctx, workspaceID, userID, &desc); err != nil {
		return fmt.Errorf("failed to create backup snapshot: %w", err)
	}
	return nil
}

func (s *SnapshotService) deleteCurrentElements(ctx context.Context, workspaceID uuid.UUID) error {
	currentElements, err := s.canvasRepo.GetElementsByWorkspace(ctx, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to get current elements: %w", err)
	}

	if len(currentElements) > 0 {
		ids := make([]uuid.UUID, len(currentElements))
		for i := range currentElements {
			ids[i] = currentElements[i].ID
		}
		if err := s.canvasRepo.BatchDeleteElements(ctx, ids); err != nil {
			return fmt.Errorf("failed to delete current elements: %w", err)
		}
	}
	return nil
}

func (s *SnapshotService) restoreElementsFromSnapshot(
	ctx context.Context,
	workspaceID, userID uuid.UUID,
	snapshot *models.CanvasSnapshot,
) error {
	elementsData, ok := snapshot.SnapshotData["elements"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid snapshot data format")
	}

	var restoredElements []models.CanvasElement
	for _, elemData := range elementsData {
		element, err := s.parseSnapshotElement(elemData, workspaceID, userID)
		if err != nil {
			continue
		}
		restoredElements = append(restoredElements, element)
	}

	if len(restoredElements) > 0 {
		if err := s.canvasRepo.BatchCreateElements(ctx, restoredElements); err != nil {
			return fmt.Errorf("failed to restore elements: %w", err)
		}
	}
	return nil
}

func (s *SnapshotService) parseSnapshotElement(elemData interface{}, workspaceID, userID uuid.UUID) (models.CanvasElement, error) {
	elemMap, ok := elemData.(map[string]interface{})
	if !ok {
		return models.CanvasElement{}, fmt.Errorf("invalid element format")
	}

	elementDataJSON, err := json.Marshal(elemMap["element_data"])
	if err != nil {
		return models.CanvasElement{}, err
	}

	var elementData models.ElementData
	if err := json.Unmarshal(elementDataJSON, &elementData); err != nil {
		return models.CanvasElement{}, err
	}

	// Parse original created_by for audit purposes
	createdBy, _ := uuid.Parse(fmt.Sprintf("%v", elemMap["created_by"]))

	var parentID *uuid.UUID
	if pid, ok := elemMap["parent_id"]; ok && pid != nil {
		parsed, _ := uuid.Parse(fmt.Sprintf("%v", pid))
		parentID = &parsed
	}

	var zIndex int
	if zVal, ok := elemMap["z_index"].(float64); ok {
		zIndex = int(zVal)
	}

	// Generate new UUID to avoid conflicts with soft-deleted elements
	return models.CanvasElement{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		ElementType: models.ElementType(fmt.Sprintf("%v", elemMap["element_type"])),
		ElementData: elementData,
		ZIndex:      zIndex,
		ParentID:    parentID,
		CreatedBy:   createdBy,
		UpdatedBy:   &userID,
	}, nil
}

// DeleteSnapshot deletes a specific snapshot
func (s *SnapshotService) DeleteSnapshot(ctx context.Context, workspaceID, snapshotID uuid.UUID) error {
	// Verify snapshot belongs to workspace
	snapshot, err := s.snapshotRepo.GetSnapshotByID(ctx, snapshotID)
	if err != nil {
		return fmt.Errorf("snapshot not found: %w", err)
	}

	if snapshot.WorkspaceID != workspaceID {
		return fmt.Errorf("snapshot does not belong to workspace")
	}

	// Don't allow deleting the only snapshot
	count, err := s.snapshotRepo.GetSnapshotCount(ctx, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to check snapshot count: %w", err)
	}

	if count <= 1 {
		return fmt.Errorf("cannot delete the only snapshot")
	}

	if err := s.snapshotRepo.DeleteSnapshot(ctx, snapshotID); err != nil {
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}

	return nil
}

// Auto-create snapshot on significant changes (helper for future use)
func (s *SnapshotService) AutoCreateSnapshot(ctx context.Context, workspaceID, userID uuid.UUID, changeDescription string) error {
	description := fmt.Sprintf("Auto: %s", changeDescription)
	_, err := s.CreateSnapshot(ctx, workspaceID, userID, &description)
	return err
}

// Private helper functions

func (s *SnapshotService) cleanupOldSnapshots(ctx context.Context, workspaceID uuid.UUID) {
	// Keep only the latest N snapshots
	_ = s.snapshotRepo.DeleteOldSnapshots(ctx, workspaceID, MaxSnapshotsPerWorkspace)
	// Errors are intentionally ignored - cleanup is best-effort
	// In production, use proper logging
}
