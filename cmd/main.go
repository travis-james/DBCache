package main

import (
	"flag"

	"github.com/travis-james/DBCache/internal/client"
	gw "github.com/travis-james/DBCache/internal/gateway"
	grpcServer "github.com/travis-james/DBCache/internal/server"
)

func main() {
	mode := flag.String("mode", "server", "Mode to run: server, dev, or client")
	flag.Parse()
	switch *mode {
	case "client":
		client.Start()
	case "dev":
		go runGRPCServer()
		go gw.RunHTTPGateway()
		select {}
	default:
		runGRPCServer()
	}
}

func runGRPCServer() {
	ss, err := grpcServer.Init()
	if err != nil {
		panic(err)
	}
	ss.StartGRPCServer()
}
