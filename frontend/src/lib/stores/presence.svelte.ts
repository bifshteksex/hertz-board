import type { UserPresence, CursorPosition } from '$lib/types/websocket';

/**
 * Store for tracking user presence in a workspace
 * Manages:
 * - User cursors
 * - User selections
 * - Active users list
 * - User colors
 */
class PresenceStore {
	private _users = $state<Map<string, UserPresence>>(new Map());
	private _currentWorkspaceId = $state<string | null>(null);
	private _version = $state(0); // Trigger reactivity on changes

	/**
	 * Get all active users in the current workspace
	 */
	get users(): UserPresence[] {
		// Access _version to make this reactive
		void this._version;
		return Array.from(this._users.values());
	}

	/**
	 * Get number of active users
	 */
	get userCount(): number {
		return this._users.size;
	}

	/**
	 * Get current workspace ID
	 */
	get workspaceId(): string | null {
		return this._currentWorkspaceId;
	}

	/**
	 * Get a specific user by ID
	 */
	getUser(userId: string): UserPresence | undefined {
		return this._users.get(userId);
	}

	/**
	 * Get cursor position for a user
	 */
	getUserCursor(userId: string): CursorPosition | undefined {
		return this._users.get(userId)?.cursor;
	}

	/**
	 * Get selected elements for a user
	 */
	getUserSelection(userId: string): string[] {
		return this._users.get(userId)?.selected_elements || [];
	}

	/**
	 * Set the current workspace
	 */
	setWorkspace(workspaceId: string): void {
		if (this._currentWorkspaceId !== workspaceId) {
			// Clear users when switching workspaces
			this._users.clear();
			this._currentWorkspaceId = workspaceId;
			this._version++; // Trigger reactivity
		}
	}

	/**
	 * Add or update a user's presence
	 */
	updateUser(presence: UserPresence): void {
		const existing = this._users.get(presence.user_id);

		if (existing) {
			// Merge with existing data
			this._users.set(presence.user_id, {
				...existing,
				...presence,
				// Update cursor only if provided
				cursor: presence.cursor !== undefined ? presence.cursor : existing.cursor,
				// Update selection only if provided
				selected_elements:
					presence.selected_elements !== undefined
						? presence.selected_elements
						: existing.selected_elements,
				last_seen: presence.last_seen || new Date().toISOString()
			});
		} else {
			// Add new user
			this._users.set(presence.user_id, {
				...presence,
				selected_elements: presence.selected_elements || [],
				last_seen: presence.last_seen || new Date().toISOString()
			});
		}
		this._version++; // Trigger reactivity
	}

	/**
	 * Update a user's cursor position
	 */
	updateCursor(userId: string, cursor: CursorPosition): void {
		const user = this._users.get(userId);
		if (user) {
			this._users.set(userId, {
				...user,
				cursor,
				last_seen: new Date().toISOString()
			});
		}
	}

	/**
	 * Update a user's selection
	 */
	updateSelection(userId: string, selectedElements: string[]): void {
		const user = this._users.get(userId);
		if (user) {
			this._users.set(userId, {
				...user,
				selected_elements: selectedElements,
				last_seen: new Date().toISOString()
			});
		}
	}

	/**
	 * Remove a user from the workspace
	 */
	removeUser(userId: string): void {
		this._users.delete(userId);
		this._version++; // Trigger reactivity
	}

	/**
	 * Check if an element is selected by any user
	 */
	isElementSelected(elementId: string): boolean {
		return this.users.some((user) => user.selected_elements.includes(elementId));
	}

	/**
	 * Get all users who have selected a specific element
	 */
	getUsersSelectingElement(elementId: string): UserPresence[] {
		return this.users.filter((user) => user.selected_elements.includes(elementId));
	}

	/**
	 * Remove inactive users (no activity for more than N seconds)
	 * @param inactivityThreshold - milliseconds of inactivity before removal
	 * @param excludeUserId - user ID to exclude from removal (e.g., current user)
	 */
	removeInactiveUsers(inactivityThreshold = 60000, excludeUserId?: string): void {
		const now = Date.now();
		const toRemove: string[] = [];

		this._users.forEach((user, userId) => {
			// Don't remove excluded user (e.g., current user)
			if (excludeUserId && userId === excludeUserId) {
				return;
			}

			const lastSeen = new Date(user.last_seen).getTime();
			if (now - lastSeen > inactivityThreshold) {
				toRemove.push(userId);
			}
		});

		toRemove.forEach((userId) => {
			this._users.delete(userId);
		});

		if (toRemove.length > 0) {
			this._version++; // Trigger reactivity
			console.log(`[Presence] Removed ${toRemove.length} inactive users`);
		}
	}

	/**
	 * Clear all users (when leaving workspace)
	 */
	clear(): void {
		this._users.clear();
		this._currentWorkspaceId = null;
		this._version++; // Trigger reactivity
	}

	/**
	 * Generate a random color for a user
	 */
	static generateUserColor(): string {
		const colors = [
			'#FF6B6B', // Red
			'#4ECDC4', // Teal
			'#45B7D1', // Blue
			'#FFA07A', // Light Salmon
			'#98D8C8', // Mint
			'#F7DC6F', // Yellow
			'#BB8FCE', // Purple
			'#85C1E2', // Sky Blue
			'#F8B739', // Orange
			'#52B788', // Green
			'#E76F51', // Burnt Orange
			'#2A9D8F', // Teal Green
			'#E9C46A', // Sand
			'#F4A261', // Sandy Brown
			'#264653' // Dark Blue
		];

		return colors[Math.floor(Math.random() * colors.length)];
	}
}

export const presenceStore = new PresenceStore();
export { PresenceStore };
