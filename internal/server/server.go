package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/travis-james/DBCache/internal/datastore"
	pb "github.com/travis-james/DBCache/pkg/proto"
)

type Server struct {
	pb.UnimplementedDBCacheServiceServer
	DB    datastore.DB
	Cache datastore.Cache
}

func (ss Server) CheckHealth(_ context.Context, _ *pb.Empty) (*pb.HealthCheckResponse, error) {
	log.Print("CheckHealth request received")
	return &pb.HealthCheckResponse{Healthy: true}, nil
}

func Start() {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDBCacheServiceServer(grpcServer, &Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
