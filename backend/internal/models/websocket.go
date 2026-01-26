package models

import (
	"time"

	"github.com/google/uuid"
)

// MessageType defines the type of WebSocket message
type MessageType string

const (
	// Connection messages
	MessageTypeJoinRoom   MessageType = "join_room"
	MessageTypeLeaveRoom  MessageType = "leave_room"
	MessageTypeUserJoined MessageType = "user_joined"
	MessageTypeUserLeft   MessageType = "user_left"

	// Presence messages
	MessageTypeCursorMove      MessageType = "cursor_move"
	MessageTypeSelectionChange MessageType = "selection_change"
	MessageTypePresenceUpdate  MessageType = "presence_update"

	// Operation messages
	MessageTypeOperation MessageType = "operation"
	MessageTypeBatch     MessageType = "batch"

	// Sync messages
	MessageTypeSyncRequest  MessageType = "sync_request"
	MessageTypeSyncResponse MessageType = "sync_response"

	// Control messages
	MessageTypeHeartbeat MessageType = "heartbeat"
	MessageTypePong      MessageType = "pong"
	MessageTypeError     MessageType = "error"
)

// WSMessage represents a WebSocket message
type WSMessage struct {
	Payload   interface{} `json:"payload,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	UserID    uuid.UUID   `json:"user_id,omitempty"`
	Type      MessageType `json:"type"`
	RequestID string      `json:"request_id,omitempty"` // For request/response matching
}

// JoinRoomPayload is the payload for join_room message
type JoinRoomPayload struct {
	WorkspaceID uuid.UUID `json:"workspace_id"`
	UserColor   string    `json:"user_color,omitempty"` // Hex color for user cursor
}

// UserJoinedPayload is broadcast when a user joins
type UserJoinedPayload struct {
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	UserColor string    `json:"user_color"`
}

// UserLeftPayload is broadcast when a user leaves
type UserLeftPayload struct {
	UserID uuid.UUID `json:"user_id"`
}

// CursorPosition represents a user's cursor position
type CursorPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// CursorMovePayload is sent when user moves cursor
type CursorMovePayload struct {
	Position CursorPosition `json:"position"`
}

// SelectionChangePayload is sent when user changes selection
type SelectionChangePayload struct {
	ElementIDs []uuid.UUID `json:"element_ids"`
}

// UserPresence represents a user's presence in the workspace
type UserPresence struct {
	UserID           uuid.UUID       `json:"user_id"`
	Cursor           *CursorPosition `json:"cursor,omitempty"`
	SelectedElements []uuid.UUID     `json:"selected_elements,omitempty"`
	LastSeen         time.Time       `json:"last_seen"`
	UserName         string          `json:"user_name"`
	UserColor        string          `json:"user_color"`
}

// PresenceUpdatePayload is broadcast to other users
type PresenceUpdatePayload struct {
	Presence UserPresence `json:"presence"`
}

// OperationType defines the type of CRDT operation
type OperationType string

const (
	OperationTypeCreate OperationType = "create"
	OperationTypeUpdate OperationType = "update"
	OperationTypeDelete OperationType = "delete"
	OperationTypeMove   OperationType = "move"
)

// OperationPayload represents a CRDT operation
type OperationPayload struct {
	ElementID   uuid.UUID     `json:"element_id"`
	WorkspaceID uuid.UUID     `json:"workspace_id"`
	UserID      uuid.UUID     `json:"user_id"`
	Data        interface{}   `json:"data,omitempty"` // Element data for create/update
	Timestamp   int64         `json:"timestamp"`      // Lamport timestamp
	OpType      OperationType `json:"op_type"`
}

// BatchPayload contains multiple operations
type BatchPayload struct {
	Operations []OperationPayload `json:"operations"`
}

// SyncRequestPayload requests synchronization
type SyncRequestPayload struct {
	WorkspaceID uuid.UUID        `json:"workspace_id"`
	StateVector map[string]int64 `json:"state_vector"` // user_id -> last_seen_timestamp
}

// SyncResponsePayload contains operations to sync
type SyncResponsePayload struct {
	StateVector map[string]int64   `json:"state_vector"` // Current state vector
	Operations  []OperationPayload `json:"operations"`
}

// ErrorPayload represents an error message
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Client represents a connected WebSocket client
type Client struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	WorkspaceID uuid.UUID
	Presence    *UserPresence
	Send        chan *WSMessage // Channel for outbound messages
	LastPing    time.Time
	UserName    string
	UserColor   string
}

// Room represents a workspace collaboration room
type Room struct {
	WorkspaceID uuid.UUID
	Clients     map[uuid.UUID]*Client // client_id -> client
	Broadcast   chan *WSMessage       // Broadcast channel
	Register    chan *Client          // Register channel
	Unregister  chan *Client          // Unregister channel
}
