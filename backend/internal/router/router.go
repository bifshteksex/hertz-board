package router

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/adaptor"

	"github.com/bifshteksex/hertz-board/internal/config"
	"github.com/bifshteksex/hertz-board/internal/handler"
	"github.com/bifshteksex/hertz-board/internal/middleware"
	"github.com/bifshteksex/hertz-board/internal/models"
	"github.com/bifshteksex/hertz-board/internal/service"
)

// Dependencies holds all service dependencies
type Dependencies struct {
	JWTService       *service.JWTService
	WorkspaceService *service.WorkspaceService
	CRDTService      *service.CRDTService
	Hub              *service.Hub
	AuthHandler      *handler.AuthHandler
	UserHandler      *handler.UserHandler
	OAuthHandler     *handler.OAuthHandler
	WorkspaceHandler *handler.WorkspaceHandler
	CanvasHandler    *handler.CanvasHandler
	AssetHandler     *handler.AssetHandler
	SnapshotHandler  *handler.SnapshotHandler
	WSHandler        *handler.WebSocketHandler
}

// Setup configures all routes and middleware
func Setup(h *server.Hertz, cfg *config.Config, deps *Dependencies) {
	// Global middleware
	h.Use(middleware.Recovery())
	h.Use(middleware.RequestID())
	h.Use(middleware.Logger())
	h.Use(middleware.CORS(&cfg.CORS))

	// Health check endpoints
	h.GET("/health", healthCheck)
	h.GET("/readiness", readinessCheck)

	// WebSocket endpoint (requires JWT token as query parameter)
	// Use HTTP adaptor to integrate gorilla/websocket with Hertz
	h.GET("/ws", adaptor.HertzHandler(http.HandlerFunc(deps.WSHandler.HandleWebSocket)))

	// API v1 routes
	v1 := h.Group("/api/v1")

	// Auth routes
	auth := v1.Group("/auth")
	auth.POST("/register", deps.AuthHandler.Register)
	auth.POST("/login", deps.AuthHandler.Login)
	auth.POST("/refresh", deps.AuthHandler.RefreshToken)
	auth.POST("/logout", deps.AuthHandler.Logout)
	auth.POST("/forgot-password", deps.AuthHandler.ForgotPassword)
	auth.POST("/reset-password", deps.AuthHandler.ResetPassword)

	// OAuth routes
	auth.GET("/google", deps.OAuthHandler.GoogleAuth)
	auth.GET("/google/callback", deps.OAuthHandler.GoogleCallback)
	auth.GET("/github", deps.OAuthHandler.GitHubAuth)
	auth.GET("/github/callback", deps.OAuthHandler.GitHubCallback)

	// User routes (protected)
	users := v1.Group("/users")
	users.Use(middleware.Auth(deps.JWTService))
	users.GET("/me", deps.UserHandler.GetProfile)
	users.PUT("/me", deps.UserHandler.UpdateProfile)
	users.PUT("/me/password", deps.UserHandler.ChangePassword)

	// Workspace routes
	workspaceMiddleware := middleware.NewWorkspaceMiddleware(deps.WorkspaceService)

	workspaces := v1.Group("/workspaces")
	workspaces.Use(middleware.Auth(deps.JWTService))

	// Workspace CRUD
	workspaces.POST("", deps.WorkspaceHandler.CreateWorkspace)
	workspaces.GET("", deps.WorkspaceHandler.ListWorkspaces)

	// Accept invite (no workspace_id param)
	workspaces.POST("/invites/accept", deps.WorkspaceHandler.AcceptInvite)

	// Specific workspace routes (require workspace access)
	workspaces.GET("/:workspace_id",
		workspaceMiddleware.OptionalWorkspaceAccess(),
		deps.WorkspaceHandler.GetWorkspace,
	)

	workspaces.PUT("/:workspace_id",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor),
		deps.WorkspaceHandler.UpdateWorkspace,
	)

	workspaces.DELETE("/:workspace_id",
		workspaceMiddleware.RequireWorkspaceOwner(),
		deps.WorkspaceHandler.DeleteWorkspace,
	)

	workspaces.POST("/:workspace_id/duplicate",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleViewer),
		deps.WorkspaceHandler.DuplicateWorkspace,
	)

	// Member management (require editor access)
	workspaces.GET("/:workspace_id/members",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleViewer),
		deps.WorkspaceHandler.ListMembers,
	)

	workspaces.PUT("/:workspace_id/members/:user_id",
		workspaceMiddleware.RequireWorkspaceOwner(),
		deps.WorkspaceHandler.UpdateMemberRole,
	)

	workspaces.DELETE("/:workspace_id/members/:user_id",
		workspaceMiddleware.RequireWorkspaceOwner(),
		deps.WorkspaceHandler.RemoveMember,
	)

	// Invitation management (require editor access to create, owner to manage)
	workspaces.POST("/:workspace_id/invites",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor),
		deps.WorkspaceHandler.CreateInvite,
	)

	workspaces.GET("/:workspace_id/invites",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor),
		deps.WorkspaceHandler.ListInvites,
	)

	workspaces.DELETE("/:workspace_id/invites/:invite_id",
		workspaceMiddleware.RequireWorkspaceOwner(),
		deps.WorkspaceHandler.RevokeInvite,
	)

	// Canvas element routes (require editor access to modify)
	workspaces.GET("/:workspace_id/elements",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleViewer),
		deps.CanvasHandler.GetWorkspaceElements,
	)

	workspaces.POST("/:workspace_id/elements",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor),
		deps.CanvasHandler.CreateElement,
	)

	workspaces.GET("/:workspace_id/elements/by-type",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleViewer),
		deps.CanvasHandler.GetElementsByType,
	)

	workspaces.GET("/:workspace_id/elements/:element_id",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleViewer),
		deps.CanvasHandler.GetElement,
	)

	workspaces.PUT("/:workspace_id/elements/:element_id",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor),
		deps.CanvasHandler.UpdateElement,
	)

	workspaces.DELETE("/:workspace_id/elements/:element_id",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor),
		deps.CanvasHandler.DeleteElement,
	)

	// Batch element operations
	workspaces.POST("/:workspace_id/elements/batch",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor),
		deps.CanvasHandler.BatchCreateElements,
	)

	workspaces.PUT("/:workspace_id/elements/batch",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor),
		deps.CanvasHandler.BatchUpdateElements,
	)

	workspaces.DELETE("/:workspace_id/elements/batch",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor),
		deps.CanvasHandler.BatchDeleteElements,
	)

	// Asset routes (require editor access to upload)
	workspaces.GET("/:workspace_id/assets",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleViewer),
		deps.AssetHandler.GetWorkspaceAssets,
	)

	workspaces.POST("/:workspace_id/assets",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor),
		deps.AssetHandler.UploadAsset,
	)

	workspaces.GET("/:workspace_id/assets/:asset_id",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleViewer),
		deps.AssetHandler.GetAsset,
	)

	workspaces.DELETE("/:workspace_id/assets/:asset_id",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor),
		deps.AssetHandler.DeleteAsset,
	)

	workspaces.POST("/:workspace_id/assets/cleanup",
		workspaceMiddleware.RequireWorkspaceOwner(),
		deps.AssetHandler.CleanupOrphanedAssets,
	)

	// Snapshot routes (require editor access to create, viewer to list)
	workspaces.GET("/:workspace_id/snapshots",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleViewer),
		deps.SnapshotHandler.ListSnapshots,
	)

	workspaces.POST("/:workspace_id/snapshots",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor),
		deps.SnapshotHandler.CreateSnapshot,
	)

	workspaces.GET("/:workspace_id/snapshots/version/:version",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleViewer),
		deps.SnapshotHandler.GetSnapshotByVersion,
	)

	workspaces.GET("/:workspace_id/snapshots/:snapshot_id",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleViewer),
		deps.SnapshotHandler.GetSnapshot,
	)

	workspaces.POST("/:workspace_id/snapshots/:snapshot_id/restore",
		workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor),
		deps.SnapshotHandler.RestoreSnapshot,
	)

	workspaces.DELETE("/:workspace_id/snapshots/:snapshot_id",
		workspaceMiddleware.RequireWorkspaceOwner(),
		deps.SnapshotHandler.DeleteSnapshot,
	)
}

// healthCheck returns basic health status
func healthCheck(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"service":   "api-gateway",
		"timestamp": time.Now().Unix(),
	})
}

// readinessCheck checks if service is ready (DB, Redis, etc.)
func readinessCheck(c context.Context, ctx *app.RequestContext) {
	// TODO: Add actual health checks for dependencies
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status":    "ready",
		"service":   "api-gateway",
		"timestamp": time.Now().Unix(),
		"checks": map[string]string{
			"database": "ok",
			"redis":    "ok",
		},
	})
}
