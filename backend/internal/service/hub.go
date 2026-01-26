package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bifshteksex/hertz-board/internal/models"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	maxClientsPerRoom   = 100 // Maximum clients allowed in a room
	roomCleanupInterval = 5 * time.Minute
	// channelBufferSize is the buffer size for broadcast and other channels
	channelBufferSize = 256
)

// Hub maintains the set of active rooms and clients
type Hub struct {
	// Rooms indexed by workspace ID
	rooms map[uuid.UUID]*models.Room

	// Redis client for pub/sub
	redis *redis.Client

	// Context for Redis operations
	ctx context.Context

	// Mutex for rooms map
	mu sync.RWMutex
}

// NewHub creates a new Hub
func NewHub(redisClient *redis.Client) *Hub {
	hub := &Hub{
		rooms: make(map[uuid.UUID]*models.Room),
		redis: redisClient,
		ctx:   context.Background(),
	}

	// Start room cleanup goroutine
	go hub.cleanupEmptyRooms()

	// Start Redis subscription
	go hub.subscribeToRedis()

	return hub
}

// Register registers a client to a room
func (h *Hub) Register(client *models.Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	workspaceID := client.WorkspaceID
	room, exists := h.rooms[workspaceID]

	if !exists {
		// Create new room
		room = &models.Room{
			WorkspaceID: workspaceID,
			Clients:     make(map[uuid.UUID]*models.Client),
			Broadcast:   make(chan *models.WSMessage, channelBufferSize),
			Register:    make(chan *models.Client),
			Unregister:  make(chan *models.Client),
		}
		h.rooms[workspaceID] = room

		// Start room goroutine
		go h.runRoom(room)

		log.Printf("Created new room for workspace %s", workspaceID)
	}

	// Check room capacity
	if len(room.Clients) >= maxClientsPerRoom {
		h.sendErrorToClient(client, "room_full", "Room has reached maximum capacity")
		return
	}

	// Register client to room
	room.Register <- client
}

// Unregister unregisters a client from a room
func (h *Hub) Unregister(client *models.Client) {
	h.mu.RLock()
	room, exists := h.rooms[client.WorkspaceID]
	h.mu.RUnlock()

	if exists {
		room.Unregister <- client
	}
}

// BroadcastToRoom broadcasts a message to all clients in a room except the sender
func (h *Hub) BroadcastToRoom(workspaceID uuid.UUID, msg *models.WSMessage, excludeClientID uuid.UUID) {
	h.mu.RLock()
	room, exists := h.rooms[workspaceID]
	h.mu.RUnlock()

	if exists {
		// Add exclude client ID to message metadata
		msgCopy := *msg
		room.Broadcast <- &msgCopy
	}

	// Publish to Redis for other server instances
	h.publishToRedis(workspaceID, msg, excludeClientID)
}

// runRoom manages a single room
func (h *Hub) runRoom(room *models.Room) {
	for {
		select {
		case client := <-room.Register:
			// Add client to room
			room.Clients[client.ID] = client

			log.Printf("Client %s joined room %s (%d total clients)",
				client.UserID, room.WorkspaceID, len(room.Clients))

			// Send list of existing users to new client
			h.sendExistingPresences(client, room)

			// Broadcast user_joined to other clients
			joinMsg := &models.WSMessage{
				Type:      models.MessageTypeUserJoined,
				UserID:    client.UserID,
				Timestamp: time.Now(),
				Payload: models.UserJoinedPayload{
					UserID:    client.UserID,
					UserName:  client.UserName,
					UserColor: client.UserColor,
				},
			}
			h.broadcastToRoomClients(room, joinMsg, client.ID)

		case client := <-room.Unregister:
			if _, ok := room.Clients[client.ID]; ok {
				// Remove client from room
				delete(room.Clients, client.ID)
				close(client.Send)

				log.Printf("Client %s left room %s (%d remaining clients)",
					client.UserID, room.WorkspaceID, len(room.Clients))

				// Broadcast user_left to other clients
				leaveMsg := &models.WSMessage{
					Type:      models.MessageTypeUserLeft,
					UserID:    client.UserID,
					Timestamp: time.Now(),
					Payload: models.UserLeftPayload{
						UserID: client.UserID,
					},
				}
				h.broadcastToRoomClients(room, leaveMsg, uuid.Nil)

				// If room is empty, it will be cleaned up by cleanupEmptyRooms
			}

		case message := <-room.Broadcast:
			// Broadcast message to all clients in room
			h.broadcastToRoomClients(room, message, uuid.Nil)
		}
	}
}

