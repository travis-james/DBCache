package server

import (
	"context"
	"testing"

	"github.com/travis-james/DBCache/internal/genproto"
)

func TestCheckHealth(t *testing.T) {
	testServer := Server{}
	got, err := testServer.CheckHealth(context.Background(), &genproto.HealthCheckRequest{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(got)
}
