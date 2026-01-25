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

type AssetHandler struct {
	assetService *service.AssetService
}

func NewAssetHandler(assetService *service.AssetService) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
	}
}

// UploadAsset godoc
// @Summary Upload an asset file
// @Description Uploads an image or file to the workspace
// @Tags assets
// @Accept multipart/form-data
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param file formData file true "File to upload"
// @Success 201 {object} models.AssetResponse
//
// @Router /api/v1/workspaces/{workspace_id}/assets [post]
func (h *AssetHandler) UploadAsset(ctx context.Context, c *app.RequestContext) {
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

	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "No file uploaded"})
		return
	}

	// Validate content type
	contentType := fileHeader.Header.Get("Content-Type")
	if !h.assetService.ValidateContentType(contentType) {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Unsupported file type. Only images are allowed."})
		return
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to open uploaded file: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to process file"})
		return
	}
	defer file.Close()

	// Upload asset
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Invalid user ID format"})
		return
	}

	asset, err := h.assetService.UploadAsset(
		ctx,
		workspaceID,
		userUUID,
		fileHeader.Filename,
		contentType,
		fileHeader.Size,
		file,
	)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to upload asset: %v", err)
		if fileHeader.Size > 10*1024*1024 {
			c.JSON(http.StatusRequestEntityTooLarge, map[string]interface{}{"error": "File too large. Maximum size is 10MB."})
		} else {
			c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, asset.ToResponse())
}

// GetAsset godoc
// @Summary Get an asset by ID
// @Description Retrieves asset metadata
// @Tags assets
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param asset_id path string true "Asset ID"
// @Success 200 {object} models.AssetResponse
//
// @Router /api/v1/workspaces/{workspace_id}/assets/{asset_id} [get]
func (h *AssetHandler) GetAsset(ctx context.Context, c *app.RequestContext) {
	assetIDStr := c.Param("asset_id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid asset ID"})
		return
	}

	asset, err := h.assetService.GetAsset(ctx, assetID)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get asset: %v", err)
		c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Asset not found"})
		return
	}

	c.JSON(http.StatusOK, asset.ToResponse())
}

// GetWorkspaceAssets godoc
// @Summary Get all assets in a workspace
// @Description Retrieves all assets for a workspace
// @Tags assets
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {object} map[string][]models.AssetResponse
//
// @Router /api/v1/workspaces/{workspace_id}/assets [get]
func (h *AssetHandler) GetWorkspaceAssets(ctx context.Context, c *app.RequestContext) {
	workspaceIDStr := c.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid workspace ID"})
		return
	}

	assets, err := h.assetService.GetWorkspaceAssets(ctx, workspaceID)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to get workspace assets: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to get assets"})
		return
	}

	// Convert to response
	responses := make([]models.AssetResponse, len(assets))
	for i := range assets {
		responses[i] = assets[i].ToResponse()
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"assets": responses,
		"total":  len(responses),
	})
}

// DeleteAsset godoc
// @Summary Delete an asset
// @Description Soft deletes an asset
// @Tags assets
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Param asset_id path string true "Asset ID"
//
// @Router /api/v1/workspaces/{workspace_id}/assets/{asset_id} [delete]
func (h *AssetHandler) DeleteAsset(ctx context.Context, c *app.RequestContext) {
	assetIDStr := c.Param("asset_id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid asset ID"})
		return
	}

	if err := h.assetService.DeleteAsset(ctx, assetID); err != nil {
		hlog.CtxErrorf(ctx, "Failed to delete asset: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "Asset deleted successfully"})
}

// CleanupOrphanedAssets godoc
// @Summary Cleanup orphaned assets
// @Description Deletes assets not referenced by any canvas element
// @Tags assets
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {object} map[string]interface{}
//
// @Router /api/v1/workspaces/{workspace_id}/assets/cleanup [post]
func (h *AssetHandler) CleanupOrphanedAssets(ctx context.Context, c *app.RequestContext) {
	workspaceIDStr := c.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid workspace ID"})
		return
	}

	count, err := h.assetService.CleanupOrphanedAssets(ctx, workspaceID)
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to cleanup orphaned assets: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to cleanup assets"})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Orphaned assets cleaned up successfully",
		"count":   count,
	})
}
