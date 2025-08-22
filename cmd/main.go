package main

import (
	"flag"

	"github.com/travis-james/DBCache/internal/client"
	"github.com/travis-james/DBCache/internal/gateway"
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
		go runHTTPGateway()
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

func runHTTPGateway() {
	gw, err := gateway.NewHTTPGateway()
	if err != nil {
		panic(err)
	}
	gw.StartHTTPServer()
}
