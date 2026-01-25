package models

import (
	"time"

	"github.com/google/uuid"
)

//nolint:govet // fieldalignment: struct field order optimized for readability
type User struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Email         string    `json:"email" db:"email"`
	PasswordHash  *string   `json:"-" db:"password_hash"`
	Name          string    `json:"name" db:"name"`
	AvatarURL     *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	Provider      string    `json:"provider" db:"provider"`
	ProviderID    *string   `json:"-" db:"provider_id"`
	EmailVerified bool      `json:"email_verified" db:"email_verified"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

//nolint:govet // fieldalignment: struct field order optimized for readability
type RefreshToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	TokenHash string    `json:"-" db:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

//nolint:govet // fieldalignment: struct field order optimized for readability
type PasswordResetToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	TokenHash string     `json:"-" db:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UsedAt    *time.Time `json:"used_at,omitempty" db:"used_at"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest represents the update profile request
type UpdateProfileRequest struct {
	Name      *string `json:"name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// ChangePasswordRequest represents the change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ForgotPasswordRequest represents the forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

//nolint:govet // fieldalignment: struct field order optimized for readability
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User   *User      `json:"user"`
	Tokens *TokenPair `json:"tokens"`
}
