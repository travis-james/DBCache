package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/travis-james/DBCache/internal/config"
	"github.com/travis-james/DBCache/internal/datastore"
	"github.com/travis-james/DBCache/internal/datastore/postgres"
	"github.com/travis-james/DBCache/internal/datastore/redis"
	pb "github.com/travis-james/DBCache/pkg/protobuf"
)

type Server struct {
	pb.UnimplementedDBCacheServiceServer
	GRPCServer *grpc.Server
	DB         datastore.DB
	Cache      datastore.Cache
	Config     *config.Config
}

func Init() (*Server, error) {
	config, err := config.Load()
	if err != nil {
		return &Server{}, err
	}

	pa, err := postgres.NewPostgres(config)
	if err != nil {
		return &Server{}, err
	}

	ra, err := redis.NewRedis(config)
	if err != nil {
		return &Server{}, err
	}

	return &Server{
		DB:     &pa,
		Cache:  &ra,
		Config: config,
	}, nil
}

func (ss Server) CheckHealth(_ context.Context, _ *pb.Empty) (*pb.HealthCheckResponse, error) {
	log.Print("CheckHealth request received")
	return &pb.HealthCheckResponse{Healthy: true}, nil
}

func (ss *Server) StartGRPCServer() {
	log.Printf("cachepw %s, grpc port %s", ss.Config.CachePw, ss.Config.GRPCPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", ss.Config.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ss.GRPCServer = grpc.NewServer()
	pb.RegisterDBCacheServiceServer(ss.GRPCServer, ss)
	log.Printf("server listening at %v", lis.Addr())
	if err := ss.GRPCServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (ss *Server) Close() error {
	ss.Cache.Close()
	ss.DB.Close()
	if ss.GRPCServer != nil {
		ss.GRPCServer.Stop()
	}
	return nil
}
