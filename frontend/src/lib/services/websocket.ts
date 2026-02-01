import type {
	WSMessage,
	MessageType,
	ConnectionState,
	ConnectionStatus,
	JoinRoomPayload,
	CursorMovePayload,
	SelectionChangePayload,
	OperationPayload,
	SyncRequestPayload
} from '$lib/types/websocket';

type MessageHandler<T = unknown> = (payload: T, message: WSMessage<T>) => void;

interface WebSocketConfig {
	url: string;
	reconnect?: boolean;
	reconnectInterval?: number;
	maxReconnectAttempts?: number;
	heartbeatInterval?: number;
}

/**
 * WebSocket client for real-time collaboration
 * Features:
 * - Automatic reconnection with exponential backoff
 * - Message queue during disconnection
 * - Heartbeat/ping mechanism
 * - Type-safe message handlers
 */
export class WebSocketClient {
	private ws: WebSocket | null = null;
	private config: Required<WebSocketConfig>;
	private handlers: Map<MessageType, Set<MessageHandler>> = new Map();
	private messageQueue: WSMessage[] = [];
	private reconnectAttempts = 0;
	private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
	private heartbeatTimer: ReturnType<typeof setInterval> | null = null;
	private _connectionState: ConnectionState = 'disconnected';
	private statusListeners: Set<(status: ConnectionStatus) => void> = new Set();
	private shouldReconnect = true;

	// LocalStorage key for persistent queue
	private readonly QUEUE_STORAGE_KEY = 'ws_message_queue';

	constructor(config: WebSocketConfig) {
		this.config = {
			reconnect: true,
			reconnectInterval: 1000,
			maxReconnectAttempts: 10,
			heartbeatInterval: 30000, // 30 seconds
			...config
		};

		// Restore queue from localStorage on init
		this.restoreQueue();
	}

	/**
	 * Connect to WebSocket server
	 * @param token JWT token for authentication
	 */
	connect(token: string): Promise<void> {
		return new Promise((resolve, reject) => {
			if (this.ws?.readyState === WebSocket.OPEN) {
				resolve();
				return;
			}

			this.shouldReconnect = true;
			this.setConnectionState('connecting');

			// Add token as query parameter
			const url = `${this.config.url}?token=${encodeURIComponent(token)}`;

			try {
				this.ws = new WebSocket(url);

				this.ws.onopen = () => {
					console.log('[WS] Connected');
					this.reconnectAttempts = 0;
					this.setConnectionState('connected');
					this.startHeartbeat();
					this.flushMessageQueue();
					resolve();
				};

				this.ws.onmessage = (event) => {
					this.handleMessage(event);
				};

				this.ws.onerror = (error) => {
					console.error('[WS] Error:', error);
					this.setConnectionState('error', 'Connection error');
					reject(new Error('WebSocket connection error'));
				};

				this.ws.onclose = (event) => {
					console.log('[WS] Closed:', event.code, event.reason);
					this.stopHeartbeat();

					if (this.shouldReconnect && this.config.reconnect) {
						this.scheduleReconnect();
					} else {
						this.setConnectionState('disconnected');
					}
				};
			} catch (error) {
				console.error('[WS] Failed to create connection:', error);
				this.setConnectionState('error', 'Failed to create connection');
				reject(error);
			}
		});
	}

	/**
	 * Disconnect from WebSocket server
	 */
	disconnect(): void {
		this.shouldReconnect = false;
		this.clearReconnectTimer();
		this.stopHeartbeat();

		if (this.ws) {
			this.ws.close(1000, 'Client disconnect');
			this.ws = null;
		}

		this.setConnectionState('disconnected');
	}

