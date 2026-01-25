package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/uuid"

	"github.com/bifshteksex/hertz-board/internal/models"
	"github.com/bifshteksex/hertz-board/internal/service"
)

type SnapshotHandler struct {
	snapshotService *service.SnapshotService
}

func NewSnapshotHandler(snapshotService *service.SnapshotService) *SnapshotHandler {
	return &SnapshotHandler{
		snapshotService: snapshotService,
	}
}

// CreateSnapshot godoc
// @Summary Create a canvas snapshot
// @Description Creates a new version snapshot of the current canvas state
// @Tags snapshots
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param request body models.CreateSnapshotRequest false "Snapshot description"
// @Success 201 {object} models.SnapshotResponse
//
// @Router /api/v1/workspaces/{workspace_id}/snapshots [post]
func (h *SnapshotHandler) CreateSnapshot(ctx context.Context, c *app.RequestContext) {
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

	var req models.CreateSnapshotRequest
	if bindErr := c.BindJSON(&req); bindErr != nil {
		// Description is optional, so it's OK if body is empty
		req.Description = nil
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Invalid user ID format"})
		return
	}

	snapshot, err := h.snapshotService.CreateSnapshot(ctx, workspaceID, userUUID, req.Description)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to create snapshot: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, snapshot.ToResponse())
}

// ListSnapshots godoc
// @Summary List canvas snapshots
// @Description Retrieves all snapshots for a workspace with pagination
// @Tags snapshots
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param limit query int false "Number of results" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} models.SnapshotListResponse
//
// @Router /api/v1/workspaces/{workspace_id}/snapshots [get]
func (h *SnapshotHandler) ListSnapshots(ctx context.Context, c *app.RequestContext) {
	workspaceIDStr := c.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid workspace ID"})
		return
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	if limit <= 0 {
		limit = 20
	}

	snapshots, total, err := h.snapshotService.ListSnapshots(ctx, workspaceID, limit, offset)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to list snapshots: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to get snapshots"})
		return
	}

	// Convert to response
	responses := make([]models.SnapshotResponse, len(snapshots))
	for i, snapshot := range snapshots {
		responses[i] = snapshot.ToResponse()
	}

	c.JSON(http.StatusOK, models.SnapshotListResponse{
		Snapshots: responses,
		Total:     total,
	})
}

// GetSnapshot godoc
// @Summary Get a snapshot by ID
// @Description Retrieves a specific snapshot with full data
// @Tags snapshots
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param snapshot_id path string true "Snapshot ID"
// @Success 200 {object} models.SnapshotDetailResponse
//
// @Router /api/v1/workspaces/{workspace_id}/snapshots/{snapshot_id} [get]
func (h *SnapshotHandler) GetSnapshot(ctx context.Context, c *app.RequestContext) {
	handleGetByID(ctx, c, "snapshot_id", func(ctx context.Context, id uuid.UUID) (interface{}, error) {
		snapshot, err := h.snapshotService.GetSnapshot(ctx, id)
		if err != nil {
			return nil, err
		}
		return snapshot.ToDetailResponse(), nil
	}, "Failed to get snapshot")
}

// GetSnapshotByVersion godoc
// @Summary Get a snapshot by version number
// @Description Retrieves a specific snapshot by its version number
// @Tags snapshots
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param version path int true "Version number"
// @Success 200 {object} models.SnapshotDetailResponse
//
// @Router /api/v1/workspaces/{workspace_id}/snapshots/version/{version} [get]
func (h *SnapshotHandler) GetSnapshotByVersion(ctx context.Context, c *app.RequestContext) {
	workspaceIDStr := c.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid workspace ID"})
		return
	}

	versionStr := c.Param("version")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid version number"})
		return
	}

	snapshot, err := h.snapshotService.GetSnapshotByVersion(ctx, workspaceID, version)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get snapshot by version: %v", err)
		c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Snapshot not found"})
		return
	}

	c.JSON(http.StatusOK, snapshot.ToDetailResponse())
}

// RestoreSnapshot godoc
// @Summary Restore canvas to a snapshot
// @Description Restores the canvas to a specific snapshot version
// @Tags snapshots
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param snapshot_id path string true "Snapshot ID"
//
// @Router /api/v1/workspaces/{workspace_id}/snapshots/{snapshot_id}/restore [post]
func (h *SnapshotHandler) RestoreSnapshot(ctx context.Context, c *app.RequestContext) {
	workspaceIDStr := c.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid workspace ID"})
		return
	}

	snapshotIDStr := c.Param("snapshot_id")
	snapshotID, err := uuid.Parse(snapshotIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid snapshot ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "User not authenticated"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Invalid user ID format"})
		return
	}

	if err := h.snapshotService.RestoreSnapshot(ctx, workspaceID, userUUID, snapshotID); err != nil {
		hlog.CtxErrorf(ctx, "Failed to restore snapshot: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "Snapshot restored successfully"})
}

// DeleteSnapshot godoc
// @Summary Delete a snapshot
// @Description Deletes a specific snapshot
// @Tags snapshots
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param snapshot_id path string true "Snapshot ID"
//
// @Router /api/v1/workspaces/{workspace_id}/snapshots/{snapshot_id} [delete]
func (h *SnapshotHandler) DeleteSnapshot(ctx context.Context, c *app.RequestContext) {
	workspaceIDStr := c.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid workspace ID"})
		return
	}

	snapshotIDStr := c.Param("snapshot_id")
	snapshotID, err := uuid.Parse(snapshotIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid snapshot ID"})
		return
	}

	if err := h.snapshotService.DeleteSnapshot(ctx, workspaceID, snapshotID); err != nil {
		hlog.CtxErrorf(ctx, "Failed to delete snapshot: %v", err)
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "Snapshot deleted successfully"})
}
