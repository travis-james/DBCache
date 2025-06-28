package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/travis-james/DBCache/internal/genproto"
)

type Server struct {
	pb.UnimplementedDBCacheServiceServer
}

func (ss Server) CheckHealth(_ context.Context, _ *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
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