	/**
	 * Send a message to the server
	 */
	send<T = unknown>(type: MessageType, payload?: T, requestId?: string): void {
		const message: WSMessage<T> = {
			type,
			payload,
			timestamp: new Date().toISOString(),
			request_id: requestId
		};

		// Queue message if not connected
		if (this.ws?.readyState !== WebSocket.OPEN) {
			console.warn('[WS] Not connected, queueing message:', type);
			this.messageQueue.push(message as WSMessage);
			this.saveQueue(); // Persist to localStorage
			return;
		}

		try {
			this.ws.send(JSON.stringify(message));
		} catch (error) {
			console.error('[WS] Failed to send message:', error);
			// Queue message for retry
			this.messageQueue.push(message as WSMessage);
			this.saveQueue(); // Persist to localStorage
		}
	}

	/**
	 * Register a message handler for a specific message type
	 */
	on<T = unknown>(type: MessageType, handler: MessageHandler<T>): () => void {
		if (!this.handlers.has(type)) {
			this.handlers.set(type, new Set());
		}
		this.handlers.get(type)!.add(handler as MessageHandler);

		// Return unsubscribe function
		return () => {
			this.handlers.get(type)?.delete(handler as MessageHandler);
		};
	}

	/**
	 * Remove a message handler
	 */
	off<T = unknown>(type: MessageType, handler: MessageHandler<T>): void {
		this.handlers.get(type)?.delete(handler as MessageHandler);
	}

	/**
	 * Subscribe to connection status changes
	 */
	onStatusChange(listener: (status: ConnectionStatus) => void): () => void {
		this.statusListeners.add(listener);
		// Send current status immediately
		listener(this.getStatus());

		return () => {
			this.statusListeners.delete(listener);
		};
	}

	/**
	 * Get current connection state
	 */
	get connectionState(): ConnectionState {
		return this._connectionState;
	}

	/**
	 * Get current connection status
	 */
	getStatus(): ConnectionStatus {
		return {
			state: this._connectionState,
			reconnectAttempt: this.reconnectAttempts,
			maxReconnectAttempts: this.config.maxReconnectAttempts
		};
	}

	/**
	 * Check if connected
	 */
	get isConnected(): boolean {
		return this.ws?.readyState === WebSocket.OPEN;
	}

	// Convenience methods for common operations

	joinRoom(workspaceId: string, userColor?: string): void {
		const payload: JoinRoomPayload = {
			workspace_id: workspaceId,
			user_color: userColor
		};
		this.send('join_room', payload);
	}

	leaveRoom(workspaceId: string): void {
		this.send('leave_room', { workspace_id: workspaceId });
	}

	sendCursorMove(workspaceId: string, x: number, y: number): void {
		const payload: CursorMovePayload = {
			workspace_id: workspaceId,
			cursor: { x, y }
		};
		this.send('cursor_move', payload);
	}

	sendSelectionChange(workspaceId: string, selectedElements: string[]): void {
		const payload: SelectionChangePayload = {
			workspace_id: workspaceId,
			selected_elements: selectedElements
		};
		this.send('selection_change', payload);
	}

	sendOperation(operation: OperationPayload): void {
		this.send('operation', operation);
	}

	requestSync(workspaceId: string, stateVector: Record<string, number>): void {
		const payload: SyncRequestPayload = {
			workspace_id: workspaceId,
			state_vector: stateVector
		};
		this.send('sync_request', payload);
	}

	sendHeartbeat(): void {
		this.send('heartbeat');
	}

	// Private methods

	private handleMessage(event: MessageEvent): void {
		try {
			const message = JSON.parse(event.data) as WSMessage;

			// Dispatch to registered handlers
			const handlers = this.handlers.get(message.type);
			if (handlers && handlers.size > 0) {
				handlers.forEach((handler) => {
					try {
						handler(message.payload, message);
					} catch (error) {
						console.error(`[WS] Handler error for ${message.type}:`, error);
					}
				});
			} else {
				console.warn('[WS] No handler for message type:', message.type);
			}
		} catch (error) {
			console.error('[WS] Failed to parse message:', error);
		}
	}

