package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bifshteksex/hertzboard/internal/models"
	"github.com/bifshteksex/hertzboard/internal/repository"

	"github.com/google/uuid"
)

type WorkspaceService struct {
	workspaceRepo *repository.WorkspaceRepository
	userRepo      *repository.UserRepository
	emailService  *EmailService
}

func NewWorkspaceService(
	workspaceRepo *repository.WorkspaceRepository,
	userRepo *repository.UserRepository,
	emailService *EmailService,
) *WorkspaceService {
	return &WorkspaceService{
		workspaceRepo: workspaceRepo,
		userRepo:      userRepo,
		emailService:  emailService,
	}
}

// --- Workspace CRUD ---

// CreateWorkspace creates a new workspace with the user as owner
func (s *WorkspaceService) CreateWorkspace(ctx context.Context, req *models.CreateWorkspaceRequest, ownerID uuid.UUID) (*models.Workspace, error) {
	workspace := &models.Workspace{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
		IsPublic:    req.IsPublic,
		Settings:    req.Settings,
	}

	if workspace.Settings == nil {
		workspace.Settings = make(map[string]interface{})
	}

	if err := s.workspaceRepo.CreateWorkspace(ctx, workspace); err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	return workspace, nil
}

// GetWorkspace retrieves a workspace by ID
func (s *WorkspaceService) GetWorkspace(ctx context.Context, id uuid.UUID) (*models.Workspace, error) {
	workspace, err := s.workspaceRepo.GetWorkspaceByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}

	if workspace == nil {
		return nil, fmt.Errorf("workspace not found")
	}

	return workspace, nil
}

// GetWorkspaceWithRole retrieves workspace with user's role
func (s *WorkspaceService) GetWorkspaceWithRole(ctx context.Context, workspaceID, userID uuid.UUID) (*models.WorkspaceWithRole, error) {
	workspace, err := s.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	member, err := s.workspaceRepo.GetMember(ctx, workspaceID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get member: %w", err)
	}

	if member == nil {
		// Check if workspace is public
		if !workspace.IsPublic {
			return nil, fmt.Errorf("access denied")
		}
		// Public workspace, viewer role
		return &models.WorkspaceWithRole{
			Workspace: *workspace,
			UserRole:  models.WorkspaceRoleViewer,
		}, nil
	}

	// Get owner info
	owner, err := s.userRepo.GetByID(ctx, workspace.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner: %w", err)
	}

	return &models.WorkspaceWithRole{
		Workspace: *workspace,
		UserRole:  member.Role,
		Owner:     owner,
	}, nil
}

// UpdateWorkspace updates workspace information
func (s *WorkspaceService) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, req *models.UpdateWorkspaceRequest) (*models.Workspace, error) {
	workspace, err := s.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.Name != nil {
		workspace.Name = *req.Name
	}
	if req.Description != nil {
		workspace.Description = req.Description
	}
	if req.IsPublic != nil {
		workspace.IsPublic = *req.IsPublic
	}
	if req.ThumbnailURL != nil {
		workspace.ThumbnailURL = req.ThumbnailURL
	}
	if req.Settings != nil {
		workspace.Settings = req.Settings
	}

	if err := s.workspaceRepo.UpdateWorkspace(ctx, workspace); err != nil {
		return nil, fmt.Errorf("failed to update workspace: %w", err)
	}

	return workspace, nil
}

// DeleteWorkspace soft deletes a workspace
func (s *WorkspaceService) DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) error {
	if err := s.workspaceRepo.SoftDeleteWorkspace(ctx, workspaceID); err != nil {
		return fmt.Errorf("failed to delete workspace: %w", err)
	}

	return nil
}

// ListUserWorkspaces retrieves all workspaces accessible to user
func (s *WorkspaceService) ListUserWorkspaces(ctx context.Context, userID uuid.UUID, filter models.WorkspaceListFilter) (*models.WorkspaceListResponse, error) {
	workspaces, total, err := s.workspaceRepo.ListWorkspacesByUser(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list workspaces: %w", err)
	}

	// Convert to response format
	response := &models.WorkspaceListResponse{
		Workspaces: make([]models.WorkspaceResponse, 0, len(workspaces)),
		Total:      total,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
	}

	for _, ws := range workspaces {
		// Get owner info
		owner, err := s.userRepo.GetByID(ctx, ws.OwnerID)
		if err != nil {
			continue // Skip on error
		}

		wsResp := models.WorkspaceResponse{
			ID:           ws.ID,
			Name:         ws.Name,
			Description:  ws.Description,
			OwnerID:      ws.OwnerID,
			ThumbnailURL: ws.ThumbnailURL,
			IsPublic:     ws.IsPublic,
			Settings:     ws.Settings,
			CreatedAt:    ws.CreatedAt,
			UpdatedAt:    ws.UpdatedAt,
			UserRole:     &ws.UserRole,
		}

		if owner != nil {
			wsResp.Owner = &models.UserResponse{
				ID:        owner.ID,
				Email:     owner.Email,
				Name:      owner.Name,
				AvatarURL: owner.AvatarURL,
			}
		}

		response.Workspaces = append(response.Workspaces, wsResp)
	}

	return response, nil
}

