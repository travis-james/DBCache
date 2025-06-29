package server

import (
	"context"
	"testing"

	pb "github.com/travis-james/DBCache/pkg/proto"
)

func TestCheckHealth(t *testing.T) {
	testServer := Server{}
	got, err := testServer.CheckHealth(context.Background(), &pb.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(got)
}
