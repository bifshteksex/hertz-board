package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/bifshteksex/hertz-board/internal/models"
	"github.com/bifshteksex/hertz-board/internal/repository"
)

type CanvasService struct {
	canvasRepo    *repository.CanvasRepository
	workspaceRepo *repository.WorkspaceRepository
	cacheService  *CanvasCacheService
}

func NewCanvasService(
	canvasRepo *repository.CanvasRepository,
	workspaceRepo *repository.WorkspaceRepository,
	cacheService *CanvasCacheService,
) *CanvasService {
	return &CanvasService{
		canvasRepo:    canvasRepo,
		workspaceRepo: workspaceRepo,
		cacheService:  cacheService,
	}
}

// CreateElement creates a new canvas element
func (s *CanvasService) CreateElement(
	ctx context.Context,
	workspaceID, userID uuid.UUID,
	req models.CreateElementRequest,
) (*models.CanvasElement, error) {
	// Validate element type
	if !req.ElementType.Valid() {
		return nil, fmt.Errorf("invalid element type: %s", req.ElementType)
	}

	// Validate element data
	if len(req.ElementData) == 0 {
		return nil, fmt.Errorf("element_data is required")
	}

	// Create element
	element := &models.CanvasElement{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		ElementType: req.ElementType,
		ElementData: req.ElementData,
		ZIndex:      req.ZIndex,
		ParentID:    req.ParentID,
		CreatedBy:   userID,
		UpdatedBy:   &userID,
	}

	// Validate parent exists if specified
	if req.ParentID != nil {
		parent, err := s.canvasRepo.GetElementByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent element not found: %w", err)
		}
		if parent.WorkspaceID != workspaceID {
			return nil, fmt.Errorf("parent element belongs to different workspace")
		}
	}

	if err := s.canvasRepo.CreateElement(ctx, element); err != nil {
		return nil, fmt.Errorf("failed to create element: %w", err)
	}

	// Invalidate workspace cache
	if s.cacheService != nil {
		_ = s.cacheService.InvalidateWorkspaceElements(ctx, workspaceID)
	}

	return element, nil
}

// GetElement retrieves a canvas element by ID
func (s *CanvasService) GetElement(ctx context.Context, id uuid.UUID) (*models.CanvasElement, error) {
	element, err := s.canvasRepo.GetElementByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get element: %w", err)
	}

	return element, nil
}

// GetWorkspaceElements retrieves all elements for a workspace
func (s *CanvasService) GetWorkspaceElements(ctx context.Context, workspaceID uuid.UUID) ([]models.CanvasElement, error) {
	// Try cache first
	if s.cacheService != nil {
		if cachedElements, found := s.cacheService.GetWorkspaceElements(ctx, workspaceID); found {
			return cachedElements, nil
		}
	}

	// Cache miss - fetch from database
	elements, err := s.canvasRepo.GetElementsByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace elements: %w", err)
	}

	// Store in cache for next time
	if s.cacheService != nil {
		_ = s.cacheService.SetWorkspaceElements(ctx, workspaceID, elements)
	}

	return elements, nil
}

// UpdateElement updates a canvas element
func (s *CanvasService) UpdateElement(
	ctx context.Context,
	id, userID uuid.UUID,
	req models.UpdateElementRequest,
) (*models.CanvasElement, error) {
	// Get existing element
	element, err := s.canvasRepo.GetElementByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("element not found: %w", err)
	}

	// Apply partial updates
	if req.ElementData != nil {
		element.ElementData = *req.ElementData
	}
	if req.ZIndex != nil {
		element.ZIndex = *req.ZIndex
	}
	if req.ParentID != nil {
		// Validate parent exists
		parent, err := s.canvasRepo.GetElementByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent element not found: %w", err)
		}
		if parent.WorkspaceID != element.WorkspaceID {
			return nil, fmt.Errorf("parent element belongs to different workspace")
		}
		element.ParentID = req.ParentID
	}

	element.UpdatedBy = &userID

	if err := s.canvasRepo.UpdateElement(ctx, element); err != nil {
		return nil, fmt.Errorf("failed to update element: %w", err)
	}

	// Invalidate caches
	if s.cacheService != nil {
		_ = s.cacheService.InvalidateWorkspaceElements(ctx, element.WorkspaceID)
		_ = s.cacheService.InvalidateElement(ctx, id)
	}

	return element, nil
}

