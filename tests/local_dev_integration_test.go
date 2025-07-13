//go:build integration
// +build integration

package tests

import (
	"context"
	"path/filepath"
	"testing"

	"os"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/travis-james/DBCache/internal/datastore/redis"
	"github.com/travis-james/DBCache/internal/server"
)

func TestMain(m *testing.M) {
	// Load env vars from file
	envPath := filepath.Join("..", "dbcache.env")
	err := godotenv.Load(envPath)
	if err != nil {
		panic("Failed to load env file: " + err.Error())
	}
	os.Exit(m.Run())
}

func TestVerifyLocalRedisWorks(t *testing.T) {
	// Verifies that get and set work for /internal/datastore/redis
	ss, err := server.Init()
	require.Nil(t, err)

	ra, ok := ss.Cache.(*redis.RedisAdapter)
	require.True(t, ok)

	ctx := context.Background()
	// Verify redis instance has seeded data.
	got, err := ra.Get(ctx, "users:1")
	require.Nil(t, err)
	expected := `{"id":1, "name":"Alice", "email":"alice@example.com", "age":30}`
	assert.Equal(t, expected, got)

	// Verify can put/get data in redis instance.
	var (
		putKey = "foo"
		putVal = "bar"
	)
	err = ra.Set(ctx, putKey, putVal, 0)
	require.Nil(t, err)

	gotVal, err := ra.Get(ctx, putKey)
	require.Nil(t, err)
	assert.Equal(t, putVal, gotVal)
}