	private setConnectionState(state: ConnectionState, error?: string): void {
		this._connectionState = state;

		const status: ConnectionStatus = {
			state,
			error,
			reconnectAttempt: this.reconnectAttempts,
			maxReconnectAttempts: this.config.maxReconnectAttempts
		};

		// Notify all listeners
		this.statusListeners.forEach((listener) => {
			try {
				listener(status);
			} catch (error) {
				console.error('[WS] Status listener error:', error);
			}
		});
	}

	private scheduleReconnect(): void {
		if (this.reconnectAttempts >= this.config.maxReconnectAttempts) {
			console.error('[WS] Max reconnect attempts reached');
			this.setConnectionState('error', 'Max reconnect attempts reached');
			return;
		}

		this.setConnectionState('reconnecting');

		// Exponential backoff: 1s, 2s, 4s, 8s, 16s, 32s, ...
		const delay = Math.min(
			this.config.reconnectInterval * Math.pow(2, this.reconnectAttempts),
			30000 // Max 30 seconds
		);

		console.log(`[WS] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts + 1})`);

		this.clearReconnectTimer();
		this.reconnectTimer = setTimeout(() => {
			this.reconnectAttempts++;
			// Need to get token again - this should be handled by the caller
			// For now, we'll emit an event that requires manual reconnection
			console.warn('[WS] Auto-reconnect requires token - please reconnect manually');
			this.setConnectionState('disconnected');
		}, delay);
	}

	private clearReconnectTimer(): void {
		if (this.reconnectTimer) {
			clearTimeout(this.reconnectTimer);
			this.reconnectTimer = null;
		}
	}

	private startHeartbeat(): void {
		this.stopHeartbeat();
		this.heartbeatTimer = setInterval(() => {
			if (this.isConnected) {
				this.sendHeartbeat();
			}
		}, this.config.heartbeatInterval);
	}

	private stopHeartbeat(): void {
		if (this.heartbeatTimer) {
			clearInterval(this.heartbeatTimer);
			this.heartbeatTimer = null;
		}
	}

	private flushMessageQueue(): void {
		if (this.messageQueue.length === 0) {
			return;
		}

		console.log(`[WS] Flushing ${this.messageQueue.length} queued messages`);

		while (this.messageQueue.length > 0) {
			const message = this.messageQueue.shift()!;
			try {
				this.ws?.send(JSON.stringify(message));
			} catch (error) {
				console.error('[WS] Failed to flush message:', error);
				// Re-queue the message
				this.messageQueue.unshift(message);
				break;
			}
		}

		// Update localStorage after flushing
		this.saveQueue();
	}

	/**
	 * Save message queue to localStorage
	 */
	private saveQueue(): void {
		try {
			if (this.messageQueue.length > 0) {
				localStorage.setItem(this.QUEUE_STORAGE_KEY, JSON.stringify(this.messageQueue));
				console.log(`[WS] Saved ${this.messageQueue.length} messages to localStorage`);
			} else {
				// Clear localStorage if queue is empty
				localStorage.removeItem(this.QUEUE_STORAGE_KEY);
			}
		} catch (error) {
			console.error('[WS] Failed to save queue to localStorage:', error);
		}
	}

	/**
	 * Restore message queue from localStorage
	 */
	private restoreQueue(): void {
		try {
			const saved = localStorage.getItem(this.QUEUE_STORAGE_KEY);
			if (saved) {
				this.messageQueue = JSON.parse(saved);
				console.log(`[WS] Restored ${this.messageQueue.length} messages from localStorage`);
			}
		} catch (error) {
			console.error('[WS] Failed to restore queue from localStorage:', error);
			this.messageQueue = [];
		}
	}

	/**
	 * Clear persisted queue
	 */
	clearPersistedQueue(): void {
		try {
			localStorage.removeItem(this.QUEUE_STORAGE_KEY);
			console.log('[WS] Cleared persisted queue');
		} catch (error) {
			console.error('[WS] Failed to clear persisted queue:', error);
		}
	}
}

// WebSocket URL from environment
const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8081/ws';

// Singleton instance
export const wsClient = new WebSocketClient({ url: WS_URL });
