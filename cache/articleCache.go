// cache_operations.go
package cache

import (
	"time"

	"github.com/aliftoriq/go-crud/initializer"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetKey(ctx *gin.Context, key string, value interface{}, duration time.Duration) (err error) {
	return initializer.RedisClient.SetEx(ctx, key, value, duration).Err()
}

func GetValueByKey(ctx *gin.Context, key string) (value string, exists bool, err error) {
	cmd := initializer.RedisClient.Get(ctx, key)
	value, err = cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}

		return "", false, err
	}

	exists = true
	return
}
