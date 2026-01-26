package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bifshteksex/hertz-board/internal/models"
	"github.com/bifshteksex/hertz-board/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking in production
		return true
	},
}

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512 KB

	// clientSendBufferSize is the buffer size for client send channel
	clientSendBufferSize = 256
)

type WebSocketHandler struct {
	hub        *service.Hub
	jwtService *service.JWTService
}

func NewWebSocketHandler(hub *service.Hub, jwtService *service.JWTService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:        hub,
		jwtService: jwtService,
	}
}

// HandleWebSocket handles WebSocket connections using gorilla/websocket
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get token from query parameter
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing authentication token", http.StatusUnauthorized)
		return
	}

	// Validate JWT token
	claims, err := h.jwtService.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
		return
	}

	// Get user ID from claims
	userID := claims.UserID

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	// Create client
	client := &models.Client{
		ID:       uuid.New(),
		UserID:   userID,
		Send:     make(chan *models.WSMessage, clientSendBufferSize),
		LastPing: time.Now(),
	}

	// Handle the connection
	h.handleConnection(conn, client, claims.Username)
}

// handleConnection manages the WebSocket connection lifecycle
func (h *WebSocketHandler) handleConnection(conn *websocket.Conn, client *models.Client, username string) {
	defer func() {
		conn.Close()
	}()

	// Configure connection
	conn.SetReadLimit(maxMessageSize)
	if err := conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("Failed to set read deadline: %v", err)
		return
	}
	conn.SetPongHandler(func(string) error {
		client.LastPing = time.Now()
		if err := conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Printf("Failed to set read deadline in pong handler: %v", err)
		}
		return nil
	})

	// Start goroutines for read and write
	go h.writePump(conn, client)
	h.readPump(conn, client, username)
}

// readPump reads messages from the WebSocket connection
func (h *WebSocketHandler) readPump(conn *websocket.Conn, client *models.Client, username string) {
	defer func() {
		// Unregister client when connection closes
		if client.WorkspaceID != uuid.Nil {
			h.hub.Unregister(client)
		}
	}()

	for {
		var msg models.WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Set user ID from client
		msg.UserID = client.UserID
		msg.Timestamp = time.Now()

		// Handle message based on type
		h.handleMessage(client, username, &msg)
	}
}

// writePump writes messages to the WebSocket connection
func (h *WebSocketHandler) writePump(conn *websocket.Conn, client *models.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if err := conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("Failed to set write deadline: %v", err)
				return
			}
			if !ok {
				// Channel closed
				if err := conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("Failed to write close message: %v", err)
				}
				return
			}

			err := conn.WriteJSON(message)
			if err != nil {
				log.Printf("Write error: %v", err)
				return
			}

		case <-ticker.C:
			if err := conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("Failed to set write deadline: %v", err)
				return
			}
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (h *WebSocketHandler) handleMessage(client *models.Client, username string, msg *models.WSMessage) {
	switch msg.Type {
	case models.MessageTypeJoinRoom:
		h.handleJoinRoom(client, username, msg)

	case models.MessageTypeLeaveRoom:
		h.handleLeaveRoom(client)

	case models.MessageTypeCursorMove:
		h.handleCursorMove(client, msg)

	case models.MessageTypeSelectionChange:
		h.handleSelectionChange(client, msg)

	case models.MessageTypeOperation:
		h.handleOperation(client, msg)

	case models.MessageTypeBatch:
		h.handleBatch(client, msg)

	case models.MessageTypeSyncRequest:
		h.handleSyncRequest(client, msg)

	case models.MessageTypeHeartbeat:
		// Respond with pong
		client.Send <- &models.WSMessage{
			Type:      models.MessageTypePong,
			Timestamp: time.Now(),
		}

	case models.MessageTypeUserJoined, models.MessageTypeUserLeft, models.MessageTypePresenceUpdate,
		models.MessageTypeSyncResponse, models.MessageTypePong, models.MessageTypeError:
		// These message types are sent by the server, not received from clients
		// Just log and ignore
		log.Printf("Received server-only message type from client: %s", msg.Type)

	default:
		log.Printf("Unknown message type: %s", msg.Type)
		h.sendError(client, "unknown_message_type", fmt.Sprintf("Unknown message type: %s", msg.Type))
	}
}

// handleJoinRoom handles join_room messages
func (h *WebSocketHandler) handleJoinRoom(client *models.Client, username string, msg *models.WSMessage) {
	// Parse payload
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		h.sendError(client, "invalid_payload", "Invalid join_room payload")
		return
	}

	workspaceIDStr, ok := payload["workspace_id"].(string)
	if !ok {
		h.sendError(client, "invalid_workspace_id", "Invalid workspace_id")
		return
	}

	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		h.sendError(client, "invalid_workspace_id", "Invalid workspace_id format")
		return
	}

	// Get user color or generate one
	userColor, _ := payload["user_color"].(string)
	if userColor == "" {
		userColor = generateUserColor(client.UserID)
	}

	// Update client info
	client.WorkspaceID = workspaceID
	client.UserName = username
	client.UserColor = userColor
	client.Presence = &models.UserPresence{
		UserID:    client.UserID,
		UserName:  username,
		UserColor: userColor,
		LastSeen:  time.Now(),
	}

	// Register client to hub
	h.hub.Register(client)

	log.Printf("User %s joined workspace %s", client.UserID, workspaceID)
}

