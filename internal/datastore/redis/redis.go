package redis

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/travis-james/DBCache/internal/config"
	"github.com/travis-james/DBCache/internal/datastore"
)

// RedisAdapter contains a Redis Client and implements the Cache
// interface in datastore.
type RedisAdapter struct {
	Client *redis.Client
}

// NewRedis returns a Redic Client connection based on the input
// config.
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

// Close the Redis connection.
func (ra RedisAdapter) Close() error {
	return ra.Client.Close()
}

// Ping the Redis cache.
func (ra RedisAdapter) Ping(ctx context.Context) (string, error) {
	return ra.Client.Ping(ctx).Result()
}

// Get the associated data in cache based on the input key.
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

// Set the cache with the given key : value.
func (ra *RedisAdapter) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return ra.Client.Set(ctx, key, value, ttl).Err()
}

// Delete an entry/entries from cache. Returns the number of keys
// deleted.
func (ra *RedisAdapter) Delete(ctx context.Context, keys []string) (int64, error) {
	return ra.Client.Del(ctx, keys...).Result()
}

// Flush remove all entries from cache and returns the number
// of items removed.
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

// NumberOfItems returns the number of items in a given cache.
func (ra *RedisAdapter) NumberOfItems(ctx context.Context) (int64, error) {
	return ra.Client.DBSize(ctx).Result()
}
