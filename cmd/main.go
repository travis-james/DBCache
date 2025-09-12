package main

import (
	"flag"

	"github.com/travis-james/DBCache/internal/client"
	gw "github.com/travis-james/DBCache/internal/gateway"
	"github.com/travis-james/DBCache/internal/metrics"
	grpcServer "github.com/travis-james/DBCache/internal/server"
)

func main() {
	mode := flag.String("mode", "server", "Mode to run: server, dev, or client")
	flag.Parse()
	switch *mode {
	case "client":
		client.Start()
	case "dev":
		go grpcServer.RunGRPCServer()
		go gw.RunHTTPGateway()
		go metrics.RunPrometheusServer()
		select {}
	default:
		grpcServer.RunGRPCServer()
	}
}