// handleLeaveRoom handles leave_room messages
func (h *WebSocketHandler) handleLeaveRoom(client *models.Client) {
	if client.WorkspaceID != uuid.Nil {
		h.hub.Unregister(client)
		client.WorkspaceID = uuid.Nil
	}
}

// handleCursorMove handles cursor movement
func (h *WebSocketHandler) handleCursorMove(client *models.Client, msg *models.WSMessage) {
	if client.WorkspaceID == uuid.Nil {
		return
	}

	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return
	}

	position, ok := payload["position"].(map[string]interface{})
	if !ok {
		return
	}

	x, _ := position["x"].(float64)
	y, _ := position["y"].(float64)

	// Update client presence
	if client.Presence != nil {
		client.Presence.Cursor = &models.CursorPosition{X: x, Y: y}
		client.Presence.LastSeen = time.Now()
	}

	// Broadcast to room
	h.hub.BroadcastToRoom(client.WorkspaceID, &models.WSMessage{
		Type:      models.MessageTypePresenceUpdate,
		UserID:    client.UserID,
		Timestamp: time.Now(),
		Payload: models.PresenceUpdatePayload{
			Presence: *client.Presence,
		},
	}, client.ID)
}

// handleSelectionChange handles selection changes
func (h *WebSocketHandler) handleSelectionChange(client *models.Client, msg *models.WSMessage) {
	if client.WorkspaceID == uuid.Nil {
		return
	}

	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return
	}

	elementIDsRaw, ok := payload["element_ids"].([]interface{})
	if !ok {
		return
	}

	elementIDs := make([]uuid.UUID, 0, len(elementIDsRaw))
	for _, idRaw := range elementIDsRaw {
		if idStr, ok := idRaw.(string); ok {
			if id, err := uuid.Parse(idStr); err == nil {
				elementIDs = append(elementIDs, id)
			}
		}
	}

	// Update client presence
	if client.Presence != nil {
		client.Presence.SelectedElements = elementIDs
		client.Presence.LastSeen = time.Now()
	}

	// Broadcast to room
	h.hub.BroadcastToRoom(client.WorkspaceID, &models.WSMessage{
		Type:      models.MessageTypePresenceUpdate,
		UserID:    client.UserID,
		Timestamp: time.Now(),
		Payload: models.PresenceUpdatePayload{
			Presence: *client.Presence,
		},
	}, client.ID)
}

// handleOperation handles CRDT operations
func (h *WebSocketHandler) handleOperation(client *models.Client, msg *models.WSMessage) {
	if client.WorkspaceID == uuid.Nil {
		return
	}

	// Broadcast operation to other clients
	h.hub.BroadcastToRoom(client.WorkspaceID, msg, client.ID)

	// TODO: Store operation in database for persistence
}

// handleBatch handles batch operations
func (h *WebSocketHandler) handleBatch(client *models.Client, msg *models.WSMessage) {
	if client.WorkspaceID == uuid.Nil {
		return
	}

	// Broadcast batch to other clients
	h.hub.BroadcastToRoom(client.WorkspaceID, msg, client.ID)

	// TODO: Store operations in database for persistence
}

// handleSyncRequest handles sync requests
func (h *WebSocketHandler) handleSyncRequest(client *models.Client, msg *models.WSMessage) {
	if client.WorkspaceID == uuid.Nil {
		return
	}

	// TODO: Implement sync logic
	// For now, send empty response
	client.Send <- &models.WSMessage{
		Type:      models.MessageTypeSyncResponse,
		Timestamp: time.Now(),
		Payload: models.SyncResponsePayload{
			Operations:  []models.OperationPayload{},
			StateVector: make(map[string]int64),
		},
		RequestID: msg.RequestID,
	}
}

// sendError sends an error message to the client
func (h *WebSocketHandler) sendError(client *models.Client, code, message string) {
	client.Send <- &models.WSMessage{
		Type:      models.MessageTypeError,
		Timestamp: time.Now(),
		Payload: models.ErrorPayload{
			Code:    code,
			Message: message,
		},
	}
}

// generateUserColor generates a consistent color for a user based on their ID
func generateUserColor(userID uuid.UUID) string {
	colors := []string{
		"#FF6B6B", "#4ECDC4", "#45B7D1", "#FFA07A",
		"#98D8C8", "#F7DC6F", "#BB8FCE", "#85C1E2",
		"#F8B739", "#52B788", "#E76F51", "#2A9D8F",
	}

	// Use user ID bytes to select color
	bytes := userID[:]
	index := int(bytes[0]) % len(colors)
	return colors[index]
}
