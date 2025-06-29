package client

import (
	"context"
	"log"
	"time"

	pb "github.com/travis-james/DBCache/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Start() {
	addr := "localhost:50051"
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewDBCacheServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.CheckHealth(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not check health: %v", err)
	}
	log.Printf("Healthy: %t", r.GetHealthy())
}
