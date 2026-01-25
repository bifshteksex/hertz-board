package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/bifshteksex/hertz-board/internal/models"
)

const (
	// Cache key patterns
	workspaceElementsKey = "workspace:%s:elements"
	elementKey           = "element:%s"

	// Cache TTLs
	workspaceElementsTTL = 5 * time.Minute
	elementTTL           = 10 * time.Minute
)

type CanvasCacheService struct {
	redis *redis.Client
}

func NewCanvasCacheService(redisClient *redis.Client) *CanvasCacheService {
	return &CanvasCacheService{
		redis: redisClient,
	}
}

// GetWorkspaceElements retrieves workspace elements from cache
func (s *CanvasCacheService) GetWorkspaceElements(ctx context.Context, workspaceID uuid.UUID) ([]models.CanvasElement, bool) {
	key := fmt.Sprintf(workspaceElementsKey, workspaceID)

	data, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}

	var elements []models.CanvasElement
	if err := json.Unmarshal(data, &elements); err != nil {
		return nil, false
	}

	return elements, true
}

// SetWorkspaceElements stores workspace elements in cache
func (s *CanvasCacheService) SetWorkspaceElements(ctx context.Context, workspaceID uuid.UUID, elements []models.CanvasElement) error {
	key := fmt.Sprintf(workspaceElementsKey, workspaceID)

	data, err := json.Marshal(elements)
	if err != nil {
		return fmt.Errorf("failed to marshal elements: %w", err)
	}

	if err := s.redis.Set(ctx, key, data, workspaceElementsTTL).Err(); err != nil {
		return fmt.Errorf("failed to cache elements: %w", err)
	}

	return nil
}

// InvalidateWorkspaceElements removes workspace elements from cache
func (s *CanvasCacheService) InvalidateWorkspaceElements(ctx context.Context, workspaceID uuid.UUID) error {
	key := fmt.Sprintf(workspaceElementsKey, workspaceID)

	if err := s.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to invalidate cache: %w", err)
	}

	return nil
}

// GetElement retrieves a single element from cache
func (s *CanvasCacheService) GetElement(ctx context.Context, elementID uuid.UUID) (*models.CanvasElement, bool) {
	key := fmt.Sprintf(elementKey, elementID)

	data, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}

	var element models.CanvasElement
	if err := json.Unmarshal(data, &element); err != nil {
		return nil, false
	}

	return &element, true
}

// SetElement stores a single element in cache
func (s *CanvasCacheService) SetElement(ctx context.Context, element *models.CanvasElement) error {
	key := fmt.Sprintf(elementKey, element.ID)

	data, err := json.Marshal(element)
	if err != nil {
		return fmt.Errorf("failed to marshal element: %w", err)
	}

	if err := s.redis.Set(ctx, key, data, elementTTL).Err(); err != nil {
		return fmt.Errorf("failed to cache element: %w", err)
	}

	return nil
}

// InvalidateElement removes a single element from cache
func (s *CanvasCacheService) InvalidateElement(ctx context.Context, elementID uuid.UUID) error {
	key := fmt.Sprintf(elementKey, elementID)

	if err := s.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to invalidate element cache: %w", err)
	}

	return nil
}

// InvalidateMultipleElements removes multiple elements from cache
func (s *CanvasCacheService) InvalidateMultipleElements(ctx context.Context, elementIDs []uuid.UUID) error {
	if len(elementIDs) == 0 {
		return nil
	}

	keys := make([]string, len(elementIDs))
	for i, id := range elementIDs {
		keys[i] = fmt.Sprintf(elementKey, id)
	}

	if err := s.redis.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("failed to invalidate element caches: %w", err)
	}

	return nil
}

// InvalidateWorkspaceCache invalidates all cache for a workspace
func (s *CanvasCacheService) InvalidateWorkspaceCache(ctx context.Context, workspaceID uuid.UUID) error {
	// Invalidate workspace elements list
	if err := s.InvalidateWorkspaceElements(ctx, workspaceID); err != nil {
		return err
	}

	// Note: Individual elements will expire naturally or be invalidated on update
	return nil
}

// WarmupCache pre-loads workspace elements into cache
func (s *CanvasCacheService) WarmupCache(ctx context.Context, workspaceID uuid.UUID, elements []models.CanvasElement) error {
	// Cache the full list
	if err := s.SetWorkspaceElements(ctx, workspaceID, elements); err != nil {
		return err
	}

	// Optionally cache individual elements (for frequently accessed ones)
	// This could be done selectively based on access patterns

	return nil
}
