package queue

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
	"github.com/nitishm/go-rejson"
	"log"
	"os"
)

var ctx = context.Background()
var instance *redis.Client
var rh *rejson.Handler

type Command struct {
	Gamemode, Main, Secondary, Fill string
}

func Add(i *discordgo.InteractionCreate) error {
	var args Command
	err := mapstructure.Decode(ParseSlashCommand(i), &args)
	jsonArgs, _ := json.Marshal(args)
	instance.Set(i.Member.User.ID, jsonArgs, -1)
	byteArray, _ := instance.Get(i.Member.User.ID).Bytes()
	json.Unmarshal(byteArray, args)
	log.Printf("%#v", args)
	return err
}

func Connect() *redis.Client {
	if instance == nil {
		log.Println("Attempting to connect to the Redis DB...")
		rh = rejson.NewReJSONHandler()
		instance = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_DB_HOST") + ":" + os.Getenv("REDIS_DB_PORT"),
			Password: os.Getenv("REDIS_DB_PASSWORD"), // no password set
			DB:       0,                              // use default DB
		})
		rh.SetGoRedisClient(instance)
		log.Println("Connected to the Redis DB!")
	} else {
		log.Println("You are already connected to Redis DB")
		log.Println(instance)
		log.Println(instance != nil)
	}
	return instance
}

func ParseSlashCommand(i *discordgo.InteractionCreate) map[string]interface{} {
	var options = make(map[string]interface{})
	for _, option := range i.ApplicationCommandData().Options {
		options[option.Name] = option.Value
	}

	return options
}
