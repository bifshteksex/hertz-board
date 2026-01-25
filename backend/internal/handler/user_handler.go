package handler

import (
	"context"

	"github.com/bifshteksex/hertzboard/internal/models"
	"github.com/bifshteksex/hertzboard/internal/repository"
	"github.com/bifshteksex/hertzboard/internal/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler handles user-related endpoints
type UserHandler struct {
	userRepo    *repository.UserRepository
	authService *service.AuthService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userRepo *repository.UserRepository, authService *service.AuthService) *UserHandler {
	return &UserHandler{
		userRepo:    userRepo,
		authService: authService,
	}
}

// GetProfile returns the current user's profile
func (h *UserHandler) GetProfile(c context.Context, ctx *app.RequestContext) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized",
		})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := h.userRepo.GetByID(c, uid)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to get user",
		})
		return
	}

	if user == nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"error": "User not found",
		})
		return
	}

	ctx.JSON(consts.StatusOK, user)
}

// UpdateProfile updates the current user's profile
func (h *UserHandler) UpdateProfile(c context.Context, ctx *app.RequestContext) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized",
		})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "Invalid user ID",
		})
		return
	}

	var req models.UpdateProfileRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Get current user
	user, err := h.userRepo.GetByID(c, uid)
	if err != nil || user == nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to get user",
		})
		return
	}

	// Update fields
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if err := h.userRepo.Update(c, user); err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to update profile",
		})
		return
	}

	ctx.JSON(consts.StatusOK, user)
}

// ChangePassword changes the current user's password
func (h *UserHandler) ChangePassword(c context.Context, ctx *app.RequestContext) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"error": "Unauthorized",
		})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "Invalid user ID",
		})
		return
	}

	var req models.ChangePasswordRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Get current user
	user, err := h.userRepo.GetByID(c, uid)
	if err != nil || user == nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to get user",
		})
		return
	}

	// Check if user has a password (not OAuth only)
	if user.PasswordHash == nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": "User registered with OAuth, password change not available",
		})
		return
	}

	// Verify old password
	if err := verifyPassword(*user.PasswordHash, req.OldPassword); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": "Invalid old password",
		})
		return
	}

	// Hash new password
	newHash, err := hashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to hash password",
		})
		return
	}

	// Update password
	if err := h.userRepo.UpdatePassword(c, uid, newHash); err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to update password",
		})
		return
	}

	// Delete all refresh tokens (logout all sessions)
	if err := h.userRepo.DeleteUserRefreshTokens(c, uid); err != nil {
		// Log error but don't fail the request
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"message": "Password changed successfully",
	})
}

// Helper functions
func hashPassword(password string) (string, error) {
	// This should use the same function from auth_service
	// For now, importing golang.org/x/crypto/bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func verifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
