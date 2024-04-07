package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nikhilsiwach28/MyCode.git/config"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(redisConfig config.RedisConfig) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	// Ping the Redis server to check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &RedisService{
		client: client,
	}
}

func (r *RedisService) Set(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set value in Redis: %v", err)
	}

	if v, err := r.Get(key); err != nil {
		fmt.Println("ERROR = ", err)
	} else {
		fmt.Println("Value = ", v)
	}

	return nil
}
func (r *RedisService) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// Key does not exist
			return "", fmt.Errorf("key '%s' not found in Redis", key)
		}
		// Other error occurred
		return "", fmt.Errorf("failed to get value from Redis: %v", err)
	}
	return value, nil
}
