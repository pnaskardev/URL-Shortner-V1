package cache

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client = nil
var RedisSyncOnce sync.Once

func NewRedisClient() {
	RedisSyncOnce.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
			Protocol: 2,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			panic("REDIS CONNECTION FAILED: " + err.Error())
		}

		redisClient = client
	})
}

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		NewRedisClient()
	}

	return redisClient

}
