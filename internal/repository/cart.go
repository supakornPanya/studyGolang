package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"study-golang-backend/internal/domain/entity"
	"study-golang-backend/internal/domain/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type CartRepository struct {
	redisClient *redis.Client
	ctx         context.Context
}

// NewCachedRepository for warps and return value
func NewCartRepository(redisClient *redis.Client) repository.CartRepository {
	return &CartRepository{
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

// Helper for format query in redis
func (r *CartRepository) getCacheKey(id int) string {
	return fmt.Sprintf("cart: %d", id)
}

// Save
func (r *CartRepository) Save(id int, cart *entity.CartItems) (*entity.CartItems, error) {
	// Crate key
	key := r.getCacheKey(id)

	// Tracking Time
	if cart.CreatedAt.IsZero() {
		cart.CreatedAt = time.Now()
	}
	cart.UpdatedAt = time.Now()

	// Convert the CartItems struct into a JSON string
	data, err := json.Marshal(cart)
	if err != nil {
		return nil, err
	}

	// Set data to redis with expiry time
	err = r.redisClient.Set(r.ctx, key, data, 1*time.Minute).Err()

	if err != nil {
		return nil, err
	}

	return cart, nil
}

// GetByID
func (r *CartRepository) GetByID(id int) (*entity.CartItems, error) {
	key := r.getCacheKey(id)

	// Get value from cache
	val, err := r.redisClient.Get(r.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return &entity.CartItems{
			UserID: int(id),
			ListProductID: []uint64{},
		}, nil
	} else if err != nil {
		return nil, err	
	}

	// Get value from cache and convert to CartItems struct
	var item entity.CartItems
	if err := json.Unmarshal([]byte(val), &item); err == nil {
		return &item, nil
	}

	return nil, err
}

// Delete
func (r *CartRepository) Delete(id int) error {
	key := r.getCacheKey(id)
	return r.redisClient.Del(r.ctx, key).Err()
}
