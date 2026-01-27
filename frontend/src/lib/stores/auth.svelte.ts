import { api } from '$lib/services/api';
import type { User } from '$lib/types/api';

class AuthStore {
	private _user = $state<User | null>(null);
	private _isLoading = $state(true);
	private _isInitialized = $state(false);

	get user() {
		return this._user;
	}

	get isAuthenticated() {
		return this._user !== null;
	}

	get isLoading() {
		return this._isLoading;
	}

	get isInitialized() {
		return this._isInitialized;
	}

	async initialize() {
		if (this._isInitialized) {
			return;
		}

		this._isLoading = true;

		try {
			// Check if we have an access token
			const token = api.getAccessToken();
			if (token) {
				// Try to get current user
				const user = await api.getCurrentUser();
				this._user = user;
			}
		} catch (error) {
			console.error('Failed to initialize auth:', error);
			// Clear invalid tokens
			api.clearTokens();
			this._user = null;
		} finally {
			this._isLoading = false;
			this._isInitialized = true;
		}
	}

	async login(email: string, password: string) {
		this._isLoading = true;
		try {
			const response = await api.login({ email, password });
			this._user = response.user;
			return response;
		} finally {
			this._isLoading = false;
		}
	}

	async register(email: string, password: string, name: string) {
		this._isLoading = true;
		try {
			const response = await api.register({ email, password, name });
			this._user = response.user;
			return response;
		} finally {
			this._isLoading = false;
		}
	}

	async logout() {
		this._isLoading = true;
		try {
			await api.logout();
		} catch (error) {
			console.error('Logout failed:', error);
		} finally {
			this._user = null;
			this._isLoading = false;
		}
	}

	async updateProfile(data: { name?: string; avatar_url?: string }) {
		if (!this._user) {
			throw new Error('Not authenticated');
		}

		this._isLoading = true;
		try {
			const updatedUser = await api.updateProfile(data);
			this._user = updatedUser;
			return updatedUser;
		} finally {
			this._isLoading = false;
		}
	}

	async changePassword(currentPassword: string, newPassword: string) {
		this._isLoading = true;
		try {
			await api.changePassword({
				current_password: currentPassword,
				new_password: newPassword
			});
		} finally {
			this._isLoading = false;
		}
	}

	setUser(user: User | null) {
		this._user = user;
	}
}

export const authStore = new AuthStore();
