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
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/travis-james/DBCache/internal/config"
	"github.com/travis-james/DBCache/internal/datastore"
	"github.com/travis-james/DBCache/internal/datastore/postgres"
	"github.com/travis-james/DBCache/internal/datastore/redis"
	mm "github.com/travis-james/DBCache/internal/metrics"
	pb "github.com/travis-james/DBCache/pkg/protobuf"
)

var (
	// Error messages.
	ERR_FAILED_TO_VALIDATE_DATASTORES = "failed to verify health of datastores"
	ERR_CONFIRM_FLUSH                 = `"confirm" needs to be set to true to flush cache`
)

// Server contains all the components to run the GRPC db cache
// server.
type Server struct {
	pb.UnimplementedDBCacheServiceServer // Returns unimplemented errors for any rpc call not implemented.
	GRPCServer                           *grpc.Server
	DB                                   datastore.DB
	Cache                                datastore.Cache
	Config                               *config.Config
	metricsManager                       *mm.MetricsManager
}

// Init a grpc db cache server. Currently hardwired for postgres
// and redis.
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
		DB:             &pa,
		Cache:          &ra,
		Config:         config,
		metricsManager: mm.MetricsManagerInit(&ra),
	}, nil
}

// StartGRPCServer will start a grpc server with the parameters
// provided in config.
func (ss *Server) StartGRPCServer() error {
	log.Printf("cachepw %s, grpc port %s", ss.Config.CachePw, ss.Config.GRPCPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", ss.Config.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	ss.GRPCServer = grpc.NewServer()
	pb.RegisterDBCacheServiceServer(ss.GRPCServer, ss)
	log.Printf("grpc server listening at %v", lis.Addr())
	if err := ss.GRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

// RunGRPCServer is a convenience function to wrap Init and
// StartGRPCServer.
func RunGRPCServer() {
	ss, err := Init()
	if err != nil {
		panic(err)
	}
	ss.StartGRPCServer()
}

// Close the cache, db, grpc server connections, and
// metrics polling.
func (ss *Server) Close() {
	ss.Cache.Close()
	ss.DB.Close()
	if ss.GRPCServer != nil {
		ss.GRPCServer.Stop()
	}
	if ss.metricsManager != nil {
		ss.metricsManager.Close()
	}
}

// CheckHealth pings the DB and Cache.
func (ss Server) CheckHealth(ctx context.Context, _ *emptypb.Empty) (*pb.HealthCheckResponse, error) {
	log.Print("Ping'ing cache...")
	cacheResult, cacheErr := ss.Cache.Ping(ctx)
	if cacheErr != nil {
		return &pb.HealthCheckResponse{
			Healthy:    false,
			CacheError: cacheErr.Error(),
		}, errors.New(ERR_FAILED_TO_VALIDATE_DATASTORES)
	}
	if cacheResult != "PONG" {
		return &pb.HealthCheckResponse{
			Healthy:    false,
			CacheError: fmt.Sprintf("expected PONG, received %s", cacheResult),
		}, errors.New(ERR_FAILED_TO_VALIDATE_DATASTORES)
	}

	log.Print("Ping'ing db...")
	dbErr := ss.DB.Ping()
	if dbErr != nil {
		return &pb.HealthCheckResponse{
			Healthy: false,
			DbError: dbErr.Error(),
		}, errors.New(ERR_FAILED_TO_VALIDATE_DATASTORES)
	}

	return &pb.HealthCheckResponse{
		Healthy: true,
	}, nil
}

// GetData retrieves data from cache if available.
// Else run fallback query and cache the result.
func (ss Server) GetData(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	data, ttl, err := ss.Cache.Get(ctx, req.GetQueryId())
	if err == nil {
		ss.metricsManager.CacheHits.Inc()
		return &pb.GetResponse{
			FromCache:  true,
			Data:       data,
			TtlSeconds: ttl,
		}, nil
	} else if !errors.Is(err, datastore.ErrCacheMiss) {
		return nil, status.Error(codes.Internal, fmt.Sprint("error in querying cache: ", err.Error()))
	}

	// Cache miss.
	ss.metricsManager.CacheMisses.Inc()
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

// convertArgs is a helper function. Protobuf Fallbackquery is a slice
// of strings, but the QueryRows function expects a slice of any.
// Maybe check if protobuf allows defining an array/slice of
// any.
func convertArgs(args []string) []any {
	result := make([]any, len(args))
	for i, arg := range args {
		result[i] = arg
	}
	return result
}

// InsertData will insert the request data directly into the cache.
// !!! This call does not write to the database, cache only !!!
func (ss Server) InsertData(ctx context.Context, req *pb.InsertRequest) (*pb.InsertResponse, error) {
	err := ss.Cache.Set(ctx, req.GetQueryId(), req.GetData(), 0)
	if err != nil {
		return &pb.InsertResponse{Success: false}, status.Error(codes.Internal, fmt.Sprint("failed to insert into cache: ", err.Error()))
	}
	return &pb.InsertResponse{Success: true}, nil
}

// InvalidateCache removes the given key:val pair in the request from
// cache.
func (ss Server) InvalidateCache(ctx context.Context, req *pb.InvalidateRequest) (*pb.InvalidateResponse, error) {
	numberOfKeys, err := ss.Cache.Delete(ctx, req.Keys)
	if err != nil {
		return nil, err
	}
	return &pb.InvalidateResponse{
		EntriesRemoved: numberOfKeys,
	}, nil
}

// FlushCache removes all entries from cache.
func (ss Server) FlushCache(ctx context.Context, req *pb.FlushRequest) (*pb.FlushResponse, error) {
	if !req.GetConfirm() {
		return nil, errors.New(ERR_CONFIRM_FLUSH)
	}
	var entriesRemoved int64 = ss.Cache.Flush(ctx)
	return &pb.FlushResponse{NumberOfEntriesRemoved: entriesRemoved}, nil
}
