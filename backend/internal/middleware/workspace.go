package middleware

import (
	"context"
	"net/http"

	"github.com/bifshteksex/hertz-board/internal/models"
	"github.com/bifshteksex/hertz-board/internal/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
)

type WorkspaceMiddleware struct {
	workspaceService *service.WorkspaceService
}

func NewWorkspaceMiddleware(workspaceService *service.WorkspaceService) *WorkspaceMiddleware {
	return &WorkspaceMiddleware{
		workspaceService: workspaceService,
	}
}

// RequireWorkspaceAccess checks if user has required access level to workspace
func (m *WorkspaceMiddleware) RequireWorkspaceAccess(requiredRole models.WorkspaceRole) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// Get user ID from context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		// Get workspace ID from path parameter
		workspaceIDStr := c.Param("workspace_id")
		if workspaceIDStr == "" {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Workspace ID is required",
			})
			c.Abort()
			return
		}

		workspaceID, err := uuid.Parse(workspaceIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Invalid workspace ID",
			})
			c.Abort()
			return
		}

		// Check permission
		uid, ok := userID.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid user ID",
			})
			c.Abort()
			return
		}
		if err := m.workspaceService.CheckPermission(ctx, workspaceID, uid, requiredRole); err != nil {
			c.JSON(http.StatusForbidden, map[string]interface{}{
				"error": "Access denied",
			})
			c.Abort()
			return
		}

		// Store workspace ID in context for handlers
		c.Set("workspace_id", workspaceID)
		c.Next(ctx)
	}
}

// RequireWorkspaceOwner checks if user is the owner of workspace
func (m *WorkspaceMiddleware) RequireWorkspaceOwner() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// Get user ID from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		// Get workspace ID from path parameter
		workspaceIDStr := c.Param("workspace_id")
		if workspaceIDStr == "" {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Workspace ID is required",
			})
			c.Abort()
			return
		}

		workspaceID, err := uuid.Parse(workspaceIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Invalid workspace ID",
			})
			c.Abort()
			return
		}

		// Check if owner
		uid, ok := userID.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid user ID",
			})
			c.Abort()
			return
		}
		isOwner, err := m.workspaceService.IsOwner(ctx, workspaceID, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to check ownership",
			})
			c.Abort()
			return
		}

		if !isOwner {
			c.JSON(http.StatusForbidden, map[string]interface{}{
				"error": "Only workspace owner can perform this action",
			})
			c.Abort()
			return
		}

		// Store workspace ID in context
		c.Set("workspace_id", workspaceID)
		c.Next(ctx)
	}
}

// OptionalWorkspaceAccess allows both authenticated and public access
func (m *WorkspaceMiddleware) OptionalWorkspaceAccess() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// Get workspace ID from path parameter
		workspaceIDStr := c.Param("workspace_id")
		if workspaceIDStr == "" {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Workspace ID is required",
			})
			c.Abort()
			return
		}

		workspaceID, err := uuid.Parse(workspaceIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Invalid workspace ID",
			})
			c.Abort()
			return
		}

		// Get workspace
		workspace, err := m.workspaceService.GetWorkspace(ctx, workspaceID)
		if err != nil {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Workspace not found",
			})
			c.Abort()
			return
		}

		// Check if user is authenticated
		userID, authenticated := c.Get("user_id")

		// If workspace is private and user is not authenticated, deny access
		if !workspace.IsPublic && !authenticated {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// If workspace is private and user is authenticated, check membership
		if !workspace.IsPublic && authenticated {
			uid, ok := userID.(uuid.UUID)
			if !ok {
				c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid user ID",
				})
				c.Abort()
				return
			}
			if err := m.workspaceService.CheckPermission(ctx, workspaceID, uid, models.WorkspaceRoleViewer); err != nil {
				c.JSON(http.StatusForbidden, map[string]interface{}{
					"error": "Access denied",
				})
				c.Abort()
				return
			}
		}

		// Store workspace ID in context
		c.Set("workspace_id", workspaceID)
		c.Next(ctx)
	}
}
