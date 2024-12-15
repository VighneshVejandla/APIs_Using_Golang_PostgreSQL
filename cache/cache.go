// cache/cache.go
package cache

import (
    "context"
    "github.com/go-redis/redis/v8"
)

// GetCache retrieves the value from Redis using the provided key
func GetCache(client *redis.Client, key string) (string, error) {
    ctx := context.Background()
    val, err := client.Get(ctx, key).Result()
    if err == redis.Nil {
        return "", nil // Key does not exist
    } else if err != nil {
        return "", err // Some other error
    }
    return val, nil
}
