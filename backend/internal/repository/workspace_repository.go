package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bifshteksex/hertzboard/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkspaceRepository struct {
	db *pgxpool.Pool
}

func NewWorkspaceRepository(db *pgxpool.Pool) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

// --- Workspace CRUD ---

// CreateWorkspace creates a new workspace and adds creator as owner
func (r *WorkspaceRepository) CreateWorkspace(ctx context.Context, workspace *models.Workspace) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	settingsJSON, err := json.Marshal(workspace.Settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Create workspace
	query := `
		INSERT INTO workspaces (id, name, description, owner_id, thumbnail_url, is_public, settings)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`
	err = tx.QueryRow(ctx, query,
		workspace.ID,
		workspace.Name,
		workspace.Description,
		workspace.OwnerID,
		workspace.ThumbnailURL,
		workspace.IsPublic,
		settingsJSON,
	).Scan(&workspace.CreatedAt, &workspace.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert workspace: %w", err)
	}

	// Add creator as owner member
	memberQuery := `
		INSERT INTO workspace_members (workspace_id, user_id, role)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, memberQuery, workspace.ID, workspace.OwnerID, models.WorkspaceRoleOwner)
	if err != nil {
		return fmt.Errorf("failed to add owner as member: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetWorkspaceByID retrieves a workspace by ID (excluding soft-deleted)
func (r *WorkspaceRepository) GetWorkspaceByID(ctx context.Context, id uuid.UUID) (*models.Workspace, error) {
	query := `
		SELECT id, name, description, owner_id, thumbnail_url, is_public, settings, deleted_at, created_at, updated_at
		FROM workspaces
		WHERE id = $1 AND deleted_at IS NULL
	`

	var workspace models.Workspace
	var settingsJSON []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&workspace.ID,
		&workspace.Name,
		&workspace.Description,
		&workspace.OwnerID,
		&workspace.ThumbnailURL,
		&workspace.IsPublic,
		&settingsJSON,
		&workspace.DeletedAt,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}

	if err := json.Unmarshal(settingsJSON, &workspace.Settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	return &workspace, nil
}

// UpdateWorkspace updates workspace fields
func (r *WorkspaceRepository) UpdateWorkspace(ctx context.Context, workspace *models.Workspace) error {
	settingsJSON, err := json.Marshal(workspace.Settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	query := `
		UPDATE workspaces
		SET name = $1, description = $2, is_public = $3, thumbnail_url = $4, settings = $5
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING updated_at
	`

	err = r.db.QueryRow(ctx, query,
		workspace.Name,
		workspace.Description,
		workspace.IsPublic,
		workspace.ThumbnailURL,
		settingsJSON,
		workspace.ID,
	).Scan(&workspace.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("workspace not found")
		}
		return fmt.Errorf("failed to update workspace: %w", err)
	}

	return nil
}

// SoftDeleteWorkspace marks workspace as deleted
func (r *WorkspaceRepository) SoftDeleteWorkspace(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE workspaces
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete workspace: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("workspace not found")
	}

	return nil
}

// ListWorkspacesByUser retrieves workspaces accessible to user with filters
func (r *WorkspaceRepository) ListWorkspacesByUser(
	ctx context.Context,
	userID uuid.UUID,
	filter models.WorkspaceListFilter,
) ([]models.WorkspaceWithRole, int, error) {
	// Build query with filters
	query := `
		SELECT DISTINCT
			w.id, w.name, w.description, w.owner_id, w.thumbnail_url,
			w.is_public, w.settings, w.created_at, w.updated_at,
			wm.role,
			COUNT(*) OVER() as total_count
		FROM workspaces w
		INNER JOIN workspace_members wm ON w.id = wm.workspace_id
		WHERE w.deleted_at IS NULL
			AND wm.user_id = $1
	`

	args := []interface{}{userID}
	argCount := 1

	// Apply filters
	if filter.OwnedOnly {
		query += " AND w.owner_id = $1"
	} else if filter.SharedOnly {
		query += " AND w.owner_id != $1"
	}

	if filter.Query != "" {
		argCount++
		query += fmt.Sprintf(" AND w.name ILIKE $%d", argCount)
		args = append(args, "%"+filter.Query+"%")
	}

	// Sorting
	sortBy := "created_at"
	if filter.SortBy == "updated_at" || filter.SortBy == "name" {
		sortBy = filter.SortBy
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY w.%s %s", sortBy, sortOrder)

	// Pagination
	limit := 20
	if filter.Limit > 0 && filter.Limit <= 100 {
		limit = filter.Limit
	}

	offset := 0
	if filter.Offset > 0 {
		offset = filter.Offset
	}

	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, limit)

	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list workspaces: %w", err)
	}
	defer rows.Close()

	var workspaces []models.WorkspaceWithRole
	var totalCount int

	for rows.Next() {
		var ws models.WorkspaceWithRole
		var settingsJSON []byte

		err := rows.Scan(
			&ws.ID,
			&ws.Name,
			&ws.Description,
			&ws.OwnerID,
			&ws.ThumbnailURL,
			&ws.IsPublic,
			&settingsJSON,
			&ws.CreatedAt,
			&ws.UpdatedAt,
			&ws.UserRole,
			&totalCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan workspace: %w", err)
		}

		if err := json.Unmarshal(settingsJSON, &ws.Settings); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal settings: %w", err)
		}

		workspaces = append(workspaces, ws)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating workspaces: %w", err)
	}

	return workspaces, totalCount, nil
}

// --- Workspace Members ---

// AddMember adds a user to workspace with specified role
func (r *WorkspaceRepository) AddMember(ctx context.Context, member *models.WorkspaceMember) error {
	query := `
		INSERT INTO workspace_members (id, workspace_id, user_id, role, invited_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING joined_at
		ON CONFLICT (workspace_id, user_id) DO NOTHING
	`

	err := r.db.QueryRow(ctx, query,
		member.ID,
		member.WorkspaceID,
		member.UserID,
		member.Role,
		member.InvitedBy,
	).Scan(&member.JoinedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("member already exists")
		}
		return fmt.Errorf("failed to add member: %w", err)
	}

	return nil
}

// GetMember retrieves member information
func (r *WorkspaceRepository) GetMember(ctx context.Context, workspaceID, userID uuid.UUID) (*models.WorkspaceMember, error) {
	query := `
		SELECT id, workspace_id, user_id, role, invited_by, joined_at
		FROM workspace_members
		WHERE workspace_id = $1 AND user_id = $2
	`

	var member models.WorkspaceMember
	err := r.db.QueryRow(ctx, query, workspaceID, userID).Scan(
		&member.ID,
		&member.WorkspaceID,
		&member.UserID,
		&member.Role,
		&member.InvitedBy,
		&member.JoinedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get member: %w", err)
	}

	return &member, nil
}

// UpdateMemberRole updates member's role in workspace
func (r *WorkspaceRepository) UpdateMemberRole(ctx context.Context, workspaceID, userID uuid.UUID, role models.WorkspaceRole) error {
	query := `
		UPDATE workspace_members
		SET role = $1
		WHERE workspace_id = $2 AND user_id = $3
	`

	result, err := r.db.Exec(ctx, query, role, workspaceID, userID)
	if err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

// RemoveMember removes a user from workspace
func (r *WorkspaceRepository) RemoveMember(ctx context.Context, workspaceID, userID uuid.UUID) error {
	query := `
		DELETE FROM workspace_members
		WHERE workspace_id = $1 AND user_id = $2
	`

	result, err := r.db.Exec(ctx, query, workspaceID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

// ListMembers retrieves all members of a workspace
func (r *WorkspaceRepository) ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceMemberWithUser, error) {
	query := `
		SELECT
			wm.id, wm.workspace_id, wm.user_id, wm.role, wm.invited_by, wm.joined_at,
			u.id, u.email, u.name, u.avatar_url
		FROM workspace_members wm
		INNER JOIN users u ON wm.user_id = u.id
		WHERE wm.workspace_id = $1
		ORDER BY wm.joined_at ASC
	`

	rows, err := r.db.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list members: %w", err)
	}
	defer rows.Close()

	var members []models.WorkspaceMemberWithUser

	for rows.Next() {
		var m models.WorkspaceMemberWithUser

		err := rows.Scan(
			&m.ID,
			&m.WorkspaceID,
			&m.UserID,
			&m.Role,
			&m.InvitedBy,
			&m.JoinedAt,
			&m.User.ID,
			&m.User.Email,
			&m.User.Name,
			&m.User.AvatarURL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}

		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating members: %w", err)
	}

	return members, nil
}

// --- Workspace Invites ---

// CreateInvite creates a new workspace invitation
func (r *WorkspaceRepository) CreateInvite(ctx context.Context, invite *models.WorkspaceInvite) error {
	query := `
		INSERT INTO workspace_invites (id, workspace_id, email, role, token_hash, expires_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`

	err := r.db.QueryRow(ctx, query,
		invite.ID,
		invite.WorkspaceID,
		invite.Email,
		invite.Role,
		invite.TokenHash,
		invite.ExpiresAt,
		invite.CreatedBy,
	).Scan(&invite.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create invite: %w", err)
	}

	return nil
}

// GetInviteByToken retrieves an invite by token hash
func (r *WorkspaceRepository) GetInviteByToken(ctx context.Context, tokenHash string) (*models.WorkspaceInvite, error) {
	query := `
		SELECT id, workspace_id, email, role, token_hash, expires_at, created_by, created_at, accepted_at, accepted_by
		FROM workspace_invites
		WHERE token_hash = $1
	`

	var invite models.WorkspaceInvite
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&invite.ID,
		&invite.WorkspaceID,
		&invite.Email,
		&invite.Role,
		&invite.TokenHash,
		&invite.ExpiresAt,
		&invite.CreatedBy,
		&invite.CreatedAt,
		&invite.AcceptedAt,
		&invite.AcceptedBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get invite: %w", err)
	}

	return &invite, nil
}

// MarkInviteAsAccepted marks an invitation as accepted
func (r *WorkspaceRepository) MarkInviteAsAccepted(ctx context.Context, inviteID, userID uuid.UUID) error {
	query := `
		UPDATE workspace_invites
		SET accepted_at = CURRENT_TIMESTAMP, accepted_by = $1
		WHERE id = $2 AND accepted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, userID, inviteID)
	if err != nil {
		return fmt.Errorf("failed to mark invite as accepted: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("invite not found or already accepted")
	}

	return nil
}

// RevokeInvite deletes an invitation
func (r *WorkspaceRepository) RevokeInvite(ctx context.Context, inviteID uuid.UUID) error {
	query := `
		DELETE FROM workspace_invites
		WHERE id = $1 AND accepted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, inviteID)
	if err != nil {
		return fmt.Errorf("failed to revoke invite: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("invite not found or already accepted")
	}

	return nil
}

// ListPendingInvites retrieves all pending invitations for a workspace
func (r *WorkspaceRepository) ListPendingInvites(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceInvite, error) {
	query := `
		SELECT id, workspace_id, email, role, token_hash, expires_at, created_by, created_at, accepted_at, accepted_by
		FROM workspace_invites
		WHERE workspace_id = $1 AND accepted_at IS NULL AND expires_at > CURRENT_TIMESTAMP
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list pending invites: %w", err)
	}
	defer rows.Close()

	var invites []models.WorkspaceInvite

	for rows.Next() {
		var invite models.WorkspaceInvite

		err := rows.Scan(
			&invite.ID,
			&invite.WorkspaceID,
			&invite.Email,
			&invite.Role,
			&invite.TokenHash,
			&invite.ExpiresAt,
			&invite.CreatedBy,
			&invite.CreatedAt,
			&invite.AcceptedAt,
			&invite.AcceptedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invite: %w", err)
		}

		invites = append(invites, invite)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating invites: %w", err)
	}

	return invites, nil
}

// CleanupExpiredInvites removes expired invitations
func (r *WorkspaceRepository) CleanupExpiredInvites(ctx context.Context) error {
	query := `
		DELETE FROM workspace_invites
		WHERE expires_at < CURRENT_TIMESTAMP AND accepted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired invites: %w", err)
	}

	return nil
}

// GetInviteByWorkspaceAndEmail checks if there's a pending invite for email in workspace
func (r *WorkspaceRepository) GetInviteByWorkspaceAndEmail(
	ctx context.Context,
	workspaceID uuid.UUID,
	email string,
) (*models.WorkspaceInvite, error) {
	query := `
		SELECT id, workspace_id, email, role, token_hash, expires_at, created_by, created_at, accepted_at, accepted_by
		FROM workspace_invites
		WHERE workspace_id = $1 AND email = $2 AND accepted_at IS NULL AND expires_at > CURRENT_TIMESTAMP
		ORDER BY created_at DESC
		LIMIT 1
	`

	var invite models.WorkspaceInvite
	err := r.db.QueryRow(ctx, query, workspaceID, email).Scan(
		&invite.ID,
		&invite.WorkspaceID,
		&invite.Email,
		&invite.Role,
		&invite.TokenHash,
		&invite.ExpiresAt,
		&invite.CreatedBy,
		&invite.CreatedAt,
		&invite.AcceptedAt,
		&invite.AcceptedBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get invite: %w", err)
	}

	return &invite, nil
}
