//go:build integration
// +build integration

package tests

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"os"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/travis-james/DBCache/internal/datastore/postgres"
	"github.com/travis-james/DBCache/internal/datastore/redis"
	"github.com/travis-james/DBCache/internal/server"
)

func TestMain(m *testing.M) {
	fmt.Println("🔧 TestMain is running")
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
	got, _, err := ra.Get(ctx, "users:1")
	require.Nil(t, err)
	expected := `{"id":1, "name":"Alice", "email":"alice@example.com", "age":30}`
	assert.Equal(t, expected, string(got))

	// Verify can put/get data in redis instance.
	var (
		putKey = "foo"
		putVal = []byte("bar")
	)
	err = ra.Set(ctx, putKey, putVal, 0)
	require.Nil(t, err)

	gotVal, _, err := ra.Get(ctx, putKey)
	require.Nil(t, err)
	assert.Equal(t, putVal, gotVal)
}

func TestVerifyLocalPostGresWorks(t *testing.T) {
	// Verifies that get and set work for /internal/datastore/redis
	ss, err := server.Init()
	require.Nil(t, err)

	pg, ok := ss.DB.(*postgres.PostgresAdapter)
	require.True(t, ok)

	// Query data
	rows, err := pg.DB.Query("SELECT id, name, email, age FROM users")
	require.Nil(t, err)
	defer rows.Close()

	// Iterate over rows
	var (
		receivedValues = map[int]string{}
		expectedValues = map[int]string{
			1: "Name=Alice, Email=alice@example.com, Age=30",
			2: "Name=Oomoto, Email=nocchi@perfume.com, Age=36",
			3: "Name=Kashino, Email=kashiyuka@europe.gov, Age=22",
			4: "Name=Nishiwaki, Email=aoaoan@neo.net, Age=49",
		}
	)
	for rows.Next() {
		var id int
		var name, email string
		var age int

		err := rows.Scan(&id, &name, &email, &age)
		require.Nil(t, err)

		receivedValues[id] = fmt.Sprintf("Name=%s, Email=%s, Age=%d", name, email, age)

	}
	assert.Equal(t, len(expectedValues), len(receivedValues))
	// Check values....
	for index := range expectedValues {
		assert.Equal(t, expectedValues[index], string(receivedValues[index]))
	}
	// Check for errors after iteration
	err = rows.Err()
	assert.Nil(t, err)
}
