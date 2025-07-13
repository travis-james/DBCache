package redis

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/travis-james/DBCache/internal/config"
)

type RedisAdapter struct {
	Client *redis.Client
}

func NewRedis(cc *config.Config) (RedisAdapter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cc.CacheAddr,
		Username: cc.CacheUser,
		Password: cc.CachePw,
		DB:       0, // Use default DB for now.
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return RedisAdapter{}, fmt.Errorf("redis ping failed: %v", err)
	}
	return RedisAdapter{Client: client}, nil
}

func (ra *RedisAdapter) Get(ctx context.Context, key string) (string, error) {
	return ra.Client.Get(ctx, key).Result()
}

func (ra *RedisAdapter) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return ra.Client.Set(ctx, key, value, ttl).Err()
}

func (ra RedisAdapter) Close() error {
	return ra.Client.Close()
}
