// API Response Types

export interface User {
	id: string;
	email: string;
	name: string;
	avatar_url?: string;
	provider: 'email' | 'google' | 'github';
	email_verified: boolean;
	created_at: string;
	updated_at: string;
}

export interface TokenPair {
	access_token: string;
	refresh_token: string;
}

export interface AuthResponse {
	user: User;
	tokens: TokenPair;
}

export interface Workspace {
	id: string;
	name: string;
	description?: string;
	thumbnail_url?: string;
	owner_id: string;
	is_public: boolean;
	public_access_token?: string;
	created_at: string;
	updated_at: string;
	// Extended fields from API
	owner_name?: string;
	user_role?: WorkspaceRole; // Backend returns user_role
	role?: WorkspaceRole; // Alias for compatibility
	member_count?: number;
}

export type WorkspaceRole = 'owner' | 'editor' | 'viewer';

export interface WorkspaceMember {
	id: string;
	workspace_id: string;
	user_id: string;
	role: WorkspaceRole;
	joined_at: string;
	user_name?: string;
	user_email?: string;
	user_avatar?: string;
}

export interface WorkspaceInvitation {
	id: string;
	workspace_id: string;
	workspace_name?: string;
	inviter_id: string;
	inviter_name?: string;
	invitee_email: string;
	role: WorkspaceRole;
	status: 'pending' | 'accepted' | 'declined' | 'expired';
	expires_at: string;
	created_at: string;
}

export interface WorkspaceListResponse {
	workspaces: Workspace[];
	total: number;
	limit: number;
	offset: number;
}

// Canvas Element Types
export type ElementType =
	| 'text'
	| 'rectangle'
	| 'ellipse'
	| 'triangle'
	| 'line'
	| 'arrow'
	| 'sticky'
	| 'image'
	| 'freehand'
	| 'list'
	| 'connector';

export interface ElementStyle {
	backgroundColor?: string;
	strokeColor?: string;
	strokeWidth?: number;
	fontFamily?: string;
	fontSize?: number;
	fontWeight?: string;
	textAlign?: 'left' | 'center' | 'right';
	opacity?: number;
	borderRadius?: number;
	color?: string;
	listType?: 'bullet' | 'numbered' | 'checkbox';
	connectorType?: 'straight' | 'curved' | 'elbow';
}

export interface ConnectorData {
	startElementId?: string;
	endElementId?: string;
	startX: number;
	startY: number;
	endX: number;
	endY: number;
	startArrow?: boolean;
	endArrow?: boolean;
	label?: string;
}

export interface CanvasElement {
	id: string;
	workspace_id: string;
	type: ElementType;
	content?: string;
	html_content?: string;
	pos_x: number;
	pos_y: number;
	width?: number;
	height?: number;
	rotation?: number;
	z_index: number;
	style?: ElementStyle;
	locked?: boolean;
	parent_id?: string;
	group_id?: string;
	image_url?: string;
	path_data?: string;
	connector_data?: ConnectorData;
	created_by?: string;
	created_at?: string;
	updated_at?: string;
	version?: number;
}

export interface Asset {
	id: string;
	workspace_id: string;
	filename: string;
	file_size: number;
	mime_type: string;
	url: string;
	thumbnail_url?: string;
	uploaded_by: string;
	created_at: string;
}

export interface Snapshot {
	id: string;
	workspace_id: string;
	name: string;
	description?: string;
	element_count: number;
	created_by: string;
	created_at: string;
}

// Request Types
export interface RegisterRequest {
	email: string;
	password: string;
	name: string;
}

export interface LoginRequest {
	email: string;
	password: string;
}

export interface RefreshTokenRequest {
	refresh_token: string;
}

export interface UpdateProfileRequest {
	name?: string;
	avatar_url?: string;
}

export interface ChangePasswordRequest {
	current_password: string;
	new_password: string;
}

export interface ForgotPasswordRequest {
	email: string;
}

export interface ResetPasswordRequest {
	token: string;
	password: string;
}

export interface CreateWorkspaceRequest {
	name: string;
	description?: string;
}

export interface UpdateWorkspaceRequest {
	name?: string;
	description?: string;
	thumbnail_url?: string;
}

export interface InviteMemberRequest {
	email: string;
	role: WorkspaceRole;
}

export interface UpdateMemberRoleRequest {
	role: WorkspaceRole;
}

export interface SetPublicAccessRequest {
	is_public: boolean;
}

export interface CreateElementRequest {
	type: ElementType;
	content: string;
	pos_x: number;
	pos_y: number;
	width: number;
	height: number;
	rotation?: number;
	style?: ElementStyle;
	parent_id?: string;
}

export interface UpdateElementRequest {
	content?: string;
	pos_x?: number;
	pos_y?: number;
	width?: number;
	height?: number;
	rotation?: number;
	z_index?: number;
	style?: ElementStyle;
	locked?: boolean;
	parent_id?: string;
}

export interface CreateSnapshotRequest {
	name: string;
	description?: string;
}

// Error Types
export interface ApiError {
	error: string;
	details?: string;
}

// Pagination
export interface PaginationParams {
	limit?: number;
	offset?: number;
}

// Workspace Filters
export interface WorkspaceFilters extends PaginationParams {
	query?: string;
	owned_only?: boolean;
	shared_only?: boolean;
	sort_by?: 'name' | 'created_at' | 'updated_at';
	sort_order?: 'asc' | 'desc';
}
