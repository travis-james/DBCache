package redis

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/travis-james/DBCache/internal/config"
	"github.com/travis-james/DBCache/internal/datastore"
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

	return RedisAdapter{
		Client: client,
	}, nil
}

func (ra RedisAdapter) Close() error {
	return ra.Client.Close()
}

func (ra RedisAdapter) Ping(ctx context.Context) (string, error) {
	return ra.Client.Ping(ctx).Result()
}

func (ra *RedisAdapter) Get(ctx context.Context, key string) ([]byte, int64, error) {
	data, err := ra.Client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, 0, datastore.ErrCacheMiss
	} else if err != nil {
		return nil, 0, err
	}
	time, err := ra.Client.TTL(ctx, key).Result()
	if err != nil {
		return nil, 0, err
	}
	return data, int64(time.Seconds()), nil
}

func (ra *RedisAdapter) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return ra.Client.Set(ctx, key, value, ttl).Err()
}

func (ra *RedisAdapter) Delete(ctx context.Context, keys []string) (int64, error) {
	return ra.Client.Del(ctx, keys...).Result()
}

func (ra *RedisAdapter) Flush(ctx context.Context) int64 {
	numberOfKeys := ra.Client.DBSize(ctx).Val()
	ra.Client.FlushDB(ctx)
	after := ra.Client.DBSize(ctx).Val() // I assume this returns the size of the currently selected db, and not the all different tables.
	if after != 0 {
		// something went wrong...?
		return 0
	}
	return numberOfKeys
}

func (ra *RedisAdapter) NumberOfItems(ctx context.Context) (int64, error) {
	return ra.Client.DBSize(ctx).Result()
}
