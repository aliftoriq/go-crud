package initializer

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectToRedis() {

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	RedisClient = client
}
