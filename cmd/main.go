package main

import (
	"flag"

	"github.com/travis-james/DBCache/internal/client"
	"github.com/travis-james/DBCache/internal/datastore/postgres"
	"github.com/travis-james/DBCache/internal/datastore/redis"
	"github.com/travis-james/DBCache/internal/server"
)

func main() {
	mode := flag.String("mode", "server", "Mode to run: server or client")
	flag.Parse()
	switch *mode {
	case "client":
		client.Start()
	default:
		redis.CheckRedis()
		postgres.CheckPost()

		server.Start()
	}
}
