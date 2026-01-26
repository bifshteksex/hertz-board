package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/bifshteksex/hertz-board/internal/models"
	"github.com/bifshteksex/hertz-board/internal/repository"

	"github.com/google/uuid"
)

const (
	// maxOperationsToFetch is the maximum number of operations to fetch from the database
	maxOperationsToFetch = 1000
)

// LamportClock implements a Lamport timestamp for ordering operations
type LamportClock struct {
	counter int64
	mu      sync.Mutex
}

// NewLamportClock creates a new Lamport clock
func NewLamportClock() *LamportClock {
	return &LamportClock{
		counter: 0,
	}
}

// Tick increments the clock and returns the new value
func (lc *LamportClock) Tick() int64 {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.counter++
	return lc.counter
}

// Update updates the clock based on a received timestamp
func (lc *LamportClock) Update(receivedTimestamp int64) int64 {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	if receivedTimestamp > lc.counter {
		lc.counter = receivedTimestamp
	}
	lc.counter++
	return lc.counter
}

// Get returns the current clock value
func (lc *LamportClock) Get() int64 {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	return lc.counter
}

// CRDTService handles CRDT-based synchronization
type CRDTService struct {
	elementRepo   *repository.ElementRepository
	operationRepo *repository.OperationRepository
	clock         *LamportClock
	ctx           context.Context
}

// NewCRDTService creates a new CRDT service
func NewCRDTService(
	elementRepo *repository.ElementRepository,
	operationRepo *repository.OperationRepository,
) *CRDTService {
	return &CRDTService{
		elementRepo:   elementRepo,
		operationRepo: operationRepo,
		clock:         NewLamportClock(),
		ctx:           context.Background(),
	}
}

