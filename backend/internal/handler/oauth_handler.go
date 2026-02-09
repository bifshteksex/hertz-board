package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/url"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/bifshteksex/hertz-board/internal/models"
	"github.com/bifshteksex/hertz-board/internal/service"
)

const (
	stateExpiration = 10 * time.Minute
	stateTokenBytes = 16
)

// OAuthState stores state information for OAuth flow
type OAuthState struct {
	Expiry      time.Time
	RedirectURI string // Optional: for desktop app deep linking
}

// OAuthHandler handles OAuth endpoints
type OAuthHandler struct {
	oauthService *service.OAuthService
	states       map[string]*OAuthState // todo: Redis
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(oauthService *service.OAuthService) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		states:       make(map[string]*OAuthState),
	}
}

// GoogleAuth redirects to Google OAuth
func (h *OAuthHandler) GoogleAuth(c context.Context, ctx *app.RequestContext) {
	// Get optional redirect_uri for desktop app deep linking
	redirectURI := ctx.Query("redirect_uri")

	state := h.generateState()
	h.states[state] = &OAuthState{
		Expiry:      time.Now().Add(stateExpiration),
		RedirectURI: redirectURI,
	}

	url := h.oauthService.GetGoogleAuthURL(state)
	ctx.Redirect(consts.StatusTemporaryRedirect, []byte(url))
}

// GoogleCallback handles Google OAuth callback
func (h *OAuthHandler) GoogleCallback(c context.Context, ctx *app.RequestContext) {
	h.handleOAuthCallback(c, ctx, h.oauthService.GoogleCallback)
}

// GitHubAuth redirects to GitHub OAuth
func (h *OAuthHandler) GitHubAuth(c context.Context, ctx *app.RequestContext) {
	// Get optional redirect_uri for desktop app deep linking
	redirectURI := ctx.Query("redirect_uri")

	state := h.generateState()
	h.states[state] = &OAuthState{
		Expiry:      time.Now().Add(stateExpiration),
		RedirectURI: redirectURI,
	}

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

	// Validate state and get redirect URI
	redirectURI, valid := h.validateState(state)
	if !valid {
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

	const (
		refreshTokenMaxAge = 7 * 24 * 60 * 60 // 7 days in seconds
	)

	// Check if this is a desktop app callback (deep link)
	if redirectURI != "" && strings.HasPrefix(redirectURI, "hertzboard://") {
		// Desktop app - redirect with tokens in URL

		// Validate redirect URI for security
		if redirectURI != "hertzboard://oauth/callback" {
			ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
				"error": "Invalid redirect URI",
			})
			return
		}

		// Build redirect URL with tokens
		params := url.Values{}
		params.Add("access_token", resp.Tokens.AccessToken)
		params.Add("refresh_token", resp.Tokens.RefreshToken)

		redirectURL := redirectURI + "?" + params.Encode()
		ctx.Redirect(consts.StatusTemporaryRedirect, []byte(redirectURL))
		return
	}

	// Web app - set cookies and redirect to home
	ctx.SetCookie(
		"access_token",
		resp.Tokens.AccessToken,
		int(time.Until(resp.Tokens.ExpiresAt).Seconds()),
		"/",
		"",
		protocol.CookieSameSiteLaxMode,
		true,
		false,
	)

	ctx.SetCookie(
		"refresh_token",
		resp.Tokens.RefreshToken,
		refreshTokenMaxAge,
		"/",
		"",
		protocol.CookieSameSiteLaxMode,
		true,
		false,
	)

	// Redirect to frontend
	ctx.Redirect(consts.StatusTemporaryRedirect, []byte("/"))
}

// generateState generates a random state for OAuth
func (h *OAuthHandler) generateState() string {
	b := make([]byte, stateTokenBytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// validateState validates OAuth state parameter and returns redirect URI
func (h *OAuthHandler) validateState(state string) (string, bool) {
	// Clean up expired states
	now := time.Now()
	for s, oauthState := range h.states {
		if now.After(oauthState.Expiry) {
			delete(h.states, s)
		}
	}

	// Check if state exists and is not expired
	oauthState, exists := h.states[state]
	if !exists || now.After(oauthState.Expiry) {
		return "", false
	}

	// Get redirect URI before deleting state
	redirectURI := oauthState.RedirectURI

	// Delete used state
	delete(h.states, state)
	return redirectURI, true
}
