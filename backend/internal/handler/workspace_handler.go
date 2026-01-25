package handler

import (
	"context"
	"net/http"

	"github.com/bifshteksex/hertzboard/internal/models"
	"github.com/bifshteksex/hertzboard/internal/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
)

type WorkspaceHandler struct {
	workspaceService *service.WorkspaceService
}

func NewWorkspaceHandler(workspaceService *service.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{
		workspaceService: workspaceService,
	}
}

// CreateWorkspace creates a new workspace
// POST /api/v1/workspaces
func (h *WorkspaceHandler) CreateWorkspace(ctx context.Context, c *app.RequestContext) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.CreateWorkspaceRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	workspace, err := h.workspaceService.CreateWorkspace(ctx, &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"workspace": workspace,
	})
}

// ListWorkspaces lists all workspaces accessible to user
// GET /api/v1/workspaces
func (h *WorkspaceHandler) ListWorkspaces(ctx context.Context, c *app.RequestContext) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var filter models.WorkspaceListFilter
	if err := c.BindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid query parameters",
		})
		return
	}

	// Set defaults
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if filter.SortBy == "" {
		filter.SortBy = "updated_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	response, err := h.workspaceService.ListUserWorkspaces(ctx, userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetWorkspace retrieves a specific workspace
// GET /api/v1/workspaces/:workspace_id
func (h *WorkspaceHandler) GetWorkspace(ctx context.Context, c *app.RequestContext) {
	workspaceID := c.MustGet("workspace_id").(uuid.UUID)
	userID, authenticated := c.Get("user_id")

	if !authenticated {
		// Public access - just return workspace
		workspace, err := h.workspaceService.GetWorkspace(ctx, workspaceID)
		if err != nil {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Workspace not found",
			})
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"workspace": workspace,
		})
		return
	}

	// Authenticated - return with role
	uid := userID.(uuid.UUID)
	workspace, err := h.workspaceService.GetWorkspaceWithRole(ctx, workspaceID, uid)
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"workspace": workspace,
	})
}

// UpdateWorkspace updates workspace information
// PUT /api/v1/workspaces/:workspace_id
func (h *WorkspaceHandler) UpdateWorkspace(ctx context.Context, c *app.RequestContext) {
	workspaceID := c.MustGet("workspace_id").(uuid.UUID)

	var req models.UpdateWorkspaceRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	workspace, err := h.workspaceService.UpdateWorkspace(ctx, workspaceID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"workspace": workspace,
	})
}

// DeleteWorkspace deletes a workspace
// DELETE /api/v1/workspaces/:workspace_id
func (h *WorkspaceHandler) DeleteWorkspace(ctx context.Context, c *app.RequestContext) {
	workspaceID := c.MustGet("workspace_id").(uuid.UUID)

	if err := h.workspaceService.DeleteWorkspace(ctx, workspaceID); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Workspace deleted successfully",
	})
}

// DuplicateWorkspace creates a copy of a workspace
// POST /api/v1/workspaces/:workspace_id/duplicate
func (h *WorkspaceHandler) DuplicateWorkspace(ctx context.Context, c *app.RequestContext) {
	workspaceID := c.MustGet("workspace_id").(uuid.UUID)
	userID := c.MustGet("user_id").(uuid.UUID)

	workspace, err := h.workspaceService.DuplicateWorkspace(ctx, workspaceID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"workspace": workspace,
	})
}

// --- Member Management ---

// ListMembers retrieves all members of a workspace
// GET /api/v1/workspaces/:workspace_id/members
func (h *WorkspaceHandler) ListMembers(ctx context.Context, c *app.RequestContext) {
	workspaceID := c.MustGet("workspace_id").(uuid.UUID)

	members, err := h.workspaceService.GetMembers(ctx, workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"members": members,
	})
}

// UpdateMemberRole updates a member's role
// PUT /api/v1/workspaces/:workspace_id/members/:user_id
func (h *WorkspaceHandler) UpdateMemberRole(ctx context.Context, c *app.RequestContext) {
	workspaceID := c.MustGet("workspace_id").(uuid.UUID)

	memberUserIDStr := c.Param("user_id")
	memberUserID, err := uuid.Parse(memberUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid user ID",
		})
		return
	}

	var req models.UpdateMemberRoleRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.workspaceService.UpdateMemberRole(ctx, workspaceID, memberUserID, req.Role); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Member role updated successfully",
	})
}

// RemoveMember removes a member from workspace
// DELETE /api/v1/workspaces/:workspace_id/members/:user_id
func (h *WorkspaceHandler) RemoveMember(ctx context.Context, c *app.RequestContext) {
	workspaceID := c.MustGet("workspace_id").(uuid.UUID)

	memberUserIDStr := c.Param("user_id")
	memberUserID, err := uuid.Parse(memberUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid user ID",
		})
		return
	}

	if err := h.workspaceService.RemoveMember(ctx, workspaceID, memberUserID); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Member removed successfully",
	})
}

// --- Invitations ---

// CreateInvite creates a workspace invitation
// POST /api/v1/workspaces/:workspace_id/invites
func (h *WorkspaceHandler) CreateInvite(ctx context.Context, c *app.RequestContext) {
	workspaceID := c.MustGet("workspace_id").(uuid.UUID)
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.InviteToWorkspaceRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	tokenResponse, err := h.workspaceService.CreateInvite(ctx, workspaceID, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, tokenResponse)
}

// ListInvites retrieves all pending invitations for a workspace
// GET /api/v1/workspaces/:workspace_id/invites
func (h *WorkspaceHandler) ListInvites(ctx context.Context, c *app.RequestContext) {
	workspaceID := c.MustGet("workspace_id").(uuid.UUID)

	invites, err := h.workspaceService.GetPendingInvites(ctx, workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"invites": invites,
	})
}

// RevokeInvite revokes a pending invitation
// DELETE /api/v1/workspaces/:workspace_id/invites/:invite_id
func (h *WorkspaceHandler) RevokeInvite(ctx context.Context, c *app.RequestContext) {
	inviteIDStr := c.Param("invite_id")
	inviteID, err := uuid.Parse(inviteIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid invite ID",
		})
		return
	}

	if err := h.workspaceService.RevokeInvite(ctx, inviteID); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Invite revoked successfully",
	})
}

// AcceptInvite accepts a workspace invitation
// POST /api/v1/workspaces/invites/accept
func (h *WorkspaceHandler) AcceptInvite(ctx context.Context, c *app.RequestContext) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.AcceptInviteRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	workspace, err := h.workspaceService.AcceptInvite(ctx, req.Token, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"workspace": workspace,
		"message":   "Invitation accepted successfully",
	})
}
