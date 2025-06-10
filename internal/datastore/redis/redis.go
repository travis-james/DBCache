package redis

import (
	"context"
	"fmt"
	"os"

	redis "github.com/redis/go-redis/v9"
)

func CheckRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PW"),
		DB:       0, // Use default DB for now.
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

	val, err = client.Get(ctx, "key3").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key3", val)
}