// DeleteElement soft deletes a canvas element
func (s *CanvasService) DeleteElement(ctx context.Context, id uuid.UUID) error {
	// Check if element has children (for groups)
	children, err := s.canvasRepo.GetChildElements(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check child elements: %w", err)
	}

	// If element has children, delete them too (cascade)
	if len(children) > 0 {
		childIDs := make([]uuid.UUID, len(children))
		for i := range children {
			childIDs[i] = children[i].ID
		}
		if err := s.canvasRepo.BatchDeleteElements(ctx, childIDs); err != nil {
			return fmt.Errorf("failed to delete child elements: %w", err)
		}
	}

	if err := s.canvasRepo.DeleteElement(ctx, id); err != nil {
		return fmt.Errorf("failed to delete element: %w", err)
	}

	// Invalidate caches
	if s.cacheService != nil {
		element, _ := s.canvasRepo.GetElementByID(ctx, id)
		if element != nil {
			_ = s.cacheService.InvalidateWorkspaceElements(ctx, element.WorkspaceID)
		}
		_ = s.cacheService.InvalidateElement(ctx, id)
	}

	return nil
}

// Batch operations

const (
	maxBatchSize = 100
)

// BatchCreateElements creates multiple canvas elements
func (s *CanvasService) BatchCreateElements(
	ctx context.Context,
	workspaceID, userID uuid.UUID,
	req models.BatchCreateRequest,
) ([]models.CanvasElement, error) {
	if len(req.Elements) == 0 {
		return nil, fmt.Errorf("no elements to create")
	}

	if len(req.Elements) > maxBatchSize {
		return nil, fmt.Errorf("cannot create more than %d elements at once", maxBatchSize)
	}

	elements := make([]models.CanvasElement, len(req.Elements))
	for i, createReq := range req.Elements {
		// Validate element type
		if !createReq.ElementType.Valid() {
			return nil, fmt.Errorf("invalid element type at index %d: %s", i, createReq.ElementType)
		}

		// Validate element data
		if len(createReq.ElementData) == 0 {
			return nil, fmt.Errorf("element_data is required at index %d", i)
		}

		elements[i] = models.CanvasElement{
			ID:          uuid.New(),
			WorkspaceID: workspaceID,
			ElementType: createReq.ElementType,
			ElementData: createReq.ElementData,
			ZIndex:      createReq.ZIndex,
			ParentID:    createReq.ParentID,
			CreatedBy:   userID,
			UpdatedBy:   &userID,
		}
	}

	if err := s.canvasRepo.BatchCreateElements(ctx, elements); err != nil {
		return nil, fmt.Errorf("failed to batch create elements: %w", err)
	}

	// Invalidate workspace cache
	if s.cacheService != nil {
		_ = s.cacheService.InvalidateWorkspaceElements(ctx, workspaceID)
	}

	return elements, nil
}

// BatchUpdateElements updates multiple canvas elements
func (s *CanvasService) BatchUpdateElements(
	ctx context.Context,
	workspaceID, userID uuid.UUID,
	req models.BatchUpdateRequest,
) ([]models.CanvasElement, error) {
	if len(req.Updates) == 0 {
		return nil, fmt.Errorf("no elements to update")
	}

	if len(req.Updates) > maxBatchSize {
		return nil, fmt.Errorf("cannot update more than %d elements at once", maxBatchSize)
	}

	// Fetch existing elements
	elements := make([]models.CanvasElement, len(req.Updates))
	for i, update := range req.Updates {
		element, err := s.canvasRepo.GetElementByID(ctx, update.ID)
		if err != nil {
			return nil, fmt.Errorf("element %s not found: %w", update.ID, err)
		}

		// Verify workspace
		if element.WorkspaceID != workspaceID {
			return nil, fmt.Errorf("element %s does not belong to workspace %s", update.ID, workspaceID)
		}

		// Apply partial updates
		if update.ElementData != nil {
			element.ElementData = *update.ElementData
		}
		if update.ZIndex != nil {
			element.ZIndex = *update.ZIndex
		}
		if update.ParentID != nil {
			element.ParentID = update.ParentID
		}

		element.UpdatedBy = &userID
		elements[i] = *element
	}

	if err := s.canvasRepo.BatchUpdateElements(ctx, elements); err != nil {
		return nil, fmt.Errorf("failed to batch update elements: %w", err)
	}

	// Invalidate caches
	if s.cacheService != nil {
		_ = s.cacheService.InvalidateWorkspaceElements(ctx, workspaceID)
		elementIDs := make([]uuid.UUID, len(req.Updates))
		for i, update := range req.Updates {
			elementIDs[i] = update.ID
		}
		_ = s.cacheService.InvalidateMultipleElements(ctx, elementIDs)
	}

	return elements, nil
}

