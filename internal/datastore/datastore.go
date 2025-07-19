package datastore

import (
	"context"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

type DB interface {
	Ping() error
	QueryRows(query string, args ...any) ([]byte, error)
	// Exec(query string, args ...any) (sql.Result, error)
	Close() error
}

type Cache interface {
	Ping(ctx context.Context) (string, error)
	Get(ctx context.Context, key string) ([]byte, int64, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	// Del(ctx context.Context, key string) error
	Close() error
}
