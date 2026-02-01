// WebSocket Message Types for Real-time Collaboration

export type MessageType =
	// Connection
	| 'join_room'
	| 'leave_room'
	| 'user_joined'
	| 'user_left'
	// Presence
	| 'cursor_move'
	| 'selection_change'
	| 'presence_update'
	// Operations (CRDT)
	| 'operation'
	| 'batch'
	// Sync
	| 'sync_request'
	| 'sync_response'
	// Control
	| 'heartbeat'
	| 'pong'
	| 'error';

export type OperationType = 'create' | 'update' | 'delete' | 'move';

export interface WSMessage<T = unknown> {
	type: MessageType;
	payload?: T;
	timestamp: string; // ISO string
	user_id?: string;
	request_id?: string; // For request/response pattern
}

// Connection Payloads

export interface JoinRoomPayload {
	workspace_id: string;
	user_color?: string; // Hex color for cursor
}

export interface LeaveRoomPayload {
	workspace_id: string;
}

export interface UserJoinedPayload {
	user_id: string;
	user_name: string;
	user_color: string;
	cursor?: CursorPosition;
}

export interface UserLeftPayload {
	user_id: string;
}

// Presence Payloads

export interface CursorPosition {
	x: number;
	y: number;
}

export interface CursorMovePayload {
	workspace_id: string;
	cursor: CursorPosition;
}

export interface SelectionChangePayload {
	workspace_id: string;
	selected_elements: string[]; // Element IDs
}

export interface UserPresence {
	user_id: string;
	user_name: string;
	user_color: string;
	cursor?: CursorPosition;
	selected_elements: string[];
	last_seen: string; // ISO timestamp
}

export interface PresenceUpdatePayload {
	presence: UserPresence;
}

// Operation Payloads (CRDT)

export interface OperationPayload {
	element_id: string;
	workspace_id: string;
	user_id: string;
	op_type: OperationType;
	data?: Record<string, unknown>; // Element data or partial update
	timestamp: number; // Lamport clock
}

export interface BatchPayload {
	operations: OperationPayload[];
}

// Sync Payloads

export interface SyncRequestPayload {
	workspace_id: string;
	state_vector: Record<string, number>; // user_id -> last_timestamp
}

export interface SyncResponsePayload {
	workspace_id: string;
	state_vector: Record<string, number>;
	operations: OperationPayload[];
}

// Control Payloads

export interface ErrorPayload {
	code: string;
	message: string;
	details?: string;
}

// WebSocket Connection State

export type ConnectionState =
	| 'connecting'
	| 'connected'
	| 'disconnected'
	| 'reconnecting'
	| 'error';

export interface ConnectionStatus {
	state: ConnectionState;
	error?: string;
	reconnectAttempt?: number;
	maxReconnectAttempts?: number;
}
