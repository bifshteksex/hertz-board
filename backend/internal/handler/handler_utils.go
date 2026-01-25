package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/uuid"
)

// parseIDParam parses a UUID from a request parameter
func parseIDParam(c *app.RequestContext, paramName string) (uuid.UUID, error) {
	idStr := c.Param(paramName)
	return uuid.Parse(idStr)
}

// handleGetByID is a generic handler for getting a resource by ID
func handleGetByID(
	ctx context.Context,
	c *app.RequestContext,
	paramName string,
	fetchFunc func(context.Context, uuid.UUID) (interface{}, error),
	errorMsg string,
) {
	id, err := parseIDParam(c, paramName)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid " + paramName})
		return
	}

	result, err := fetchFunc(ctx, id)
	if err != nil {
		hlog.CtxErrorf(ctx, "%s: %v", errorMsg, err)
		c.JSON(http.StatusNotFound, map[string]interface{}{"error": errorMsg})
		return
	}

	c.JSON(http.StatusOK, result)
}

// handleDeleteByID is a generic handler for deleting a resource by ID
func handleDeleteByID(
	ctx context.Context,
	c *app.RequestContext,
	paramName string,
	deleteFunc func(context.Context, uuid.UUID) error,
	errorMsg string,
	successMsg string,
) {
	id, err := parseIDParam(c, paramName)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid " + paramName})
		return
	}

	if err := deleteFunc(ctx, id); err != nil {
		hlog.CtxErrorf(ctx, "%s: %v", errorMsg, err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": successMsg})
}

// handleElementOperation is a generic handler for element creation/update operations
func handleElementOperation(
	ctx context.Context,
	c *app.RequestContext,
	idParam string,
	requestPtr interface{},
	operationFunc func(context.Context, uuid.UUID, uuid.UUID, interface{}) (interface{}, error),
	errorMsg string,
	statusCode int,
) {
	var id uuid.UUID
	var err error

	// If idParam is not empty, parse it (for update operation)
	if idParam != "" {
		idStr := c.Param(idParam)
		id, err = uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid " + idParam})
			return
		}
	} else {
		// For create operation, parse workspace_id
		workspaceIDStr := c.Param("workspace_id")
		id, err = uuid.Parse(workspaceIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid workspace ID"})
			return
		}
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "User not authenticated"})
		return
	}

	if bindErr := c.BindJSON(requestPtr); bindErr != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request body"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Invalid user ID format"})
		return
	}

	result, err := operationFunc(ctx, id, userUUID, requestPtr)
	if err != nil {
		hlog.CtxErrorf(ctx, "%s: %v", errorMsg, err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(statusCode, result)
}

// handleBatchElementOperation is a generic handler for batch element operations
func handleBatchElementOperation(
	ctx context.Context,
	c *app.RequestContext,
	requestPtr interface{},
	operationFunc func(context.Context, uuid.UUID, uuid.UUID, interface{}) ([]interface{}, error),
	errorMsg string,
	statusCode int,
) {
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

	if bindErr := c.BindJSON(requestPtr); bindErr != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request body"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Invalid user ID format"})
		return
	}

	results, err := operationFunc(ctx, workspaceID, userUUID, requestPtr)
	if err != nil {
		hlog.CtxErrorf(ctx, "%s: %v", errorMsg, err)
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(statusCode, map[string]interface{}{
		"elements": results,
		"total":    len(results),
	})
}
