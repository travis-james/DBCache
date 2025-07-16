//go:build integration
// +build integration

package tests

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/travis-james/DBCache/internal/client"
	"github.com/travis-james/DBCache/internal/server"
	pb "github.com/travis-james/DBCache/pkg/protobuf"
)

func TestHealthCheck(t *testing.T) {
	// Scenario: Server returns true for healthy when server is running.
	// server setup
	ss, err := server.Init()
	require.Nil(t, err)
	go ss.StartGRPCServer()
	defer ss.Close()

	// client setup
	cc, err := client.Start()
	require.Nil(t, err)
	defer cc.Close()

	// What we're testing.
	r, err := cc.CheckHealth(context.Background(), &pb.Empty{})
	require.Nil(t, err)
	assert.True(t, r.GetHealthy())
}

func TestGetData(t *testing.T) {
	// server setup
	ss, err := server.Init()
	require.Nil(t, err)
	go ss.StartGRPCServer()
	defer ss.Close()

	// client setup
	cc, err := client.Start()
	require.Nil(t, err)
	defer cc.Close()

	// What we're testing.
	tests := []struct {
		testName  string
		req       *pb.GetRequest
		fromCache bool
		expected  map[string]any
	}{
		{
			testName:  "cache hit",
			req:       &pb.GetRequest{QueryId: "users:1"},
			fromCache: true,
			expected: map[string]any{
				"id":    float64(1),
				"name":  "Alice",
				"email": "alice@example.com",
				"age":   float64(30),
			},
		},
		{
			testName: "cache miss",
			req: &pb.GetRequest{
				QueryId: "users:2",
				FallbackQuery: &pb.FallbackQuery{
					Query: "SELECT name, email, age FROM users WHERE id = $1",
					Args:  []string{"3"},
				},
			},
			fromCache: false,
			expected: map[string]any{
				"name":  "Kashino",
				"email": "kashiyuka@europe.gov",
				"age":   float64(22),
			},
		},
		{
			testName: "cache hit on what was previously a miss",
			req: &pb.GetRequest{
				QueryId: "users:2",
			},
			fromCache: true,
			expected: map[string]any{
				"name":  "Kashino",
				"email": "kashiyuka@europe.gov",
				"age":   float64(22),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			r, err := cc.GetData(context.Background(), test.req)
			require.Nil(t, err)
			require.NotNil(t, r)
			require.Equal(t, test.fromCache, r.FromCache)
			if len(test.expected) > 0 {
				var actual map[string]any
				err := json.Unmarshal(r.Data, &actual)
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}
