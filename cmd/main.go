package main

import (
	"github.com/travis-james/DBCache/internal/datastore/postgres"
	"github.com/travis-james/DBCache/internal/datastore/redis"
)

func main() {
	redis.CheckRedis()
	postgres.CheckPost()
}
