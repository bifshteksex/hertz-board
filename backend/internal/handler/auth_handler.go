package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/bifshteksex/hertzboard/internal/models"
	"github.com/bifshteksex/hertzboard/internal/service"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c context.Context, ctx *app.RequestContext) {
	var req models.CreateUserRequest
	resp, statusCode, err := h.bindValidateAndExecute(ctx, &req, func() (interface{}, error) {
		return h.authService.Register(c, &req)
	})

	if err != nil {
		ctx.JSON(statusCode, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusCreated, resp)
}

// Login handles user login
func (h *AuthHandler) Login(c context.Context, ctx *app.RequestContext) {
	var req models.LoginRequest
	resp, statusCode, err := h.bindValidateAndExecute(ctx, &req, func() (interface{}, error) {
		return h.authService.Login(c, &req)
	})

	if err != nil {
		ctx.JSON(statusCode, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, resp)
}

// bindValidateAndExecute handles common request binding, validation, and execution pattern
func (h *AuthHandler) bindValidateAndExecute(
	ctx *app.RequestContext,
	req interface{},
	execute func() (interface{}, error),
) (resp interface{}, statusCode int, err error) {
	if err = ctx.BindAndValidate(req); err != nil {
		return nil, consts.StatusBadRequest, err
	}

	resp, err = execute()
	if err != nil {
		// Determine status code based on error context
		// For now, return BadRequest, but caller can override
		return nil, consts.StatusBadRequest, err
	}

	return resp, consts.StatusOK, nil
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c context.Context, ctx *app.RequestContext) {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	var req RefreshRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	tokens, err := h.authService.RefreshToken(c, req.RefreshToken)
	if err != nil {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, tokens)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c context.Context, ctx *app.RequestContext) {
	type LogoutRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	var req LogoutRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	if err := h.authService.Logout(c, req.RefreshToken); err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to logout",
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"message": "Logged out successfully",
	})
}

// ForgotPassword handles forgot password requests
func (h *AuthHandler) ForgotPassword(c context.Context, ctx *app.RequestContext) {
	var req models.ForgotPasswordRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// TODO: Send email with reset token (Task 1.7)
	_, err := h.authService.ForgotPassword(c, req.Email)
	if err != nil {
		// Don't reveal if email exists
		ctx.JSON(consts.StatusOK, map[string]interface{}{
			"message": "If the email exists, a password reset link has been sent",
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"message": "If the email exists, a password reset link has been sent",
	})
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c context.Context, ctx *app.RequestContext) {
	var req models.ResetPasswordRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	if err := h.authService.ResetPassword(c, req.Token, req.NewPassword); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"message": "Password reset successfully",
	})
}
