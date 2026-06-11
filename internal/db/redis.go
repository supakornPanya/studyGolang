package db

import (
	"context"
	"os"
	"github.com/redis/go-redis/v9"
)

// Initial Redis Connection
func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	return rdb
}