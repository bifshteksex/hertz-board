import { wsClient } from '$lib/services/websocket';
import { presenceStore, PresenceStore } from './presence.svelte';
import { crdtStore } from './crdt.svelte';
import { canvasStore } from './canvas.svelte';
import { authStore } from './auth.svelte';
import type {
	UserJoinedPayload,
	UserLeftPayload,
	PresenceUpdatePayload,
	OperationPayload,
	SyncResponsePayload,
	ErrorPayload
} from '$lib/types/websocket';
import type { CanvasElement } from '$lib/types/api';

// Forward declaration - will be set after canvas is created
let canvasWithHistoryRef: {
	applyRemoteOperation: (
		elementId: string,
		opType: 'create' | 'update' | 'delete' | 'move',
		data?: Partial<CanvasElement>
	) => void;
} | null = null;

/**
 * Collaboration Store
 * Orchestrates WebSocket, Presence, and CRDT stores
 * Handles real-time synchronization of canvas operations
 */
class CollaborationStore {
	private _isConnected = $state(false);
	private _currentWorkspaceId = $state<string | null>(null);
	private _userColor = $state<string>('#3B82F6');
	private _isSyncing = $state(false);
	private _error = $state<string | null>(null);

	// Throttle timers
	private cursorThrottle: ReturnType<typeof setTimeout> | null = null;
	private selectionThrottle: ReturnType<typeof setTimeout> | null = null;

	// Cleanup timer for inactive users
	private cleanupTimer: ReturnType<typeof setInterval> | null = null;

	// Unsubscribe functions
	private unsubscribers: (() => void)[] = [];

	constructor() {
		// Subscribe to WebSocket status changes
		wsClient.onStatusChange((status) => {
			this._isConnected = status.state === 'connected';
			if (status.state === 'error') {
				this._error = status.error || 'Connection error';
			}
		});
	}

	get isConnected(): boolean {
		return this._isConnected;
	}

	get workspaceId(): string | null {
		return this._currentWorkspaceId;
	}

	get userColor(): string {
		return this._userColor;
	}

	get isSyncing(): boolean {
		return this._isSyncing;
	}

	get error(): string | null {
		return this._error;
	}

	/**
	 * Set reference to canvasWithHistory for applying remote operations
	 * This avoids circular dependency issues
	 */
	setCanvasRef(canvas: {
		applyRemoteOperation: (
			elementId: string,
			opType: 'create' | 'update' | 'delete' | 'move',
			data?: Partial<CanvasElement>
		) => void;
	}) {
		canvasWithHistoryRef = canvas;
	}

	/**
	 * Connect to a workspace
	 */
	async connect(workspaceId: string, token: string): Promise<void> {
		try {
			this._error = null;
			this._currentWorkspaceId = workspaceId;

			// Generate user color
			this._userColor = PresenceStore.generateUserColor();

			// Initialize CRDT with user ID
			const userId = authStore.user?.id;
			if (!userId) {
				throw new Error('User not authenticated');
			}
			crdtStore.initialize(userId);

			// Connect to WebSocket
			if (!wsClient.isConnected) {
				await wsClient.connect(token);
			}

			// Register message handlers
			this.registerHandlers();

			// Join room
			wsClient.joinRoom(workspaceId, this._userColor);

			// Set workspace in presence store
			presenceStore.setWorkspace(workspaceId);

			// Request initial sync
			this.requestSync();

			// Start cleanup timer for inactive users
			this.startCleanupTimer();

			console.log('[Collaboration] Connected to workspace:', workspaceId);
		} catch (error) {
			console.error('[Collaboration] Connection failed:', error);
			this._error = error instanceof Error ? error.message : 'Connection failed';
			throw error;
		}
	}

