package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/bifshteksex/hertz-board/internal/models"
	"github.com/bifshteksex/hertz-board/internal/repository"
)

// AuthService handles authentication logic
type AuthService struct {
	userRepo   *repository.UserRepository
	jwtService *JWTService
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository, jwtService *JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req *models.CreateUserRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:         req.Email,
		PasswordHash:  &passwordHash,
		Name:          req.Name,
		Provider:      "email",
		EmailVerified: false,
	}

	if createErr := s.userRepo.Create(ctx, user); createErr != nil {
		return nil, fmt.Errorf("failed to create user: %w", createErr)
	}

	// Generate tokens
	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &models.AuthResponse{
		User:   user,
		Tokens: tokens,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user registered with OAuth
	if user.PasswordHash == nil {
		return nil, fmt.Errorf("user registered with %s, please use OAuth login", user.Provider)
	}

	// Verify password
	if !verifyPassword(*user.PasswordHash, req.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate tokens
	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &models.AuthResponse{
		User:   user,
		Tokens: tokens,
	}, nil
}

// RefreshToken refreshes access token using refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	// Hash the refresh token
	tokenHash := s.jwtService.HashRefreshToken(refreshToken)

	// Get refresh token from DB
	token, err := s.userRepo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	if token == nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Delete old refresh token
	if deleteErr := s.userRepo.DeleteRefreshToken(ctx, tokenHash); deleteErr != nil {
		return nil, fmt.Errorf("failed to delete old refresh token: %w", deleteErr)
	}

	// Generate new token pair
	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokens, nil
}

// Logout invalidates a refresh token
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := s.jwtService.HashRefreshToken(refreshToken)
	return s.userRepo.DeleteRefreshToken(ctx, tokenHash)
}

// ForgotPassword creates a password reset token
func (s *AuthService) ForgotPassword(ctx context.Context, email string) (string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		// Don't reveal if user exists
		return "", nil
	}

	// Check if user has a password (not OAuth only)
	if user.PasswordHash == nil {
		return "", fmt.Errorf("user registered with %s, password reset not available", user.Provider)
	}

	// Generate reset token
	token := uuid.New().String()
	tokenHash := hashToken(token)
	expiresAt := time.Now().Add(1 * time.Hour)

	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	if err := s.userRepo.CreatePasswordResetToken(ctx, resetToken); err != nil {
		return "", fmt.Errorf("failed to create password reset token: %w", err)
	}

	return token, nil
}

// ResetPassword resets user password using a reset token
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Hash token
	tokenHash := hashToken(token)

	// Get reset token from DB
	resetToken, err := s.userRepo.GetPasswordResetToken(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to get password reset token: %w", err)
	}
	if resetToken == nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Hash new password
	passwordHash, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := s.userRepo.UpdatePassword(ctx, resetToken.UserID, passwordHash); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	if err := s.userRepo.MarkPasswordResetTokenUsed(ctx, tokenHash); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Delete all refresh tokens for the user (logout all sessions)
	if err := s.userRepo.DeleteUserRefreshTokens(ctx, resetToken.UserID); err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	return nil
}

// generateTokenPair generates access and refresh token pair
func (s *AuthService) generateTokenPair(ctx context.Context, user *models.User) (*models.TokenPair, error) {
	// Generate access token
	accessToken, expiresAt, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, refreshHash, refreshExpiresAt, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Store refresh token in DB
	dbToken := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: refreshExpiresAt,
	}

	if err := s.userRepo.CreateRefreshToken(ctx, dbToken); err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// verifyPassword verifies a password against its hash
func verifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
