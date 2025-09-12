package datastore

import (
	"context"
	"errors"
	"time"
)

// ErrCacheMiss returned when requested key from client isn't found
// in cache.
var ErrCacheMiss = errors.New("cache miss")

// DB is a database interface so that SQL, PostgreSQL, etc. can
// be implemented for the DBCache app.
type DB interface {
	// Ping checks database connectivity.
	Ping() error

	// QueryRows executes a query and returns the row data.
	QueryRows(query string, args ...any) ([]byte, error)

	// Close the DB connection.
	Close() error
}

// Cache is an interface so (hopefully) any cache can be implemented
// for the app.
type Cache interface {
	// Ping checks cache connectivity.
	Ping(ctx context.Context) (string, error)

	// Get the data corresponding to the given key.
	Get(ctx context.Context, key string) ([]byte, int64, error)

	// Set the cache with the given key:val pair.
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete the entry/entries at keys.
	Delete(ctx context.Context, keys []string) (int64, error)

	// Close the cache connection.
	Close() error

	// Flush the cache, removing all entries.
	Flush(ctx context.Context) int64

	// NumberOfItems returns the total items in cache.
	NumberOfItems(ctx context.Context) (int64, error)
}