	/**
	 * Disconnect from workspace
	 */
	disconnect(): void {
		if (this._currentWorkspaceId) {
			wsClient.leaveRoom(this._currentWorkspaceId);
		}

		// Stop cleanup timer
		this.stopCleanupTimer();

		// Unregister handlers
		this.unsubscribers.forEach((unsub) => unsub());
		this.unsubscribers = [];

		// Clear stores
		presenceStore.clear();
		crdtStore.reset();

		// Clear state
		this._currentWorkspaceId = null;
		this._isConnected = false;
		this._error = null;

		// Disconnect WebSocket if no other workspaces
		wsClient.disconnect();

		console.log('[Collaboration] Disconnected');
	}

	/**
	 * Send cursor position (throttled)
	 */
	sendCursorMove(x: number, y: number): void {
		if (!this._currentWorkspaceId) return;

		// Throttle cursor updates to max 10 per second
		if (this.cursorThrottle) return;

		wsClient.sendCursorMove(this._currentWorkspaceId, x, y);

		this.cursorThrottle = setTimeout(() => {
			this.cursorThrottle = null;
		}, 100);
	}

	/**
	 * Send selection change (throttled)
	 */
	sendSelectionChange(selectedIds: string[]): void {
		if (!this._currentWorkspaceId) return;

		// Throttle selection updates
		if (this.selectionThrottle) {
			clearTimeout(this.selectionThrottle);
		}

		this.selectionThrottle = setTimeout(() => {
			wsClient.sendSelectionChange(this._currentWorkspaceId!, selectedIds);
			this.selectionThrottle = null;
		}, 50);
	}

	/**
	 * Send operation to other users
	 */
	sendOperation(
		elementId: string,
		opType: 'create' | 'update' | 'delete' | 'move',
		data?: Record<string, unknown>
	): void {
		if (!this._currentWorkspaceId) return;

		try {
			const operation = crdtStore.createOperation(
				elementId,
				this._currentWorkspaceId,
				opType,
				data
			);

			wsClient.sendOperation(operation);
		} catch (error) {
			console.error('[Collaboration] Failed to send operation:', error);
		}
	}

	/**
	 * Request full sync
	 */
	requestSync(): void {
		if (!this._currentWorkspaceId) return;

		this._isSyncing = true;
		const stateVector = crdtStore.stateVector;
		wsClient.requestSync(this._currentWorkspaceId, stateVector);
	}

	// Private methods

	private registerHandlers(): void {
		// User joined
		this.unsubscribers.push(
			wsClient.on<UserJoinedPayload>('user_joined', (payload) => {
				// Validate payload
				if (!payload || !payload.user_id || !payload.user_name) {
					console.warn('[Collaboration] Invalid user_joined payload:', payload);
					return;
				}

				// Ignore self - don't add own cursor to presence store
				const currentUserId = authStore.user?.id;
				if (currentUserId && payload.user_id === currentUserId) {
					console.log('[Collaboration] Ignoring user_joined for self:', payload.user_name);
					return;
				}

				console.log('[Collaboration] User joined:', payload.user_name, 'ID:', payload.user_id);
				presenceStore.updateUser({
					user_id: payload.user_id,
					user_name: payload.user_name,
					user_color: payload.user_color,
					cursor: payload.cursor,
					selected_elements: [],
					last_seen: new Date().toISOString()
				});
			})
		);

		// User left
		this.unsubscribers.push(
			wsClient.on<UserLeftPayload>('user_left', (payload) => {
				console.log('[Collaboration] User left:', payload.user_id);
				presenceStore.removeUser(payload.user_id);
			})
		);

		// Presence update
		this.unsubscribers.push(
			wsClient.on<PresenceUpdatePayload>('presence_update', (payload) => {
				// Validate payload
				if (!payload || !payload.presence || !payload.presence.user_id) {
					console.warn('[Collaboration] Invalid presence_update payload:', payload);
					return;
				}

				// Ignore self - don't add own cursor to presence store
				const currentUserId = authStore.user?.id;
				if (currentUserId && payload.presence.user_id === currentUserId) {
					// Silently ignore presence updates for self (happens frequently)
					return;
				}

				presenceStore.updateUser(payload.presence);
			})
		);

		// Operation (CRDT)
		this.unsubscribers.push(
			wsClient.on<OperationPayload>('operation', (payload) => {
				this.handleOperation(payload);
			})
		);

		// Sync response
		this.unsubscribers.push(
			wsClient.on<SyncResponsePayload>('sync_response', (payload) => {
				this.handleSyncResponse(payload);
			})
		);

		// Error
		this.unsubscribers.push(
			wsClient.on<ErrorPayload>('error', (payload) => {
				console.error('[Collaboration] Server error:', payload);
				this._error = payload.message;
			})
		);

		// Pong (heartbeat response)
		this.unsubscribers.push(
			wsClient.on('pong', () => {
				// Heartbeat acknowledged
			})
		);
	}

