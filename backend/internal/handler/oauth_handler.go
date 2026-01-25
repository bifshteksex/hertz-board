package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/bifshteksex/hertzboard/internal/models"
	"github.com/bifshteksex/hertzboard/internal/service"
)

const (
	stateExpiration = 10 * time.Minute
	stateTokenBytes = 16
)

// OAuthHandler handles OAuth endpoints
type OAuthHandler struct {
	oauthService *service.OAuthService
	states       map[string]time.Time // Simple in-memory state storage (use Redis in production)
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(oauthService *service.OAuthService) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		states:       make(map[string]time.Time),
	}
}

// GoogleAuth redirects to Google OAuth
func (h *OAuthHandler) GoogleAuth(c context.Context, ctx *app.RequestContext) {
	state := h.generateState()
	h.states[state] = time.Now().Add(stateExpiration)

	url := h.oauthService.GetGoogleAuthURL(state)
	ctx.Redirect(consts.StatusTemporaryRedirect, []byte(url))
}

// GoogleCallback handles Google OAuth callback
func (h *OAuthHandler) GoogleCallback(c context.Context, ctx *app.RequestContext) {
	h.handleOAuthCallback(c, ctx, h.oauthService.GoogleCallback)
}

// GitHubAuth redirects to GitHub OAuth
func (h *OAuthHandler) GitHubAuth(c context.Context, ctx *app.RequestContext) {
	state := h.generateState()
	h.states[state] = time.Now().Add(stateExpiration)

	url := h.oauthService.GetGitHubAuthURL(state)
	ctx.Redirect(consts.StatusTemporaryRedirect, []byte(url))
}

// GitHubCallback handles GitHub OAuth callback
func (h *OAuthHandler) GitHubCallback(c context.Context, ctx *app.RequestContext) {
	h.handleOAuthCallback(c, ctx, h.oauthService.GitHubCallback)
}

// handleOAuthCallback is a common handler for OAuth callbacks
func (h *OAuthHandler) handleOAuthCallback(
	c context.Context,
	ctx *app.RequestContext,
	callbackFunc func(context.Context, string) (*models.AuthResponse, error),
) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	// Validate state
	if !h.validateState(state) {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": "Invalid state parameter",
		})
		return
	}

	// Handle OAuth callback
	resp, err := callbackFunc(c, code)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, resp)
}

// generateState generates a random state for OAuth
func (h *OAuthHandler) generateState() string {
	b := make([]byte, stateTokenBytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// validateState validates OAuth state parameter
func (h *OAuthHandler) validateState(state string) bool {
	// Clean up expired states
	now := time.Now()
	for s, expiry := range h.states {
		if now.After(expiry) {
			delete(h.states, s)
		}
	}

	// Check if state exists and is not expired
	expiry, exists := h.states[state]
	if !exists || now.After(expiry) {
		return false
	}

	// Delete used state
	delete(h.states, state)
	return true
}
