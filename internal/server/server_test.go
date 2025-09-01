package server

// Most tests are handled via integration tests.
// This is more of just a placeholder, not sure if there's
// reason to add unit tests or not.
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
