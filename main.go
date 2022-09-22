package main

import (
	"discord-test/initializers"
	"discord-test/variables"
	"flag"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
)

var s *discordgo.Session

func init() {
	flag.Parse()
	initializers.DiscordConnect()
	initializers.DiscordAddHandlers(variables.CommandHandlers)
}
func main() {
	//initializers all the connections
	initializers.DiscordCreateSession()
	initializers.DiscordAddCommands(variables.Commands)
	initializers.PostgresCreateConnection()
	initializers.RedisCreateConnection()
	initializers.CreateGinConnection()

	//event listener
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	//remove commands
	initializers.DiscordRemoveCommands()

	log.Println("Gracefully shutting down.")
}
