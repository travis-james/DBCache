//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/travis-james/DBCache/internal/client"
	"github.com/travis-james/DBCache/internal/server"
	pb "github.com/travis-james/DBCache/pkg/proto"
)

func TestHealthCheck(t *testing.T) {
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