// BatchDeleteElements soft deletes multiple canvas elements
func (s *CanvasService) BatchDeleteElements(ctx context.Context, workspaceID uuid.UUID, req models.BatchDeleteRequest) error {
	if len(req.IDs) == 0 {
		return fmt.Errorf("no elements to delete")
	}

	if len(req.IDs) > maxBatchSize {
		return fmt.Errorf("cannot delete more than %d elements at once", maxBatchSize)
	}

	// Verify all elements belong to the workspace
	for _, id := range req.IDs {
		element, err := s.canvasRepo.GetElementByID(ctx, id)
		if err != nil {
			return fmt.Errorf("element %s not found: %w", id, err)
		}
		if element.WorkspaceID != workspaceID {
			return fmt.Errorf("element %s does not belong to workspace %s", id, workspaceID)
		}
	}

	// Delete elements and their children
	var allIDs []uuid.UUID
	for _, id := range req.IDs {
		allIDs = append(allIDs, id)

		// Get children
		children, err := s.canvasRepo.GetChildElements(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to get child elements for %s: %w", id, err)
		}

		for i := range children {
			allIDs = append(allIDs, children[i].ID)
		}
	}

	if err := s.canvasRepo.BatchDeleteElements(ctx, allIDs); err != nil {
		return fmt.Errorf("failed to batch delete elements: %w", err)
	}

	// Invalidate caches
	if s.cacheService != nil {
		_ = s.cacheService.InvalidateWorkspaceElements(ctx, workspaceID)
		_ = s.cacheService.InvalidateMultipleElements(ctx, allIDs)
	}

	return nil
}

// GetElementsByType retrieves elements of a specific type
func (s *CanvasService) GetElementsByType(
	ctx context.Context,
	workspaceID uuid.UUID,
	elementType models.ElementType,
) ([]models.CanvasElement, error) {
	if !elementType.Valid() {
		return nil, fmt.Errorf("invalid element type: %s", elementType)
	}

	elements, err := s.canvasRepo.GetElementsByType(ctx, workspaceID, elementType)
	if err != nil {
		return nil, fmt.Errorf("failed to get elements by type: %w", err)
	}

	return elements, nil
}

// GetElementCount returns the total number of elements in a workspace
func (s *CanvasService) GetElementCount(ctx context.Context, workspaceID uuid.UUID) (int, error) {
	count, err := s.canvasRepo.GetElementCount(ctx, workspaceID)
	if err != nil {
		return 0, fmt.Errorf("failed to get element count: %w", err)
	}

	return count, nil
}

// Helper functions

// ValidateElementData performs basic validation on element data
func (s *CanvasService) ValidateElementData(elementType models.ElementType, data models.ElementData) error {
	if len(data) == 0 {
		return fmt.Errorf("element_data cannot be empty")
	}

	return s.validateElementTypeSpecific(elementType, data)
}

func (s *CanvasService) validateElementTypeSpecific(elementType models.ElementType, data models.ElementData) error {
	switch elementType {
	case models.ElementTypeText:
		return s.validateTextElement(data)
	case models.ElementTypeImage:
		return s.validateImageElement(data)
	case models.ElementTypeConnector:
		return s.validateConnectorElement(data)
	case models.ElementTypeShape, models.ElementTypeDrawing, models.ElementTypeSticky, models.ElementTypeList, models.ElementTypeGroup:
		return nil
	default:
		return nil
	}
}

func (s *CanvasService) validateTextElement(data models.ElementData) error {
	if _, ok := data["content"]; !ok {
		return fmt.Errorf("text element must have 'content' field")
	}
	return nil
}

func (s *CanvasService) validateImageElement(data models.ElementData) error {
	if _, ok := data["url"]; !ok {
		return fmt.Errorf("image element must have 'url' field")
	}
	return nil
}

func (s *CanvasService) validateConnectorElement(data models.ElementData) error {
	if _, hasStart := data["start_element_id"]; !hasStart {
		if _, hasStartPoint := data["start_point"]; !hasStartPoint {
			return fmt.Errorf("connector must have either 'start_element_id' or 'start_point'")
		}
	}
	if _, hasEnd := data["end_element_id"]; !hasEnd {
		if _, hasEndPoint := data["end_point"]; !hasEndPoint {
			return fmt.Errorf("connector must have either 'end_element_id' or 'end_point'")
		}
	}
	return nil
}
