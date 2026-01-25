package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/bifshteksex/hertzboard/internal/config"
	"github.com/bifshteksex/hertzboard/internal/models"
	"github.com/bifshteksex/hertzboard/internal/repository"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// OAuthService handles OAuth authentication
type OAuthService struct {
	userRepo   *repository.UserRepository
	jwtService *JWTService
	googleCfg  *oauth2.Config
	githubCfg  *oauth2.Config
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(
	cfg *config.OAuthConfig,
	userRepo *repository.UserRepository,
	jwtService *JWTService,
) *OAuthService {
	googleCfg := &oauth2.Config{
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.ClientSecret,
		RedirectURL:  cfg.Google.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	githubCfg := &oauth2.Config{
		ClientID:     cfg.GitHub.ClientID,
		ClientSecret: cfg.GitHub.ClientSecret,
		RedirectURL:  cfg.GitHub.RedirectURL,
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}

	return &OAuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
		googleCfg:  googleCfg,
		githubCfg:  githubCfg,
	}
}

// GetGoogleAuthURL returns the Google OAuth authorization URL
func (s *OAuthService) GetGoogleAuthURL(state string) string {
	return s.googleCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// GetGitHubAuthURL returns the GitHub OAuth authorization URL
func (s *OAuthService) GetGitHubAuthURL(state string) string {
	return s.githubCfg.AuthCodeURL(state)
}

// GoogleCallback handles Google OAuth callback
func (s *OAuthService) GoogleCallback(ctx context.Context, code string) (*models.AuthResponse, error) {
	// Exchange code for token
	token, err := s.googleCfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from Google
	client := s.googleCfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	// Find or create user
	return s.findOrCreateUser(ctx, "google", userInfo.ID, userInfo.Email, userInfo.Name, userInfo.Picture)
}

// GitHubCallback handles GitHub OAuth callback
func (s *OAuthService) GitHubCallback(ctx context.Context, code string) (*models.AuthResponse, error) {
	// Exchange code for token
	token, err := s.githubCfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from GitHub
	client := s.githubCfg.Client(ctx, token)

	// Get user profile
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo struct {
		ID        int64  `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	// If email is not public, fetch it separately
	if userInfo.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			emailBody, err := io.ReadAll(emailResp.Body)
			if err == nil {
				var emails []struct {
					Email   string `json:"email"`
					Primary bool   `json:"primary"`
				}
				if json.Unmarshal(emailBody, &emails) == nil {
					for _, e := range emails {
						if e.Primary {
							userInfo.Email = e.Email
							break
						}
					}
				}
			}
		}
	}

	if userInfo.Email == "" {
		return nil, fmt.Errorf("failed to get email from GitHub")
	}

	providerID := fmt.Sprintf("%d", userInfo.ID)
	name := userInfo.Name
	if name == "" {
		name = userInfo.Email
	}

	// Find or create user
	return s.findOrCreateUser(ctx, "github", providerID, userInfo.Email, name, userInfo.AvatarURL)
}

// findOrCreateUser finds existing user or creates a new one
func (s *OAuthService) findOrCreateUser(
	ctx context.Context,
	provider, providerID, email, name, avatarURL string,
) (*models.AuthResponse, error) {
	// Try to find user by provider
	user, err := s.userRepo.GetByProvider(ctx, provider, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by provider: %w", err)
	}

	// If not found, try by email
	if user == nil {
		user, err = s.userRepo.GetByEmail(ctx, email)
		if err != nil {
			return nil, fmt.Errorf("failed to get user by email: %w", err)
		}
	}

	// Create new user if doesn't exist
	if user == nil {
		user = &models.User{
			Email:         email,
			Name:          name,
			Provider:      provider,
			ProviderID:    &providerID,
			EmailVerified: true, // OAuth emails are already verified
		}

		if avatarURL != "" {
			user.AvatarURL = &avatarURL
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Generate tokens
	accessToken, expiresAt, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshHash, refreshExpiresAt, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in DB
	dbToken := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: refreshExpiresAt,
	}

	if err := s.userRepo.CreateRefreshToken(ctx, dbToken); err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &models.AuthResponse{
		User: user,
		Tokens: &models.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    expiresAt,
		},
	}, nil
}
