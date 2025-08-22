package gateway

import (
	"log"

	"github.com/gin-gonic/gin"
	grpcInternal "github.com/travis-james/DBCache/internal/client"
)

type HttpGateway struct {
	router     *gin.Engine
	grpcClient grpcInternal.Client
}

func NewHTTPGateway() (*HttpGateway, error) {
	gClient, err := grpcInternal.Start()
	if err != nil {
		return nil, err
	}

	serv := &HttpGateway{
		router:     NewRouter(gClient),
		grpcClient: *gClient,
	}
	return serv, nil
}

func (gw *HttpGateway) StartHTTPServer() {
	addr := "localhost:8080"
	log.Printf("starting HTTP server at %v", addr)
	if err := gw.router.Run(addr); err != nil {
		panic(err)
	}
}