	private handleOperation(operation: OperationPayload): void {
		try {
			// Ignore operations from self to prevent feedback loop
			const currentUserId = authStore.user?.id;
			if (currentUserId && operation.user_id === currentUserId) {
				console.log('[Collaboration] Ignoring own operation:', operation.element_id);
				return;
			}

			// Get current element
			const currentElement = canvasStore.elements.find((el) => el.id === operation.element_id);

			// Apply operation through CRDT
			const result = crdtStore.applyOperation(operation, currentElement);

			// Use canvasWithHistory to apply operations (if reference is set)
			if (canvasWithHistoryRef) {
				if (result === null && operation.op_type === 'delete') {
					// Delete element
					canvasWithHistoryRef.applyRemoteOperation(operation.element_id, 'delete');
					console.log('[Collaboration] Applied delete operation:', operation.element_id);
				} else if (result) {
					// Create or update element
					if (currentElement) {
						// Update existing
						canvasWithHistoryRef.applyRemoteOperation(
							operation.element_id,
							operation.op_type as 'update' | 'move',
							result
						);
						console.log('[Collaboration] Applied update operation:', operation.element_id);
					} else {
						// Create new
						canvasWithHistoryRef.applyRemoteOperation(operation.element_id, 'create', result);
						console.log('[Collaboration] Applied create operation:', operation.element_id);
					}
				}
			} else {
				// Fallback to direct canvas store manipulation (old behavior)
				if (result === null && operation.op_type === 'delete') {
					const elements = canvasStore.elements.filter((el) => el.id !== operation.element_id);
					canvasStore.setElements(elements);
				} else if (result) {
					if (currentElement) {
						const elements = canvasStore.elements.map((el) => (el.id === result.id ? result : el));
						canvasStore.setElements(elements);
					} else {
						canvasStore.addElements([result]);
					}
				}
			}
		} catch (error) {
			console.error('[Collaboration] Failed to apply operation:', error);
		}
	}

	private handleSyncResponse(payload: SyncResponsePayload): void {
		try {
			console.log('[Collaboration] Sync response:', payload.operations.length, 'operations');

			// Update state vector
			crdtStore.setStateVector(payload.state_vector);

			// Apply operations in order
			payload.operations.forEach((operation) => {
				this.handleOperation(operation);
			});

			this._isSyncing = false;
			console.log('[Collaboration] Sync completed');
		} catch (error) {
			console.error('[Collaboration] Failed to sync:', error);
			this._isSyncing = false;
			this._error = 'Sync failed';
		}
	}

	/**
	 * Start cleanup timer for removing inactive users
	 */
	private startCleanupTimer(): void {
		this.stopCleanupTimer();

		// Run cleanup every 10 seconds
		this.cleanupTimer = setInterval(() => {
			// Remove users inactive for more than 30 seconds
			presenceStore.removeInactiveUsers(30000);
		}, 10000);

		console.log('[Collaboration] Started inactive user cleanup timer');
	}

	/**
	 * Stop cleanup timer
	 */
	private stopCleanupTimer(): void {
		if (this.cleanupTimer) {
			clearInterval(this.cleanupTimer);
			this.cleanupTimer = null;
			console.log('[Collaboration] Stopped inactive user cleanup timer');
		}
	}
}

export const collaborationStore = new CollaborationStore();
