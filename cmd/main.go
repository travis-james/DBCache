package main

import (
	"flag"

	"github.com/travis-james/DBCache/internal/client"
	"github.com/travis-james/DBCache/internal/server"
)

func main() {
	mode := flag.String("mode", "server", "Mode to run: server or client")
	flag.Parse()
	switch *mode {
	case "client":
		client.Start()
	default:
		ss, err := server.Init()
		if err != nil {
			panic(err)
		}
		ss.StartGRPCServer()
	}
}
