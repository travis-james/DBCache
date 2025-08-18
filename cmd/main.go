package main

import (
	"flag"
	"log"
	"net/http"

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
		go runHTTPServer()
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

func runHTTPServer() {
	gw, err := gateway.NewHTTPGateway()
	if err != nil {
		panic(err)
	}
	httpServer := http.Server{
		Addr:    ":8080",
		Handler: gw,
	}
	log.Printf("HTTP server listening at %v", httpServer.Addr)
	if err = httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
