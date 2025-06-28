package main

import (
	"flag"

	"github.com/travis-james/DBCache/internal/datastore/postgres"
	"github.com/travis-james/DBCache/internal/datastore/redis"
	"github.com/travis-james/DBCache/pkg/server"
)

func main() {
	mode := flag.String("mode", "server", "Mode to run: server or client")
	flag.Parse()
	switch *mode {
	case "client":
		// TEST
	default:
		redis.CheckRedis()
		postgres.CheckPost()

		server.Start()
	}
}
