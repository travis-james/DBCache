//go:build integration
// +build integration

package tests

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/travis-james/DBCache/internal/client"
	"github.com/travis-james/DBCache/internal/server"
	pb "github.com/travis-james/DBCache/pkg/protobuf"
	"google.golang.org/protobuf/types/known/emptypb"
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
	r, err := cc.CheckHealth(context.Background(), &emptypb.Empty{})
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
			require.Equal(t, test.fromCache, r.GetFromCache())
			if len(test.expected) > 0 {
				var actual map[string]any
				err := json.Unmarshal(r.GetData(), &actual)
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}

func TestInsertData(t *testing.T) {
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
		testName     string
		insertReq    *pb.InsertRequest
		getReq       *pb.GetRequest
		expectedData []byte
	}{
		{
			testName: "a user is successfully inserted into cache",
			insertReq: &pb.InsertRequest{
				QueryId: "users:11",
				Data:    []byte(`{"id":11, "name":"Test Person", "email":"test@testing.gov", "age":67}`),
			},
			getReq:       &pb.GetRequest{QueryId: "users:11"},
			expectedData: []byte(`{"id":11, "name":"Test Person", "email":"test@testing.gov", "age":67}`),
		},
	}
	// Run tests.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			insertResp, err := cc.InsertData(context.Background(), test.insertReq)
			require.Nil(t, err)
			require.NotNil(t, insertResp)
			require.True(t, insertResp.GetSuccess())

			// Verify data is in cache.
			getDataResp, err := cc.GetData(context.Background(), test.getReq)
			require.Nil(t, err)
			require.NotNil(t, getDataResp)
			require.True(t, getDataResp.GetFromCache())
			require.Equal(t, test.expectedData, getDataResp.GetData())
		})
	}
}

func TestInvalidateCache(t *testing.T) {
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
		testName        string
		req             *pb.InvalidateRequest
		expectedRemoved int64
	}{
		{
			testName: "invalidate a users from the cache",
			req: &pb.InvalidateRequest{
				Keys: []string{
					"users:1",
					"users:45", // doesn't exist
				},
			},
			expectedRemoved: 1,
		},
		{
			testName: "verify previous user was invalidated from cache",
			req: &pb.InvalidateRequest{
				Keys: []string{
					"users:1",
				},
			},
			expectedRemoved: 0,
		},
	}
	// Run tests.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			resp, err := cc.InvalidateCache(context.Background(), test.req)
			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Equal(t, test.expectedRemoved, resp.GetEntriesRemoved())
		})
	}
}

func TestFlushCache(t *testing.T) {
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
		testName    string
		req         *pb.FlushRequest
		expectedErr string
	}{
		{
			testName:    "confirm is false",
			req:         &pb.FlushRequest{},
			expectedErr: server.ERR_CONFIRM_FLUSH,
		},
		{
			testName: "happy path",
			req:      &pb.FlushRequest{Confirm: true},
		},
	}
	// Run tests.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			got, err := cc.FlushCache(context.Background(), test.req)
			if test.expectedErr != "" {
				require.NotNil(t, err)
				require.True(t, strings.Contains(err.Error(), test.expectedErr))
			} else {
				require.Nil(t, err)
				require.NotNil(t, got)
				require.Greater(t, got.GetNumberOfEntriesRemoved(), int64(0))
			}
		})
	}
}
