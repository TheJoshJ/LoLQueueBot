package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

var ctx = context.Background()

func Connect() {
	log.Println("Attempting to connect to the Redis DB...")
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_DB_HOST") + ":" + os.Getenv("REDIS_DB_PORT"),
		Password: os.Getenv("REDIS_DB_PASSWORD"), // no password set
		DB:       0,                              // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	log.Println("Connected to the Redis DB!")
}
