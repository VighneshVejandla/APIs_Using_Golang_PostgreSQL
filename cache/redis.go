package cache

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

// ConnectRedis establishes a connection to Redis
func ConnectRedis() (*redis.Client, error) {
	redisAddr := os.Getenv("REDIS_ADDR") // Example: localhost:6379
	redisPassword := ""                   // No password set
	redisDB := 0                          // Default DB

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis")
	return client, nil
}
