package datastore

import (
	"context"
	"database/sql"
	"time"
)

type DB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	// Exec(query string, args ...any) (sql.Result, error)
	Close() error
}

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	// Del(ctx context.Context, key string) error
	Close() error
}
