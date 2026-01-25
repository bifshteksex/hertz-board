package router

import (
	"context"
	"net/http"
	"time"

	"github.com/bifshteksex/hertzboard/internal/config"
	"github.com/bifshteksex/hertzboard/internal/handler"
	"github.com/bifshteksex/hertzboard/internal/middleware"
	"github.com/bifshteksex/hertzboard/internal/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// Dependencies holds all service dependencies
type Dependencies struct {
	JWTService   *service.JWTService
	AuthHandler  *handler.AuthHandler
	UserHandler  *handler.UserHandler
	OAuthHandler *handler.OAuthHandler
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

	// Workspace routes (will be implemented in Phase 2)
	workspaces := v1.Group("/workspaces")
	workspaces.Use(middleware.Auth(deps.JWTService))
	workspaces.GET("", placeholderHandler("list-workspaces"))
	workspaces.POST("", placeholderHandler("create-workspace"))
	workspaces.GET("/:id", placeholderHandler("get-workspace"))
	workspaces.PUT("/:id", placeholderHandler("update-workspace"))
	workspaces.DELETE("/:id", placeholderHandler("delete-workspace"))
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

// placeholderHandler returns a temporary handler for routes not yet implemented
func placeholderHandler(name string) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(http.StatusNotImplemented, map[string]interface{}{
			"error":   "Not implemented yet",
			"handler": name,
		})
	}
}
