import { api } from '$lib/services/api';
import { authStore } from '$lib/stores/auth.svelte';

export interface OAuthCallbackParams {
	access_token?: string;
	refresh_token?: string;
	error?: string;
}

/**
 * Parse OAuth callback URL and extract tokens
 */
export function parseOAuthCallback(url: string): OAuthCallbackParams {
	try {
		const urlObj = new URL(url);
		const params: OAuthCallbackParams = {};

		// Extract access_token
		const accessToken = urlObj.searchParams.get('access_token');
		if (accessToken) {
			params.access_token = accessToken;
		}

		// Extract refresh_token
		const refreshToken = urlObj.searchParams.get('refresh_token');
		if (refreshToken) {
			params.refresh_token = refreshToken;
		}

		// Extract error if present
		const error = urlObj.searchParams.get('error');
		if (error) {
			params.error = error;
		}

		return params;
	} catch (e) {
		console.error('Failed to parse OAuth callback URL:', e);
		return { error: 'Invalid callback URL' };
	}
}

/**
 * Handle OAuth callback and authenticate user
 */
export async function handleOAuthCallback(
	url: string
): Promise<{ success: boolean; error?: string }> {
	const params = parseOAuthCallback(url);

	if (params.error) {
		return { success: false, error: params.error };
	}

	if (!params.access_token || !params.refresh_token) {
		return { success: false, error: 'Missing authentication tokens' };
	}

	try {
		// Set tokens in API client
		api.setTokens(params.access_token, params.refresh_token);

		// Fetch user data and update auth store
		const user = await api.getCurrentUser();
		authStore.setUser(user);

		return { success: true };
	} catch (error) {
		console.error('OAuth authentication failed:', error);
		return {
			success: false,
			error: error instanceof Error ? error.message : 'Authentication failed'
		};
	}
}

/**
 * Setup OAuth callback listener for Tauri deep links
 */
export function setupOAuthCallbackListener(
	onSuccess: () => void,
	onError: (error: string) => void
) {
	if (typeof window === 'undefined') return;

	const handler = async (event: Event) => {
		const customEvent = event as CustomEvent<string>;
		const url = customEvent.detail;
		console.log('OAuth callback received:', url);

		const result = await handleOAuthCallback(url);

		if (result.success) {
			onSuccess();
		} else {
			onError(result.error || 'Unknown error');
		}
	};

	window.addEventListener('oauth-callback', handler);

	// Return cleanup function
	return () => {
		window.removeEventListener('oauth-callback', handler);
	};
}
