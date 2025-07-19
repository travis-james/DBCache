package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

// TODO: ping db and cache.
func (ss Server) CheckHealth(_ context.Context, _ *pb.Empty) (*pb.HealthCheckResponse, error) {
	log.Print("CheckHealth request received")
	return &pb.HealthCheckResponse{Healthy: true}, nil
}

// GetData: Retrieve data from cache if available. Else run fallback query and cache the result.
func (ss Server) GetData(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	data, ttl, err := ss.Cache.Get(ctx, req.GetQueryId())
	if err == nil {
		return &pb.GetResponse{
			FromCache:  true,
			Data:       data,
			TtlSeconds: ttl,
		}, nil
	} else if !errors.Is(err, datastore.ErrCacheMiss) {
		return nil, status.Error(codes.Internal, fmt.Sprint("error in querying cache: ", err.Error()))
	}
	// Cache miss.
	args := convertArgs(req.GetFallbackQuery().GetArgs())
	dataFromDB, err := ss.DB.QueryRows(req.GetFallbackQuery().GetQuery(), args...)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprint("failed to query rows: ", err.Error()))
	}

	// Now put that db result in cache. Hmm, bit of code duplication here and InsertData.
	err = ss.Cache.Set(ctx, req.GetQueryId(), dataFromDB, 0)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprint("failed to insert into cache: ", err.Error()))
	}

	return &pb.GetResponse{
		FromCache: false,
		Data:      dataFromDB,
	}, err
}

func (ss Server) InsertData(ctx context.Context, req *pb.InsertRequest) (*pb.InsertResponse, error) {
	err := ss.Cache.Set(ctx, req.GetQueryId(), req.GetData(), 0)
	if err != nil {
		return &pb.InsertResponse{Success: false}, status.Error(codes.Internal, fmt.Sprint("failed to insert into cache: ", err.Error()))
	}
	return &pb.InsertResponse{Success: true}, nil
}

func convertArgs(args []string) []any {
	result := make([]any, len(args))
	for i, arg := range args {
		result[i] = arg
	}
	return result
}
