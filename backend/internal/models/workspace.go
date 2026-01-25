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
//
//nolint:govet // fieldalignment: struct field order optimized for readability
type Workspace struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	Description  *string                `json:"description,omitempty"`
	OwnerID      uuid.UUID              `json:"owner_id"`
	ThumbnailURL *string                `json:"thumbnail_url,omitempty"`
	IsPublic     bool                   `json:"is_public"`
	Settings     map[string]interface{} `json:"settings"`
	DeletedAt    *time.Time             `json:"deleted_at,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// WorkspaceMember represents a user's membership in a workspace
//
//nolint:govet // fieldalignment: struct field order optimized for readability
type WorkspaceMember struct {
	ID          uuid.UUID     `json:"id"`
	WorkspaceID uuid.UUID     `json:"workspace_id"`
	UserID      uuid.UUID     `json:"user_id"`
	Role        WorkspaceRole `json:"role"`
	InvitedBy   *uuid.UUID    `json:"invited_by,omitempty"`
	JoinedAt    time.Time     `json:"joined_at"`
}

// WorkspaceInvite represents an invitation to join a workspace
//
//nolint:govet // fieldalignment: struct field order optimized for readability
type WorkspaceInvite struct {
	ID          uuid.UUID     `json:"id"`
	WorkspaceID uuid.UUID     `json:"workspace_id"`
	Email       string        `json:"email"`
	Role        WorkspaceRole `json:"role"`
	TokenHash   string        `json:"-"` // Never expose token hash
	ExpiresAt   time.Time     `json:"expires_at"`
	CreatedBy   uuid.UUID     `json:"created_by"`
	CreatedAt   time.Time     `json:"created_at"`
	AcceptedAt  *time.Time    `json:"accepted_at,omitempty"`
	AcceptedBy  *uuid.UUID    `json:"accepted_by,omitempty"`
}

// WorkspaceWithRole extends Workspace with user's role
//
//nolint:govet // fieldalignment: struct field order optimized for readability
type WorkspaceWithRole struct {
	Workspace
	UserRole WorkspaceRole `json:"user_role"`
	Owner    *User         `json:"owner,omitempty"`
}

// WorkspaceMemberWithUser extends WorkspaceMember with user details
type WorkspaceMemberWithUser struct {
	WorkspaceMember
	User User `json:"user"`
}

// --- Request DTOs ---

// CreateWorkspaceRequest represents a request to create a new workspace
//
//nolint:govet // fieldalignment: struct field order optimized for readability
type CreateWorkspaceRequest struct {
	Name        string                 `json:"name" binding:"required,min=1,max=255"`
	Description *string                `json:"description,omitempty"`
	IsPublic    bool                   `json:"is_public"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
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
//
//nolint:govet // fieldalignment: struct field order optimized for readability
type WorkspaceListFilter struct {
	Query      string `form:"q"`
	OwnedOnly  bool   `form:"owned_only"`
	SharedOnly bool   `form:"shared_only"`
	SortBy     string `form:"sort_by"`    // created_at, updated_at, name
	SortOrder  string `form:"sort_order"` // asc, desc
	Limit      int    `form:"limit"`
	Offset     int    `form:"offset"`
}

// --- Response DTOs ---

// WorkspaceResponse represents workspace data in API responses
//
//nolint:govet // fieldalignment: struct field order optimized for readability
type WorkspaceResponse struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	Description  *string                `json:"description,omitempty"`
	OwnerID      uuid.UUID              `json:"owner_id"`
	ThumbnailURL *string                `json:"thumbnail_url,omitempty"`
	IsPublic     bool                   `json:"is_public"`
	Settings     map[string]interface{} `json:"settings"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	UserRole     *WorkspaceRole         `json:"user_role,omitempty"`
	Owner        *UserResponse          `json:"owner,omitempty"`
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
//nolint:govet // fieldalignment: struct field order optimized for readability
type WorkspaceMemberResponse struct {
	ID       uuid.UUID     `json:"id"`
	User     UserResponse  `json:"user"`
	Role     WorkspaceRole `json:"role"`
	JoinedAt time.Time     `json:"joined_at"`
}

// WorkspaceInviteResponse represents workspace invite in API responses
//
//nolint:govet // fieldalignment: struct field order optimized for readability
type WorkspaceInviteResponse struct {
	ID        uuid.UUID     `json:"id"`
	Email     string        `json:"email"`
	Role      WorkspaceRole `json:"role"`
	ExpiresAt time.Time     `json:"expires_at"`
	CreatedAt time.Time     `json:"created_at"`
	CreatedBy UserResponse  `json:"created_by"`
}

// InviteTokenResponse represents response with invitation token
type InviteTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	InviteURL string    `json:"invite_url"`
}