// broadcastToRoomClients sends a message to all clients in a room except excluded one
func (h *Hub) broadcastToRoomClients(room *models.Room, msg *models.WSMessage, excludeClientID uuid.UUID) {
	for clientID, client := range room.Clients {
		if excludeClientID != uuid.Nil && clientID == excludeClientID {
			continue
		}

		select {
		case client.Send <- msg:
		default:
			// Client's send buffer is full, close the connection
			close(client.Send)
			delete(room.Clients, clientID)
			log.Printf("Client %s send buffer full, closing connection", client.UserID)
		}
	}
}

// sendExistingPresences sends the list of existing users to a newly joined client
func (h *Hub) sendExistingPresences(client *models.Client, room *models.Room) {
	for _, existingClient := range room.Clients {
		if existingClient.ID == client.ID {
			continue
		}

		// Send user_joined for each existing user
		msg := &models.WSMessage{
			Type:      models.MessageTypeUserJoined,
			UserID:    existingClient.UserID,
			Timestamp: time.Now(),
			Payload: models.UserJoinedPayload{
				UserID:    existingClient.UserID,
				UserName:  existingClient.UserName,
				UserColor: existingClient.UserColor,
			},
		}
		client.Send <- msg

		// Send presence update if available
		if existingClient.Presence != nil {
			presenceMsg := &models.WSMessage{
				Type:      models.MessageTypePresenceUpdate,
				UserID:    existingClient.UserID,
				Timestamp: time.Now(),
				Payload: models.PresenceUpdatePayload{
					Presence: *existingClient.Presence,
				},
			}
			client.Send <- presenceMsg
		}
	}
}

// sendErrorToClient sends an error message to a client
func (h *Hub) sendErrorToClient(client *models.Client, code, message string) {
	client.Send <- &models.WSMessage{
		Type:      models.MessageTypeError,
		Timestamp: time.Now(),
		Payload: models.ErrorPayload{
			Code:    code,
			Message: message,
		},
	}
}

// cleanupEmptyRooms periodically removes empty rooms
func (h *Hub) cleanupEmptyRooms() {
	ticker := time.NewTicker(roomCleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		h.mu.Lock()
		for workspaceID, room := range h.rooms {
			if len(room.Clients) == 0 {
				delete(h.rooms, workspaceID)
				log.Printf("Cleaned up empty room %s", workspaceID)
			}
		}
		h.mu.Unlock()
	}
}

// GetRoomStats returns statistics about a room
func (h *Hub) GetRoomStats(workspaceID uuid.UUID) (int, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, exists := h.rooms[workspaceID]
	if !exists {
		return 0, false
	}

	return len(room.Clients), true
}

// GetAllRoomStats returns statistics for all rooms
func (h *Hub) GetAllRoomStats() map[uuid.UUID]int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	stats := make(map[uuid.UUID]int)
	for workspaceID, room := range h.rooms {
		stats[workspaceID] = len(room.Clients)
	}

	return stats
}

// Redis Pub/Sub methods for scaling across multiple instances

type RedisMessage struct {
	WorkspaceID     uuid.UUID         `json:"workspace_id"`
	ExcludeClientID uuid.UUID         `json:"exclude_client_id"`
	Message         *models.WSMessage `json:"message"`
}

// publishToRedis publishes a message to Redis for other server instances
func (h *Hub) publishToRedis(workspaceID uuid.UUID, msg *models.WSMessage, excludeClientID uuid.UUID) {
	redisMsg := RedisMessage{
		WorkspaceID:     workspaceID,
		Message:         msg,
		ExcludeClientID: excludeClientID,
	}

	data, err := json.Marshal(redisMsg)
	if err != nil {
		log.Printf("Failed to marshal Redis message: %v", err)
		return
	}

	channel := fmt.Sprintf("workspace:%s", workspaceID)
	err = h.redis.Publish(h.ctx, channel, data).Err()
	if err != nil {
		log.Printf("Failed to publish to Redis: %v", err)
	}
}

// subscribeToRedis subscribes to Redis channels for workspace updates
func (h *Hub) subscribeToRedis() {
	pubsub := h.redis.PSubscribe(h.ctx, "workspace:*")
	defer pubsub.Close()

	ch := pubsub.Channel()

	log.Println("Started Redis subscription for workspace channels")

	for msg := range ch {
		var redisMsg RedisMessage
		err := json.Unmarshal([]byte(msg.Payload), &redisMsg)
		if err != nil {
			log.Printf("Failed to unmarshal Redis message: %v", err)
			continue
		}

		// Forward message to local room clients
		h.mu.RLock()
		room, exists := h.rooms[redisMsg.WorkspaceID]
		h.mu.RUnlock()

		if exists {
			h.broadcastToRoomClients(room, redisMsg.Message, redisMsg.ExcludeClientID)
		}
	}
}
