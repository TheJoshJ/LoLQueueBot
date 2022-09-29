package main

import (
	"discord-test/initializers"
	"discord-test/variables"
	"flag"
	"github.com/bwmarrin/discordgo"
	"log"
)

var s *discordgo.Session

func init() {
	flag.Parse()
}

func main() {

	//initializes all the connections
	initializers.RedisCreateConnection()
	go initializers.CreateGinConnection()

	initializers.DiscordConnect()
	initializers.DiscordAddHandlers(variables.CommandHandlers)
	initializers.DiscordCreateSession()
	initializers.DiscordAddCommands(variables.Commands)

	//remove commands
	initializers.DiscordRemoveCommands()

	log.Println("Gracefully shutting down.")
}
