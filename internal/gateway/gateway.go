package gateway

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/travis-james/DBCache/internal/config"
	pb "github.com/travis-james/DBCache/pkg/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RunHTTPGateway sets up an HTTP gateway over the GRPC server
// using gin and grpc-gateway.
func RunHTTPGateway() {
	config, err := config.Load()
	if err != nil {
		panic(err) // TODO: maybe handle this differently.
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	err = pb.RegisterDBCacheServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("localhost:%s", config.GRPCPort),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Any("/v1/*any", gin.WrapH(mux))
	r.Run(fmt.Sprintf(":%s", config.HTTPServerPort))
}
