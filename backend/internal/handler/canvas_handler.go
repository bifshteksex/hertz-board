package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/uuid"

	"github.com/bifshteksex/hertz-board/internal/models"
	"github.com/bifshteksex/hertz-board/internal/service"
)

type CanvasHandler struct {
	canvasService *service.CanvasService
}

func NewCanvasHandler(canvasService *service.CanvasService) *CanvasHandler {
	return &CanvasHandler{
		canvasService: canvasService,
	}
}

// GetWorkspaceElements godoc
// @Summary Get all elements in a workspace
// @Description Retrieves all canvas elements for a workspace
// @Tags canvas
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {object} models.ElementListResponse
//
// @Router /api/v1/workspaces/{workspace_id}/elements [get]
func (h *CanvasHandler) GetWorkspaceElements(ctx context.Context, c *app.RequestContext) {
	workspaceIDStr := c.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid workspace ID"})
		return
	}

	elements, err := h.canvasService.GetWorkspaceElements(ctx, workspaceID)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get workspace elements: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to get elements"})
		return
	}

	// Convert to response
	responses := make([]models.ElementResponse, len(elements))
	for i := range elements {
		responses[i] = elements[i].ToResponse()
	}

	c.JSON(http.StatusOK, models.ElementListResponse{
		Elements: responses,
		Total:    len(responses),
	})
}

// CreateElement godoc
// @Summary Create a new canvas element
// @Description Creates a new canvas element in a workspace
// @Tags canvas
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param request body models.CreateElementRequest true "Element data"
// @Success 201 {object} models.ElementResponse
//
// @Router /api/v1/workspaces/{workspace_id}/elements [post]
func (h *CanvasHandler) CreateElement(ctx context.Context, c *app.RequestContext) {
	workspaceIDStr := c.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid workspace ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "User not authenticated"})
		return
	}

	var req models.CreateElementRequest
	if bindErr := c.BindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request body"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Invalid user ID format"})
		return
	}

	element, err := h.canvasService.CreateElement(ctx, workspaceID, userUUID, req)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to create element: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, element.ToResponse())
}

// GetElement godoc
// @Summary Get a canvas element by ID
// @Description Retrieves a specific canvas element
// @Tags canvas
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param element_id path string true "Element ID"
// @Success 200 {object} models.ElementResponse
//
// @Router /api/v1/workspaces/{workspace_id}/elements/{element_id} [get]
func (h *CanvasHandler) GetElement(ctx context.Context, c *app.RequestContext) {
	elementIDStr := c.Param("element_id")
	elementID, err := uuid.Parse(elementIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid element ID"})
		return
	}

	element, err := h.canvasService.GetElement(ctx, elementID)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get element: %v", err)
		c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Element not found"})
		return
	}

	c.JSON(http.StatusOK, element.ToResponse())
}

// UpdateElement godoc
// @Summary Update a canvas element
// @Description Updates an existing canvas element
// @Tags canvas
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param element_id path string true "Element ID"
// @Param request body models.UpdateElementRequest true "Element data"
// @Success 200 {object} models.ElementResponse
//
// @Router /api/v1/workspaces/{workspace_id}/elements/{element_id} [put]
func (h *CanvasHandler) UpdateElement(ctx context.Context, c *app.RequestContext) {
	elementIDStr := c.Param("element_id")
	elementID, err := uuid.Parse(elementIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid element ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "User not authenticated"})
		return
	}

	var req models.UpdateElementRequest
	if bindErr := c.BindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request body"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Invalid user ID format"})
		return
	}

	element, err := h.canvasService.UpdateElement(ctx, elementID, userUUID, req)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to update element: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, element.ToResponse())
}

