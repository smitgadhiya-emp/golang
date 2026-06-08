package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Redis *redis.Client
	Ctx   = context.Background()
)

func InitRedis() {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		log.Println("⚠️  REDIS_URL not set; Redis-backed features (e.g. password reset) are disabled")
		return
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Printf("⚠️  invalid REDIS_URL, Redis disabled: %v", err)
		return
	}

	client := redis.NewClient(opt)

	if err := client.Ping(Ctx).Err(); err != nil {
		// Don't take the whole API down just because Redis is unreachable.
		// Redis-backed features will return errors until it recovers.
		log.Printf("⚠️  Redis ping failed, continuing without Redis: %v", err)
		Redis = client
		return
	}

	Redis = client
	log.Println("Redis Server Connected")
}
