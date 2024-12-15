package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitializeRedis(addr, password string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
}

func SetProductCache(ctx context.Context, key, value string) error {
	return redisClient.Set(ctx, key, value, 0).Err()
}

func GetProductCache(ctx context.Context, key string) (string, error) {
	return redisClient.Get(ctx, key).Result()
}