// ApplyOperation applies a CRDT operation and returns the resulting element state
func (s *CRDTService) ApplyOperation(op *models.OperationPayload) error {
	// Update Lamport clock
	s.clock.Update(op.Timestamp)

	// Store operation in database
	err := s.operationRepo.Create(s.ctx, &models.Operation{
		ID:          uuid.New(),
		WorkspaceID: op.WorkspaceID,
		ElementID:   op.ElementID,
		UserID:      op.UserID,
		OpType:      string(op.OpType),
		Data:        op.Data,
		Timestamp:   op.Timestamp,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to store operation: %w", err)
	}

	// Apply operation to element
	switch op.OpType {
	case models.OperationTypeCreate:
		return s.applyCreate(op)
	case models.OperationTypeUpdate:
		return s.applyUpdate(op)
	case models.OperationTypeDelete:
		return s.applyDelete(op)
	case models.OperationTypeMove:
		return s.applyMove(op)
	default:
		return fmt.Errorf("unknown operation type: %s", op.OpType)
	}
}

// applyCreate creates a new element
func (s *CRDTService) applyCreate(op *models.OperationPayload) error {
	// Check if element already exists (idempotent operation)
	existing, err := s.elementRepo.GetByID(s.ctx, op.ElementID)
	if err == nil && existing != nil {
		// Element exists, check timestamp for LWW
		if op.Timestamp <= existing.Version {
			// Ignore older operation
			return nil
		}
	}

	// Parse element data from operation
	dataBytes, err := json.Marshal(op.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal element data: %w", err)
	}

	var elementData map[string]interface{}
	err = json.Unmarshal(dataBytes, &elementData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal element data: %w", err)
	}

	// Extract element fields
	elementType, _ := elementData["type"].(string)
	content, _ := elementData["content"].(string)
	posX, _ := elementData["pos_x"].(float64)
	posY, _ := elementData["pos_y"].(float64)
	width, _ := elementData["width"].(float64)
	height, _ := elementData["height"].(float64)
	zIndex, _ := elementData["z_index"].(float64)
	rotation, _ := elementData["rotation"].(float64)

	// Extract style as JSON
	var styleData map[string]interface{}
	if style, ok := elementData["style"].(map[string]interface{}); ok {
		styleData = style
	}

	// Create element
	element := &models.Element{
		ID:          op.ElementID,
		WorkspaceID: op.WorkspaceID,
		Type:        elementType,
		Content:     content,
		PosX:        posX,
		PosY:        posY,
		Width:       width,
		Height:      height,
		ZIndex:      int(zIndex),
		Rotation:    rotation,
		Style:       styleData,
		Version:     op.Timestamp,
		CreatedBy:   op.UserID,
		UpdatedBy:   op.UserID,
	}

	return s.elementRepo.Create(s.ctx, element)
}

// applyUpdate updates an existing element using LWW (Last-Write-Wins)
func (s *CRDTService) applyUpdate(op *models.OperationPayload) error {
	// Get existing element
	existing, err := s.elementRepo.GetByID(s.ctx, op.ElementID)
	if err != nil {
		return fmt.Errorf("element not found: %w", err)
	}

	// Check timestamp for LWW - only apply if newer
	if op.Timestamp <= existing.Version {
		// Ignore older update
		return nil
	}

	// Parse update data
	dataBytes, err := json.Marshal(op.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal update data: %w", err)
	}

	var updateData map[string]interface{}
	err = json.Unmarshal(dataBytes, &updateData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal update data: %w", err)
	}

	// Apply updates to element (partial updates)
	if content, ok := updateData["content"].(string); ok {
		existing.Content = content
	}
	if posX, ok := updateData["pos_x"].(float64); ok {
		existing.PosX = posX
	}
	if posY, ok := updateData["pos_y"].(float64); ok {
		existing.PosY = posY
	}
	if width, ok := updateData["width"].(float64); ok {
		existing.Width = width
	}
	if height, ok := updateData["height"].(float64); ok {
		existing.Height = height
	}
	if zIndex, ok := updateData["z_index"].(float64); ok {
		existing.ZIndex = int(zIndex)
	}
	if rotation, ok := updateData["rotation"].(float64); ok {
		existing.Rotation = rotation
	}
	if style, ok := updateData["style"].(map[string]interface{}); ok {
		existing.Style = style
	}

	// Update version and user
	existing.Version = op.Timestamp
	existing.UpdatedBy = op.UserID

	return s.elementRepo.Update(s.ctx, existing)
}

// applyDelete marks an element as deleted using tombstone
func (s *CRDTService) applyDelete(op *models.OperationPayload) error {
	// Get existing element
	existing, err := s.elementRepo.GetByID(s.ctx, op.ElementID)
	if err != nil {
		// Element doesn't exist, operation is already applied
		return nil
	}

	// Check timestamp for LWW
	if op.Timestamp <= existing.Version {
		// Ignore older delete
		return nil
	}

	// Soft delete the element
	return s.elementRepo.Delete(s.ctx, op.ElementID)
}

// applyMove updates element position
func (s *CRDTService) applyMove(op *models.OperationPayload) error {
	// Get existing element
	existing, err := s.elementRepo.GetByID(s.ctx, op.ElementID)
	if err != nil {
		return fmt.Errorf("element not found: %w", err)
	}

	// Check timestamp for LWW
	if op.Timestamp <= existing.Version {
		// Ignore older move
		return nil
	}

	// Parse move data
	dataBytes, err := json.Marshal(op.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal move data: %w", err)
	}

	var moveData map[string]interface{}
	err = json.Unmarshal(dataBytes, &moveData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal move data: %w", err)
	}

	// Update position
	if posX, ok := moveData["pos_x"].(float64); ok {
		existing.PosX = posX
	}
	if posY, ok := moveData["pos_y"].(float64); ok {
		existing.PosY = posY
	}

	// Update version and user
	existing.Version = op.Timestamp
	existing.UpdatedBy = op.UserID

	return s.elementRepo.Update(s.ctx, existing)
}

// ResolveConflict resolves conflicts between concurrent operations
func (s *CRDTService) ResolveConflict(op1, op2 *models.OperationPayload) *models.OperationPayload {
	// Use Lamport timestamp for ordering
	if op1.Timestamp != op2.Timestamp {
		if op1.Timestamp > op2.Timestamp {
			return op1
		}
		return op2
	}

	// If timestamps are equal, use UserID as tiebreaker (deterministic)
	if op1.UserID.String() > op2.UserID.String() {
		return op1
	}
	return op2
}

// GetOperationsSince returns operations since a given state vector
func (s *CRDTService) GetOperationsSince(
	workspaceID uuid.UUID,
	stateVector map[string]int64,
) ([]*models.Operation, error) {
	// Get all operations for workspace
	operations, err := s.operationRepo.GetByWorkspaceID(s.ctx, workspaceID, maxOperationsToFetch)
	if err != nil {
		return nil, err
	}

	// Filter operations based on state vector
	result := make([]*models.Operation, 0)
	for _, op := range operations {
		userIDStr := op.UserID.String()
		lastSeen, exists := stateVector[userIDStr]

		// Include operation if:
		// 1. We haven't seen any operations from this user
		// 2. This operation is newer than what we've seen
		if !exists || op.Timestamp > lastSeen {
			result = append(result, op)
		}
	}

	return result, nil
}

// BuildStateVector builds a state vector from operations
func (s *CRDTService) BuildStateVector(operations []*models.Operation) map[string]int64 {
	stateVector := make(map[string]int64)

	for _, op := range operations {
		userIDStr := op.UserID.String()
		if op.Timestamp > stateVector[userIDStr] {
			stateVector[userIDStr] = op.Timestamp
		}
	}

	return stateVector
}

// GenerateTimestamp generates a new Lamport timestamp
func (s *CRDTService) GenerateTimestamp() int64 {
	return s.clock.Tick()
}
