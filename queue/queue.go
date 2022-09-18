package queue

import (
	"context"
	"encoding/json"
	"fmt"
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
	Gamemode, Primary, Secondary, Fill string
}

var queues = []string{"normal", "aram", "special"}

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

	instance.FlushAll()

	for _, queue := range queues {
		instance.RPush(queue, "Startup")
		log.Printf("%s queue created.", queue)
	}
	for _, queue := range queues {
		instance.LRem(queue, -1, "Startup")
	}

	return instance
}

func CommandConvert(i *discordgo.InteractionCreate) Command {
	var args Command
	mapstructure.Decode(ParseSlashCommand(i), &args)
	return args
}

func Add(i *discordgo.InteractionCreate, args Command) {

	status := instance.Get(i.Member.User.ID)

	if status != nil {
		instance.Del(i.Member.User.ID)
		for _, queue := range queues {
			instance.LRem(queue, -1, i.Member.User.ID)
		}
	}

	for _, queue := range queues {
		if args.Gamemode == queue {
			instance.RPush(queue, i.Member.User.ID)
		}
		if args.Gamemode == "normal" {
			jsonArgs, _ := json.Marshal(args)
			instance.Set(i.Member.User.ID, jsonArgs, -1)
		}
	}
}

func Position(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var found bool
	var pos int
	for _, queue := range queues {
		value := instance.LRange(queue, 0, -1)
		pos++
		log.Println(value.String())
		if value.String() == string(i.Member.User.ID) {
			found = true
			break
		}

	}
	if found {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("You are currently position %v in the queue", pos),
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("You are currently not in the queue."),
			},
		})
	}
}

func Leave(s *discordgo.Session, i *discordgo.InteractionCreate) {
	instance.Del(i.Member.User.ID)
	for _, queue := range queues {
		instance.LRem(queue, -1, i.Member.User.ID)
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("You have left your position in the queue."),
		},
	})
}

func ParseSlashCommand(i *discordgo.InteractionCreate) map[string]interface{} {
	var options = make(map[string]interface{})
	for _, option := range i.ApplicationCommandData().Options {
		options[option.Name] = option.Value
	}
	return options
}

func CheckCommand(cmd Command) bool {
	var allowed bool = true
	log.Printf("%#v", cmd)
	if cmd.Primary == cmd.Secondary && cmd.Fill == "no" {
		allowed = false
	}
	return allowed
}

func CommandResponse(allowed bool, i *discordgo.InteractionCreate, s *discordgo.Session) {

	status := instance.Get(i.Member.User.ID)
	if status != nil {
		instance.Del()
	}

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	msgformat := "You have entered the queue for the following positions:\n"
	margs := make([]interface{}, 0, len(options))

	if option, ok := optionMap["gamemode"]; ok {
		margs = append(margs, option.StringValue())
		msgformat += "> Gamemode: %s\n"
	}
	if option, ok := optionMap["primary"]; ok {
		margs = append(margs, option.StringValue())
		msgformat += "> Primary: %s\n"
	}
	if option, ok := optionMap["secondary"]; ok {
		margs = append(margs, option.StringValue())
		msgformat += "> Secondary: %s\n"
	}
	if opt, ok := optionMap["fill"]; ok {
		margs = append(margs, opt.StringValue())
		msgformat += "> Fill: %v\n"
	}
	if allowed {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf(
					msgformat,
					margs...,
				),
			},
		})

		Add(i, CommandConvert(i))
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("Please ensure that primary and secondary roles are different if you are not willing to fill."),
			},
		})
	}
}
