package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/uuid"

	"github.com/bifshteksex/hertz-board/internal/models"
	"github.com/bifshteksex/hertz-board/internal/service"
)

var ErrInvalidRequestType = errors.New("invalid request type")

type CanvasHandler struct {
	canvasService *service.CanvasService
}

func NewCanvasHandler(canvasService *service.CanvasService) *CanvasHandler {
	return &CanvasHandler{
		canvasService: canvasService,
	}
}

// Helper function for single element operations
func (h *CanvasHandler) processElementRequest(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
	reqPtr interface{},
	operation func(context.Context, uuid.UUID, uuid.UUID, interface{}) (*models.CanvasElement, error),
) (interface{}, error) {
	element, err := operation(ctx, id, userID, reqPtr)
	if err != nil {
		return nil, err
	}
	return element.ToResponse(), nil
}

// Helper function for batch element operations
func (h *CanvasHandler) processBatchElementRequest(
	ctx context.Context,
	workspaceID uuid.UUID,
	userID uuid.UUID,
	reqPtr interface{},
	operation func(context.Context, uuid.UUID, uuid.UUID, interface{}) ([]models.CanvasElement, error),
) ([]interface{}, error) {
	elements, err := operation(ctx, workspaceID, userID, reqPtr)
	if err != nil {
		return nil, err
	}
	results := make([]interface{}, len(elements))
	for i := range elements {
		results[i] = elements[i].ToResponse()
	}
	return results, nil
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
//
//nolint:dupl,errcheck // Similar pattern needed for create/update operations
func (h *CanvasHandler) CreateElement(ctx context.Context, c *app.RequestContext) {
	var req models.CreateElementRequest
	handleElementOperation(
		ctx, c, "", &req,
		func(ctx context.Context, id uuid.UUID, userID uuid.UUID, reqPtr interface{}) (interface{}, error) {
			createReq, ok := reqPtr.(*models.CreateElementRequest)
			if !ok {
				return nil, ErrInvalidRequestType
			}
			return h.processElementRequest(ctx, id, userID, createReq,
				func(ctx context.Context, id, userID uuid.UUID, r interface{}) (*models.CanvasElement, error) {
					return h.canvasService.CreateElement(ctx, id, userID, *r.(*models.CreateElementRequest))
				})
		},
		"Failed to create element",
		http.StatusCreated,
	)
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
	handleGetByID(ctx, c, "element_id", func(ctx context.Context, id uuid.UUID) (interface{}, error) {
		element, err := h.canvasService.GetElement(ctx, id)
		if err != nil {
			return nil, err
		}
		return element.ToResponse(), nil
	}, "Failed to get element")
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
//
//nolint:dupl,errcheck // Similar pattern needed for create/update operations
func (h *CanvasHandler) UpdateElement(ctx context.Context, c *app.RequestContext) {
	var req models.UpdateElementRequest
	handleElementOperation(
		ctx, c, "element_id", &req,
		func(ctx context.Context, id uuid.UUID, userID uuid.UUID, reqPtr interface{}) (interface{}, error) {
			updateReq, ok := reqPtr.(*models.UpdateElementRequest)
			if !ok {
				return nil, ErrInvalidRequestType
			}
			return h.processElementRequest(ctx, id, userID, updateReq,
				func(ctx context.Context, id, userID uuid.UUID, r interface{}) (*models.CanvasElement, error) {
					return h.canvasService.UpdateElement(ctx, id, userID, *r.(*models.UpdateElementRequest))
				})
		},
		"Failed to update element",
		http.StatusOK,
	)
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
	handleDeleteByID(ctx, c, "element_id", h.canvasService.DeleteElement, "Failed to delete element", "Element deleted successfully")
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
//
//nolint:dupl,errcheck // Similar pattern needed for batch create/update operations
func (h *CanvasHandler) BatchCreateElements(ctx context.Context, c *app.RequestContext) {
	var req models.BatchCreateRequest
	handleBatchElementOperation(
		ctx, c, &req,
		func(
			ctx context.Context,
			workspaceID uuid.UUID,
			userID uuid.UUID,
			reqPtr interface{},
		) ([]interface{}, error) {
			batchReq, ok := reqPtr.(*models.BatchCreateRequest)
			if !ok {
				return nil, ErrInvalidRequestType
			}
			return h.processBatchElementRequest(ctx, workspaceID, userID, batchReq,
				func(ctx context.Context, wID, uID uuid.UUID, r interface{}) ([]models.CanvasElement, error) {
					return h.canvasService.BatchCreateElements(ctx, wID, uID, *r.(*models.BatchCreateRequest))
				})
		},
		"Failed to batch create elements",
		http.StatusCreated,
	)
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
//
//nolint:dupl,errcheck // Similar pattern needed for batch create/update operations
func (h *CanvasHandler) BatchUpdateElements(ctx context.Context, c *app.RequestContext) {
	var req models.BatchUpdateRequest
	handleBatchElementOperation(
		ctx, c, &req,
		func(
			ctx context.Context,
			workspaceID uuid.UUID,
			userID uuid.UUID,
			reqPtr interface{},
		) ([]interface{}, error) {
			batchReq, ok := reqPtr.(*models.BatchUpdateRequest)
			if !ok {
				return nil, ErrInvalidRequestType
			}
			return h.processBatchElementRequest(ctx, workspaceID, userID, batchReq,
				func(ctx context.Context, wID, uID uuid.UUID, r interface{}) ([]models.CanvasElement, error) {
					return h.canvasService.BatchUpdateElements(ctx, wID, uID, *r.(*models.BatchUpdateRequest))
				})
		},
		"Failed to batch update elements",
		http.StatusOK,
	)
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