// DuplicateWorkspace creates a copy of a workspace
func (s *WorkspaceService) DuplicateWorkspace(ctx context.Context, workspaceID, userID uuid.UUID) (*models.Workspace, error) {
	// Get original workspace
	original, err := s.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	// Create new workspace
	newWorkspace := &models.Workspace{
		ID:          uuid.New(),
		Name:        original.Name + " (Copy)",
		Description: original.Description,
		OwnerID:     userID,
		IsPublic:    false, // Copies are private by default
		Settings:    original.Settings,
	}

	if err := s.workspaceRepo.CreateWorkspace(ctx, newWorkspace); err != nil {
		return nil, fmt.Errorf("failed to duplicate workspace: %w", err)
	}

	// TODO: Copy canvas elements (will be implemented in Phase 3)

	return newWorkspace, nil
}

// --- Member Management ---

// GetMembers retrieves all members of a workspace
func (s *WorkspaceService) GetMembers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceMemberResponse, error) {
	members, err := s.workspaceRepo.ListMembers(ctx, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}

	response := make([]models.WorkspaceMemberResponse, 0, len(members))
	for _, m := range members {
		response = append(response, models.WorkspaceMemberResponse{
			ID: m.ID,
			User: models.UserResponse{
				ID:        m.User.ID,
				Email:     m.User.Email,
				Name:      m.User.Name,
				AvatarURL: m.User.AvatarURL,
			},
			Role:     m.Role,
			JoinedAt: m.JoinedAt,
		})
	}

	return response, nil
}

// UpdateMemberRole updates a member's role
func (s *WorkspaceService) UpdateMemberRole(ctx context.Context, workspaceID, memberUserID uuid.UUID, role models.WorkspaceRole) error {
	// Prevent changing owner role
	workspace, err := s.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return err
	}

	if workspace.OwnerID == memberUserID && role != models.WorkspaceRoleOwner {
		return fmt.Errorf("cannot change owner's role")
	}

	if err := s.workspaceRepo.UpdateMemberRole(ctx, workspaceID, memberUserID, role); err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}

	return nil
}

// RemoveMember removes a member from workspace
func (s *WorkspaceService) RemoveMember(ctx context.Context, workspaceID, memberUserID uuid.UUID) error {
	// Prevent removing owner
	workspace, err := s.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return err
	}

	if workspace.OwnerID == memberUserID {
		return fmt.Errorf("cannot remove workspace owner")
	}

	if err := s.workspaceRepo.RemoveMember(ctx, workspaceID, memberUserID); err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	return nil
}

// --- Invitations ---

// CreateInvite creates a new workspace invitation
func (s *WorkspaceService) CreateInvite(ctx context.Context, workspaceID, createdBy uuid.UUID, req *models.InviteToWorkspaceRequest) (*models.InviteTokenResponse, error) {
	// Check if user already exists and is a member
	user, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if user != nil {
		member, _ := s.workspaceRepo.GetMember(ctx, workspaceID, user.ID)
		if member != nil {
			return nil, fmt.Errorf("user is already a member")
		}
	}

	// Check if there's already a pending invite
	existingInvite, _ := s.workspaceRepo.GetInviteByWorkspaceAndEmail(ctx, workspaceID, req.Email)
	if existingInvite != nil {
		return nil, fmt.Errorf("invitation already sent to this email")
	}

	// Generate invite token
	token := uuid.New().String()
	tokenHash := hashToken(token)

	invite := &models.WorkspaceInvite{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		Email:       req.Email,
		Role:        req.Role,
		TokenHash:   tokenHash,
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedBy:   createdBy,
	}

	if err := s.workspaceRepo.CreateInvite(ctx, invite); err != nil {
		return nil, fmt.Errorf("failed to create invite: %w", err)
	}

	// Get workspace details for email
	workspace, _ := s.GetWorkspace(ctx, workspaceID)
	creator, _ := s.userRepo.GetByID(ctx, createdBy)

	// Send invitation email
	if workspace != nil && creator != nil {
		s.emailService.SendWorkspaceInvite(req.Email, workspace.Name, creator.Name, token)
	}

	// Build invite URL (frontend route)
	inviteURL := fmt.Sprintf("/workspace/invite?token=%s", token)

	return &models.InviteTokenResponse{
		Token:     token,
		ExpiresAt: invite.ExpiresAt,
		InviteURL: inviteURL,
	}, nil
}

