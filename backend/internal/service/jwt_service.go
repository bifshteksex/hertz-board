package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/bifshteksex/hertzboard/internal/config"
)

//nolint:govet // fieldalignment: struct field order optimized for readability
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService struct {
	secret               string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// NewJWTService creates a new JWT service
func NewJWTService(cfg *config.JWTConfig) (*JWTService, error) {
	accessDuration, err := cfg.GetAccessTokenDuration()
	if err != nil {
		return nil, fmt.Errorf("invalid access token duration: %w", err)
	}

	refreshDuration, err := cfg.GetRefreshTokenDuration()
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token duration: %w", err)
	}

	return &JWTService{
		secret:               cfg.Secret,
		accessTokenDuration:  accessDuration,
		refreshTokenDuration: refreshDuration,
	}, nil
}

// GenerateAccessToken generates a new access token
func (s *JWTService) GenerateAccessToken(userID uuid.UUID, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.accessTokenDuration)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "hertzboard",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken generates a new refresh token
func (s *JWTService) GenerateRefreshToken() (token string, tokenHash string, expiresAt time.Time, err error) {
	token = uuid.New().String()
	tokenHash = hashToken(token)
	expiresAt = time.Now().Add(s.refreshTokenDuration)

	return token, tokenHash, expiresAt, nil
}

// ValidateAccessToken validates an access token and returns the claims
func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// HashRefreshToken hashes a refresh token for storage
func (s *JWTService) HashRefreshToken(token string) string {
	return hashToken(token)
}

// GetRefreshTokenDuration returns the refresh token duration
func (s *JWTService) GetRefreshTokenDuration() time.Duration {
	return s.refreshTokenDuration
}

// hashToken creates a SHA-256 hash of a token
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
