package server

import (
	"context"
	"testing"

	"google.golang.org/protobuf/types/known/emptypb"
)

func TestCheckHealth(t *testing.T) {
	testServer := Server{}
	got, err := testServer.CheckHealth(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(got)
}
