package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"study-golang-backend/internal/domain/entity"
	"study-golang-backend/internal/domain/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type CartCachedRepository struct {
	postgresRepo *CartPostgresRepository
	redisClient  *redis.Client
	ctx          context.Context
}

// NewCachedRepository for warps and return value
func NewCartCachedRepository(postgresRepo *CartPostgresRepository, redisClient *redis.Client) repository.CartRepository {
	return &CartCachedRepository{
		postgresRepo: postgresRepo,
		redisClient:  redisClient,
		ctx:          context.Background(),
	}
}

func (r *CartCachedRepository) getCacheKey(id int) string {
	return fmt.Sprintf("item: %d", id)
}

// Create
func (r *CartCachedRepository) Create(sku string, name string, qty int, price float64) (*entity.Item, error) {
	item, err := r.postgresRepo.Create(sku, name, qty, price)
	if err != nil {
		return nil, err
	}

	// Crate cache
	key := r.getCacheKey(item.ID)

	// Convert the Item struct into a JSON string
	data, err := json.Marshal(item)
	err = r.redisClient.Set(r.ctx, key, data, 1*time.Minute).Err()

	if err != nil {
		return nil, err
	}

	return item, nil
}

// GetAll
func (r *CartCachedRepository) GetAll() ([]*entity.Item, error) {
	return r.postgresRepo.GetAll()
}

// GetByID
func (r *CartCachedRepository) GetByID(id int) (*entity.Item, error) {
	key := r.getCacheKey(id)

	val, err := r.redisClient.Get(r.ctx, key).Result()

	// Cache hit -> return data from redis
	if err == nil {
		var item entity.Item
		if err := json.Unmarshal([]byte(val), &item); err == nil {
			return &item, nil
		}
	}

	// Cache miss
	item, err := r.postgresRepo.GetByID(id)

	// Return error
	if err != nil {
		return nil, err
	}

	// Set data to redis
	data, err := json.Marshal(item)
	if err == nil {
		_ = r.redisClient.Set(r.ctx, key, data, 5*time.Minute).Err()
	}
	return item, nil
}

// Update
func (r *CartCachedRepository) Update(id int, sku string, name string, qty int, price float64) (*entity.Item, error) {
	// Update the SQL database
	item, err := r.postgresRepo.Update(id, sku, name, qty, price)
	if err != nil {
		return nil, err
	}

	// Delete old cahe
	key := r.getCacheKey(id)
	r.redisClient.Del(r.ctx, key)
	return item, nil
}

// Delete
func (r *CartCachedRepository) Delete(id int) error {
	err := r.postgresRepo.Delete(id)
	if err != nil {
		return err
	}

	// Delete old cache
	key := r.getCacheKey(id)
	r.redisClient.Del(r.ctx, key)
	return nil
}
