package initializers

import (
	"github.com/go-redis/redis"
	"github.com/nitishm/go-rejson"
	"log"
	"os"
)

var Instance *redis.Client
var rh *rejson.Handler

//avaiable queues within Redis
var Queues = []string{"normal", "aram", "rotating"}

func RedisCreateConnection() *redis.Client {
	//check to see if there is already a connection. If not, connect
	if Instance == nil {
		log.Println("Attempting to connect to the Redis DB...")
		rh = rejson.NewReJSONHandler()
		Instance = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_DB_HOST") + ":" + os.Getenv("REDIS_DB_PORT"),
			Password: os.Getenv("REDIS_DB_PASSWORD"), // no password set
			DB:       0,                              // use default DB
		})
		rh.SetGoRedisClient(Instance)
		log.Println("Connected to the Redis DB!")
		//if you're already connected, don't connect
	} else {
		log.Println("You are already connected to Redis DB")
		log.Println(Instance)
		log.Println(Instance != nil)
	}
	//flush the DB of all existing data
	Instance.FlushAll()
	//return the instance for other functions to use
	return Instance
}