// AcceptInvite accepts a workspace invitation
func (s *WorkspaceService) AcceptInvite(ctx context.Context, token string, userID uuid.UUID) (*models.Workspace, error) {
	tokenHash := hashToken(token)

	invite, err := s.workspaceRepo.GetInviteByToken(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get invite: %w", err)
	}

	if invite == nil {
		return nil, fmt.Errorf("invalid or expired invitation")
	}

	// Check if already accepted
	if invite.AcceptedAt != nil {
		return nil, fmt.Errorf("invitation already accepted")
	}

	// Check if expired
	if time.Now().After(invite.ExpiresAt) {
		return nil, fmt.Errorf("invitation has expired")
	}

	// Get user to verify email matches
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.Email != invite.Email {
		return nil, fmt.Errorf("invitation email does not match your account")
	}

	// Check if already a member
	member, _ := s.workspaceRepo.GetMember(ctx, invite.WorkspaceID, userID)
	if member != nil {
		return nil, fmt.Errorf("you are already a member of this workspace")
	}

	// Add user as member
	newMember := &models.WorkspaceMember{
		ID:          uuid.New(),
		WorkspaceID: invite.WorkspaceID,
		UserID:      userID,
		Role:        invite.Role,
		InvitedBy:   &invite.CreatedBy,
	}

	if err := s.workspaceRepo.AddMember(ctx, newMember); err != nil {
		return nil, fmt.Errorf("failed to add member: %w", err)
	}

	// Mark invite as accepted
	if err := s.workspaceRepo.MarkInviteAsAccepted(ctx, invite.ID, userID); err != nil {
		return nil, fmt.Errorf("failed to mark invite as accepted: %w", err)
	}

	// Get workspace
	workspace, err := s.GetWorkspace(ctx, invite.WorkspaceID)
	if err != nil {
		return nil, err
	}

	return workspace, nil
}

// GetPendingInvites retrieves all pending invitations for a workspace
func (s *WorkspaceService) GetPendingInvites(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceInviteResponse, error) {
	invites, err := s.workspaceRepo.ListPendingInvites(ctx, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending invites: %w", err)
	}

	response := make([]models.WorkspaceInviteResponse, 0, len(invites))
	for _, inv := range invites {
		// Get creator info
		creator, err := s.userRepo.GetByID(ctx, inv.CreatedBy)
		if err != nil {
			continue
		}

		response = append(response, models.WorkspaceInviteResponse{
			ID:        inv.ID,
			Email:     inv.Email,
			Role:      inv.Role,
			ExpiresAt: inv.ExpiresAt,
			CreatedAt: inv.CreatedAt,
			CreatedBy: models.UserResponse{
				ID:        creator.ID,
				Email:     creator.Email,
				Name:      creator.Name,
				AvatarURL: creator.AvatarURL,
			},
		})
	}

	return response, nil
}

// RevokeInvite revokes a pending invitation
func (s *WorkspaceService) RevokeInvite(ctx context.Context, inviteID uuid.UUID) error {
	if err := s.workspaceRepo.RevokeInvite(ctx, inviteID); err != nil {
		return fmt.Errorf("failed to revoke invite: %w", err)
	}

	return nil
}

// --- Permissions ---

// CheckPermission checks if user has required permission level
func (s *WorkspaceService) CheckPermission(ctx context.Context, workspaceID, userID uuid.UUID, requiredRole models.WorkspaceRole) error {
	workspace, err := s.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return err
	}

	member, err := s.workspaceRepo.GetMember(ctx, workspaceID, userID)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}

	// If not a member, check if workspace is public and only viewer access is required
	if member == nil {
		if workspace.IsPublic && requiredRole == models.WorkspaceRoleViewer {
			return nil // Allow public view
		}
		return fmt.Errorf("access denied")
	}

	// Check role hierarchy: owner > editor > viewer
	if !hasPermission(member.Role, requiredRole) {
		return fmt.Errorf("insufficient permissions")
	}

	return nil
}

// IsOwner checks if user is the owner of workspace
func (s *WorkspaceService) IsOwner(ctx context.Context, workspaceID, userID uuid.UUID) (bool, error) {
	workspace, err := s.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return false, err
	}

	return workspace.OwnerID == userID, nil
}

// --- Helpers ---

func hasPermission(userRole, requiredRole models.WorkspaceRole) bool {
	roleHierarchy := map[models.WorkspaceRole]int{
		models.WorkspaceRoleViewer: 1,
		models.WorkspaceRoleEditor: 2,
		models.WorkspaceRoleOwner:  3,
	}

	return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}
