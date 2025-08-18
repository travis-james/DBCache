package gateway

import (
	grpcInternal "github.com/travis-james/DBCache/internal/client"
)

type HttpGateway struct {
	*router
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
