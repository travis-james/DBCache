package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	redis "github.com/redis/go-redis/v9"
)

func main() {
	checkRedis()
	checkPost()
}

func checkRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
	})
	ctx := context.Background()

	err := client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val)
}

func checkPost() {
	var (
		host     = "localhost"
		port     = 5432
		user     = "local_dev"
		password = "local_pw"
		dbname   = "local_db"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}
