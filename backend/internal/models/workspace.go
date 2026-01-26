package models

import (
	"time"

	"github.com/google/uuid"
)

// WorkspaceRole defines the role of a user in a workspace
type WorkspaceRole string

const (
	WorkspaceRoleOwner  WorkspaceRole = "owner"
	WorkspaceRoleEditor WorkspaceRole = "editor"
	WorkspaceRoleViewer WorkspaceRole = "viewer"
)

// Workspace represents a collaborative workspace
type Workspace struct {
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	Description  *string                `json:"description,omitempty"`
	ThumbnailURL *string                `json:"thumbnail_url,omitempty"`
	Settings     map[string]interface{} `json:"settings"`
	DeletedAt    *time.Time             `json:"deleted_at,omitempty"`
	Name         string                 `json:"name"`
	ID           uuid.UUID              `json:"id"`
	OwnerID      uuid.UUID              `json:"owner_id"`
	IsPublic     bool                   `json:"is_public"`
}

// WorkspaceMember represents a user's membership in a workspace
type WorkspaceMember struct {
	JoinedAt    time.Time     `json:"joined_at"`
	InvitedBy   *uuid.UUID    `json:"invited_by,omitempty"`
	Role        WorkspaceRole `json:"role"`
	ID          uuid.UUID     `json:"id"`
	WorkspaceID uuid.UUID     `json:"workspace_id"`
	UserID      uuid.UUID     `json:"user_id"`
}

// WorkspaceInvite represents an invitation to join a workspace
type WorkspaceInvite struct {
	ExpiresAt   time.Time     `json:"expires_at"`
	CreatedAt   time.Time     `json:"created_at"`
	AcceptedAt  *time.Time    `json:"accepted_at,omitempty"`
	AcceptedBy  *uuid.UUID    `json:"accepted_by,omitempty"`
	Email       string        `json:"email"`
	Role        WorkspaceRole `json:"role"`
	TokenHash   string        `json:"-"`
	ID          uuid.UUID     `json:"id"`
	WorkspaceID uuid.UUID     `json:"workspace_id"`
	CreatedBy   uuid.UUID     `json:"created_by"`
}

// WorkspaceWithRole extends Workspace with user's role
type WorkspaceWithRole struct {
	Owner    *User         `json:"owner,omitempty"`
	UserRole WorkspaceRole `json:"user_role"`
	Workspace
}

// WorkspaceMemberWithUser extends WorkspaceMember with user details
type WorkspaceMemberWithUser struct {
	User User `json:"user"`
	WorkspaceMember
}

// --- Request DTOs ---

// CreateWorkspaceRequest represents a request to create a new workspace
type CreateWorkspaceRequest struct {
	Description *string                `json:"description,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Name        string                 `json:"name" binding:"required,min=1,max=255"`
	IsPublic    bool                   `json:"is_public"`
}

// UpdateWorkspaceRequest represents a request to update workspace
type UpdateWorkspaceRequest struct {
	Name         *string                `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	Description  *string                `json:"description,omitempty"`
	IsPublic     *bool                  `json:"is_public,omitempty"`
	ThumbnailURL *string                `json:"thumbnail_url,omitempty"`
	Settings     map[string]interface{} `json:"settings,omitempty"`
}

// InviteToWorkspaceRequest represents a request to invite a user to workspace
type InviteToWorkspaceRequest struct {
	Email string        `json:"email" binding:"required,email"`
	Role  WorkspaceRole `json:"role" binding:"required,oneof=editor viewer"`
}

// AcceptInviteRequest represents a request to accept workspace invitation
type AcceptInviteRequest struct {
	Token string `json:"token" binding:"required"`
}

// UpdateMemberRoleRequest represents a request to update member's role
type UpdateMemberRoleRequest struct {
	Role WorkspaceRole `json:"role" binding:"required,oneof=owner editor viewer"`
}

// WorkspaceListFilter represents filters for listing workspaces
type WorkspaceListFilter struct {
	Query      string `form:"q"`
	SortBy     string `form:"sort_by"`
	SortOrder  string `form:"sort_order"`
	Limit      int    `form:"limit"`
	Offset     int    `form:"offset"`
	OwnedOnly  bool   `form:"owned_only"`
	SharedOnly bool   `form:"shared_only"`
}

// --- Response DTOs ---

// WorkspaceResponse represents workspace data in API responses
type WorkspaceResponse struct {
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	Description  *string                `json:"description,omitempty"`
	ThumbnailURL *string                `json:"thumbnail_url,omitempty"`
	Settings     map[string]interface{} `json:"settings"`
	UserRole     *WorkspaceRole         `json:"user_role,omitempty"`
	Owner        *UserResponse          `json:"owner,omitempty"`
	Name         string                 `json:"name"`
	ID           uuid.UUID              `json:"id"`
	OwnerID      uuid.UUID              `json:"owner_id"`
	IsPublic     bool                   `json:"is_public"`
}

// WorkspaceListResponse represents paginated list of workspaces
type WorkspaceListResponse struct {
	Workspaces []WorkspaceResponse `json:"workspaces"`
	Total      int                 `json:"total"`
	Limit      int                 `json:"limit"`
	Offset     int                 `json:"offset"`
}

// WorkspaceMemberResponse represents workspace member in API responses
//
//nolint:govet // Field order optimized for JSON readability, not memory alignment
type WorkspaceMemberResponse struct {
	JoinedAt time.Time     `json:"joined_at"`
	Role     WorkspaceRole `json:"role"`
	ID       uuid.UUID     `json:"id"`
	User     *UserResponse `json:"user"`
}

// WorkspaceInviteResponse represents workspace invite in API responses
//
//nolint:govet // Field order optimized for JSON readability, not memory alignment
type WorkspaceInviteResponse struct {
	ExpiresAt time.Time     `json:"expires_at"`
	CreatedAt time.Time     `json:"created_at"`
	Email     string        `json:"email"`
	Role      WorkspaceRole `json:"role"`
	ID        uuid.UUID     `json:"id"`
	CreatedBy *UserResponse `json:"created_by"`
}

// InviteTokenResponse represents response with invitation token
type InviteTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	InviteURL string    `json:"invite_url"`
}
