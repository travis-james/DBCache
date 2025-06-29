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
		ss, err := server.Init()
		if err != nil {
			panic(err)
		}
		defer ss.Close()
		// Check dbs, for test/local dev.
		if pg, ok := ss.DB.(*postgres.PostgresAdapter); ok {
			pg.CheckPost()
		}
		if pg, ok := ss.Cache.(*redis.RedisAdapter); ok {
			pg.CheckRedis()
		}
		ss.StartGRPCServer()
	}
}
