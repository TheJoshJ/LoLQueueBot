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

	//initializes everything associated with the discord bot
	initializers.DiscordConnect()
	initializers.DiscordAddHandlers(variables.CommandHandlers)
	initializers.DiscordCreateSession()
	initializers.DiscordAddCommands(variables.Commands)

	//remove commands
	initializers.DiscordRemoveCommands()

	log.Println("Gracefully shutting down.")
}
