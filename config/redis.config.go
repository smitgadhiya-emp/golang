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
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}

	Redis = redis.NewClient(opt)

	if err := Redis.Ping(Ctx).Err(); err != nil {
		panic(err)
	}

	log.Println("Redis Server Connected")
}
