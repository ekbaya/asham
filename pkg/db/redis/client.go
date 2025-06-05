package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

// InitRedis initializes the Redis client. Call this once during app startup.
func InitRedis(addr, password string, db int) {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     password,
			DB:           db,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		})

		// Optionally ping to verify connection
		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			panic(fmt.Sprintf("failed to connect to Redis: %v", err))
		}
	})
}

// GetRedis returns the singleton Redis client instance.
func GetRedis() *redis.Client {
	if redisClient == nil {
		panic("Redis client is not initialized. Call InitRedis first.")
	}
	return redisClient
}
