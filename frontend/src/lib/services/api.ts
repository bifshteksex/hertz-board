import type {
	ApiError,
	AuthResponse,
	TokenPair,
	User,
	RegisterRequest,
	LoginRequest,
	RefreshTokenRequest,
	UpdateProfileRequest,
	ChangePasswordRequest,
	ForgotPasswordRequest,
	ResetPasswordRequest,
	Workspace,
	WorkspaceListResponse,
	WorkspaceFilters,
	CreateWorkspaceRequest,
	UpdateWorkspaceRequest,
	WorkspaceMember,
	InviteMemberRequest,
	UpdateMemberRoleRequest,
	WorkspaceInvitation,
	SetPublicAccessRequest,
	CanvasElement,
	CreateElementRequest,
	UpdateElementRequest,
	Asset,
	Snapshot,
	CreateSnapshotRequest
} from '$lib/types/api';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export class ApiClient {
	private accessToken: string | null = null;
	private refreshToken: string | null = null;

	constructor() {
		// Load tokens from localStorage on initialization
		if (typeof window !== 'undefined') {
			this.accessToken = localStorage.getItem('access_token');
			this.refreshToken = localStorage.getItem('refresh_token');
		}
	}

	setTokens(accessToken: string, refreshToken: string) {
		this.accessToken = accessToken;
		this.refreshToken = refreshToken;
		if (typeof window !== 'undefined') {
			localStorage.setItem('access_token', accessToken);
			localStorage.setItem('refresh_token', refreshToken);
		}
	}

	clearTokens() {
		this.accessToken = null;
		this.refreshToken = null;
		if (typeof window !== 'undefined') {
			localStorage.removeItem('access_token');
			localStorage.removeItem('refresh_token');
		}
	}

	getAccessToken(): string | null {
		return this.accessToken;
	}

	private async request<T>(
		endpoint: string,
		options: globalThis.RequestInit = {},
		retry = true
	): Promise<T> {
		const url = `${API_BASE_URL}${endpoint}`;
		const headers: globalThis.HeadersInit = {
			'Content-Type': 'application/json',
			...options.headers
		};

		// Add authorization header if access token exists
		if (this.accessToken && !endpoint.startsWith('/auth/')) {
			headers['Authorization'] = `Bearer ${this.accessToken}`;
		}

		try {
			const response = await fetch(url, {
				...options,
				headers
			});

			// Handle 401 Unauthorized
			if (response.status === 401 && retry) {
				// Try to refresh token
				const refreshed = await this.tryRefreshToken();
				if (refreshed) {
					// Retry the request with new token
					return this.request<T>(endpoint, options, false);
				} else {
					// Clear tokens but don't redirect - let the caller handle it
					this.clearTokens();
					throw new Error('Unauthorized');
				}
			}

			// Handle other error status codes
			if (!response.ok) {
				const error: ApiError = await response.json().catch(() => ({
					error: 'Unknown error occurred'
				}));
				throw new Error(error.error || `HTTP ${response.status}: ${response.statusText}`);
			}

			// Handle 204 No Content
			if (response.status === 204) {
				return {} as T;
			}

			return response.json();
		} catch (error) {
			console.error('API request failed:', error);
			throw error;
		}
	}

	private async tryRefreshToken(): Promise<boolean> {
		if (!this.refreshToken) {
			return false;
		}

		try {
			const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ refresh_token: this.refreshToken } as RefreshTokenRequest)
			});

			if (response.ok) {
				const data: TokenPair = await response.json();
				this.setTokens(data.access_token, data.refresh_token);
				return true;
			}

			return false;
		} catch (error) {
			console.error('Token refresh failed:', error);
			return false;
		}
	}

	// Auth endpoints
	async register(data: RegisterRequest): Promise<AuthResponse> {
		const response = await this.request<AuthResponse>('/auth/register', {
			method: 'POST',
			body: JSON.stringify(data)
		});
		this.setTokens(response.tokens.access_token, response.tokens.refresh_token);
		return response;
	}

	async login(data: LoginRequest): Promise<AuthResponse> {
		const response = await this.request<AuthResponse>('/auth/login', {
			method: 'POST',
			body: JSON.stringify(data)
		});
		this.setTokens(response.tokens.access_token, response.tokens.refresh_token);
		return response;
	}

	async logout(): Promise<void> {
		await this.request('/auth/logout', {
			method: 'POST',
			body: JSON.stringify({ refresh_token: this.refreshToken })
		});
		this.clearTokens();
	}

	async forgotPassword(data: ForgotPasswordRequest): Promise<void> {
		return this.request('/auth/forgot-password', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async resetPassword(data: ResetPasswordRequest): Promise<void> {
		return this.request('/auth/reset-password', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	// User endpoints
	async getCurrentUser(): Promise<User> {
		return this.request<User>('/users/me');
	}

	async updateProfile(data: UpdateProfileRequest): Promise<User> {
		return this.request<User>('/users/me', {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	}

	async changePassword(data: ChangePasswordRequest): Promise<void> {
		return this.request('/users/me/password', {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	}

	// Workspace endpoints
	async listWorkspaces(filters?: WorkspaceFilters): Promise<WorkspaceListResponse> {
		const params = new URLSearchParams();
		if (filters) {
			Object.entries(filters).forEach(([key, value]) => {
				if (value !== undefined && value !== null) {
					params.append(key, String(value));
				}
			});
		}
		const query = params.toString();
		return this.request<WorkspaceListResponse>(`/workspaces${query ? `?${query}` : ''}`);
	}

	async getWorkspace(id: string): Promise<Workspace> {
		const response = await this.request<{ workspace: Workspace }>(`/workspaces/${id}`);
		return response.workspace;
	}

	async createWorkspace(data: CreateWorkspaceRequest): Promise<Workspace> {
		const response = await this.request<{ workspace: Workspace }>('/workspaces', {
			method: 'POST',
			body: JSON.stringify(data)
		});
		return response.workspace;
	}

	async updateWorkspace(id: string, data: UpdateWorkspaceRequest): Promise<Workspace> {
		const response = await this.request<{ workspace: Workspace }>(`/workspaces/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
		return response.workspace;
	}

	async deleteWorkspace(id: string): Promise<void> {
		return this.request(`/workspaces/${id}`, {
			method: 'DELETE'
		});
	}

	async duplicateWorkspace(id: string, name: string): Promise<Workspace> {
		const response = await this.request<{ workspace: Workspace }>(`/workspaces/${id}/duplicate`, {
			method: 'POST',
			body: JSON.stringify({ name })
		});
		return response.workspace;
	}

	// Workspace members
	async listMembers(workspaceId: string): Promise<WorkspaceMember[]> {
		return this.request<WorkspaceMember[]>(`/workspaces/${workspaceId}/members`);
	}

	async inviteMember(workspaceId: string, data: InviteMemberRequest): Promise<void> {
		return this.request(`/workspaces/${workspaceId}/members/invite`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async updateMemberRole(
		workspaceId: string,
		memberId: string,
		data: UpdateMemberRoleRequest
	): Promise<void> {
		return this.request(`/workspaces/${workspaceId}/members/${memberId}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	}

	async removeMember(workspaceId: string, memberId: string): Promise<void> {
		return this.request(`/workspaces/${workspaceId}/members/${memberId}`, {
			method: 'DELETE'
		});
	}

	// Workspace invitations
	async listInvitations(): Promise<WorkspaceInvitation[]> {
		return this.request<WorkspaceInvitation[]>('/invitations');
	}

	async acceptInvitation(invitationId: string): Promise<void> {
		return this.request(`/invitations/${invitationId}/accept`, {
			method: 'POST'
		});
	}

	async declineInvitation(invitationId: string): Promise<void> {
		return this.request(`/invitations/${invitationId}/decline`, {
			method: 'POST'
		});
	}

	// Public access
	async setPublicAccess(workspaceId: string, data: SetPublicAccessRequest): Promise<Workspace> {
		return this.request<Workspace>(`/workspaces/${workspaceId}/public`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	}

	async regeneratePublicToken(workspaceId: string): Promise<Workspace> {
		return this.request<Workspace>(`/workspaces/${workspaceId}/public/regenerate`, {
			method: 'POST'
		});
	}

	// Canvas elements
	async listElements(workspaceId: string): Promise<CanvasElement[]> {
		return this.request<CanvasElement[]>(`/workspaces/${workspaceId}/elements`);
	}

	async getElement(workspaceId: string, elementId: string): Promise<CanvasElement> {
		return this.request<CanvasElement>(`/workspaces/${workspaceId}/elements/${elementId}`);
	}

	async createElement(workspaceId: string, data: CreateElementRequest): Promise<CanvasElement> {
		return this.request<CanvasElement>(`/workspaces/${workspaceId}/elements`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async updateElement(
		workspaceId: string,
		elementId: string,
		data: UpdateElementRequest
	): Promise<CanvasElement> {
		return this.request<CanvasElement>(`/workspaces/${workspaceId}/elements/${elementId}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	}

	async deleteElement(workspaceId: string, elementId: string): Promise<void> {
		return this.request(`/workspaces/${workspaceId}/elements/${elementId}`, {
			method: 'DELETE'
		});
	}

	async batchCreateElements(
		workspaceId: string,
		elements: CreateElementRequest[]
	): Promise<CanvasElement[]> {
		return this.request<CanvasElement[]>(`/workspaces/${workspaceId}/elements/batch`, {
			method: 'POST',
			body: JSON.stringify({ elements })
		});
	}

	// Assets
	async uploadAsset(workspaceId: string, file: File): Promise<Asset> {
		const formData = new FormData();
		formData.append('file', file);

		const url = `${API_BASE_URL}/workspaces/${workspaceId}/assets`;
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${this.accessToken}`
			},
			body: formData
		});

		if (!response.ok) {
			throw new Error('Upload failed');
		}

		return response.json();
	}

	async listAssets(workspaceId: string): Promise<Asset[]> {
		return this.request<Asset[]>(`/workspaces/${workspaceId}/assets`);
	}

	async deleteAsset(workspaceId: string, assetId: string): Promise<void> {
		return this.request(`/workspaces/${workspaceId}/assets/${assetId}`, {
			method: 'DELETE'
		});
	}

	// Snapshots
	async listSnapshots(workspaceId: string): Promise<Snapshot[]> {
		return this.request<Snapshot[]>(`/workspaces/${workspaceId}/snapshots`);
	}

	async createSnapshot(workspaceId: string, data: CreateSnapshotRequest): Promise<Snapshot> {
		return this.request<Snapshot>(`/workspaces/${workspaceId}/snapshots`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async restoreSnapshot(workspaceId: string, snapshotId: string): Promise<void> {
		return this.request(`/workspaces/${workspaceId}/snapshots/${snapshotId}/restore`, {
			method: 'POST'
		});
	}

	async deleteSnapshot(workspaceId: string, snapshotId: string): Promise<void> {
		return this.request(`/workspaces/${workspaceId}/snapshots/${snapshotId}`, {
			method: 'DELETE'
		});
	}
}

// Singleton instance
export const api = new ApiClient();