// DeleteElement godoc
// @Summary Delete a canvas element
// @Description Soft deletes a canvas element
// @Tags canvas
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param element_id path string true "Element ID"
//
// @Router /api/v1/workspaces/{workspace_id}/elements/{element_id} [delete]
func (h *CanvasHandler) DeleteElement(ctx context.Context, c *app.RequestContext) {
	elementIDStr := c.Param("element_id")
	elementID, err := uuid.Parse(elementIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid element ID"})
		return
	}

	if err := h.canvasService.DeleteElement(ctx, elementID); err != nil {
		hlog.CtxErrorf(ctx, "Failed to delete element: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "Element deleted successfully"})
}

// Batch operations

// BatchCreateElements godoc
// @Summary Create multiple canvas elements
// @Description Creates multiple canvas elements in a single request
// @Tags canvas
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param request body models.BatchCreateRequest true "Elements data"
// @Success 201 {object} models.ElementListResponse
//
// @Router /api/v1/workspaces/{workspace_id}/elements/batch [post]
func (h *CanvasHandler) BatchCreateElements(ctx context.Context, c *app.RequestContext) {
	workspaceIDStr := c.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid workspace ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "User not authenticated"})
		return
	}

	var req models.BatchCreateRequest
	if bindErr := c.BindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request body"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Invalid user ID format"})
		return
	}

	elements, err := h.canvasService.BatchCreateElements(ctx, workspaceID, userUUID, req)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to batch create elements: %v", err)
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	// Convert to response
	responses := make([]models.ElementResponse, len(elements))
	for i := range elements {
		responses[i] = elements[i].ToResponse()
	}

	c.JSON(http.StatusCreated, models.ElementListResponse{
		Elements: responses,
		Total:    len(responses),
	})
}

// BatchUpdateElements godoc
// @Summary Update multiple canvas elements
// @Description Updates multiple canvas elements in a single request
// @Tags canvas
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param request body models.BatchUpdateRequest true "Elements data"
// @Success 200 {object} models.ElementListResponse
//
// @Router /api/v1/workspaces/{workspace_id}/elements/batch [put]
func (h *CanvasHandler) BatchUpdateElements(ctx context.Context, c *app.RequestContext) {
	workspaceIDStr := c.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid workspace ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "User not authenticated"})
		return
	}

	var req models.BatchUpdateRequest
	if bindErr := c.BindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request body"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Invalid user ID format"})
		return
	}

	elements, err := h.canvasService.BatchUpdateElements(ctx, workspaceID, userUUID, req)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to batch update elements: %v", err)
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	// Convert to response
	responses := make([]models.ElementResponse, len(elements))
	for i := range elements {
		responses[i] = elements[i].ToResponse()
	}

	c.JSON(http.StatusOK, models.ElementListResponse{
		Elements: responses,
		Total:    len(responses),
	})
}

// BatchDeleteElements godoc
// @Summary Delete multiple canvas elements
// @Description Deletes multiple canvas elements in a single request
// @Tags canvas
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param request body models.BatchDeleteRequest true "Element IDs"
//
// @Router /api/v1/workspaces/{workspace_id}/elements/batch [delete]
func (h *CanvasHandler) BatchDeleteElements(ctx context.Context, c *app.RequestContext) {
	workspaceIDStr := c.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid workspace ID"})
		return
	}

	var req models.BatchDeleteRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request body"})
		return
	}

	if err := h.canvasService.BatchDeleteElements(ctx, workspaceID, req); err != nil {
		hlog.CtxErrorf(ctx, "Failed to batch delete elements: %v", err)
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "Elements deleted successfully"})
}

// GetElementsByType godoc
// @Summary Get elements by type
// @Description Retrieves all elements of a specific type in a workspace
// @Tags canvas
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param type query string true "Element type"
// @Success 200 {object} models.ElementListResponse
//
// @Router /api/v1/workspaces/{workspace_id}/elements/by-type [get]
func (h *CanvasHandler) GetElementsByType(ctx context.Context, c *app.RequestContext) {
	workspaceIDStr := c.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid workspace ID"})
		return
	}

	elementTypeStr := c.Query("type")
	if elementTypeStr == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Element type is required"})
		return
	}

	elementType := models.ElementType(elementTypeStr)
	if !elementType.Valid() {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid element type"})
		return
	}

	elements, err := h.canvasService.GetElementsByType(ctx, workspaceID, elementType)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get elements by type: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to get elements"})
		return
	}

	// Convert to response
	responses := make([]models.ElementResponse, len(elements))
	for i := range elements {
		responses[i] = elements[i].ToResponse()
	}

	c.JSON(http.StatusOK, models.ElementListResponse{
		Elements: responses,
		Total:    len(responses),
	})
}
