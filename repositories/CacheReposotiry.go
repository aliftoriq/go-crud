package repositories

import (
	"time"

	"github.com/aliftoriq/go-crud/initializer"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

//go:generate mockery --outpkg mocks --name CacheRepository
type CacheRepository interface {
	SetKey(ctx *gin.Context, key string, value interface{}, duration time.Duration) (err error)
	GetValueByKey(ctx *gin.Context, key string) (value string, exists bool, err error)
}

type cacheRepository struct {
	redis *redis.Client
}

func NewCacheRepository() CacheRepository {
	return &cacheRepository{redis: initializer.RedisClient}
}

func (rc *cacheRepository) SetKey(ctx *gin.Context, key string, value interface{}, duration time.Duration) (err error) {
	return initializer.RedisClient.SetEx(ctx, key, value, duration).Err()
}

func (rc *cacheRepository) GetValueByKey(ctx *gin.Context, key string) (value string, exists bool, err error) {
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
