package main

import (
	"context"
	"discord-test/queue"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"os/signal"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", os.Getenv("DISCORD_TOKEN"), "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", false, "Remove all commands after shutting down or not")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "queue",
			Description: "Queue up to be put with a team",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "gamemode",
					Description: "select which gamemode you are wanting to play",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "normal",
							Value: "normal",
						},
						{
							Name:  "aram",
							Value: "aram",
						},
						{
							Name:  "rotating",
							Value: "rotating",
						},
					},
					Required: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "primary",
					Description: "Select your primary position",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "top",
							Value: "top",
						},
						{
							Name:  "jungle",
							Value: "jungle",
						},
						{
							Name:  "mid",
							Value: "mid",
						},
						{
							Name:  "bot",
							Value: "bot",
						},
						{
							Name:  "support",
							Value: "support",
						},
					},
					Required: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "secondary",
					Description: "Select your secondary position",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "top",
							Value: "top",
						},
						{
							Name:  "jungle",
							Value: "jungle",
						},
						{
							Name:  "mid",
							Value: "mid",
						},
						{
							Name:  "bot",
							Value: "bot",
						},
						{
							Name:  "support",
							Value: "support",
						},
					},
					Required: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "fill",
					Description: "Would you like to fill? true/false",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "yes",
							Value: "yes",
						},
						{
							Name:  "no",
							Value: "no",
						},
					},
					Required: true,
				},
			},
		},
		{
			Name:        "setup",
			Description: "Set up your profile before you queue up",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "String option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "server",
					Description: "String option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "rank",
					Description: "String option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "position",
					Description: "String option",
					Required:    false,
				},
			},
		},
		{
			Name:        "lobby",
			Description: "Create a lobby",
		},
		{
			Name:        "close",
			Description: "Close your current lobby",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"lobby": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			queue.CreateLobby(s, i)
		},
		"close": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			queue.CloseLobby(s, i)
		},
		"queue": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			queue.CommandResponse(queue.CheckCommand(queue.CommandConvert(i)), i, s)
		},
		"setup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			// This example stores the provided arguments in an []interface{}
			// which will be used to format the bot's response
			margs := make([]interface{}, 0, len(options))
			msgformat := "You have updated your profile! " +
				"Here is the information that you entered:\n"

			// Get the value from the option map.
			// When the option exists, ok = true
			if option, ok := optionMap["username"]; ok {
				margs = append(margs, option.StringValue())
				msgformat += "> Username: %s\n"
			}
			if option, ok := optionMap["server"]; ok {
				margs = append(margs, option.StringValue())
				msgformat += "> Server: %s\n"
			}
			if option, ok := optionMap["rank"]; ok {
				margs = append(margs, option.StringValue())
				msgformat += "> Rank: %s\n"
			}
			if option, ok := optionMap["position"]; ok {
				margs = append(margs, option.StringValue())
				msgformat += "> Rank: %s\n"
			}

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

		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}
func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()
	log.Println("Commands successfully added!")

	//establish connection to the PostgreSQL DB
	log.Println("Attempting to establish connection to the PostgreSQL database...")
	conn, conErr := pgx.Connect(context.Background(), os.Getenv("POSTGRES_DB_URL"))
	if conErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", conErr)
		os.Exit(1)
	}
	if conErr == nil {
		fmt.Fprintf(os.Stderr, "Connection established!\n")
	}
	defer conn.Close(context.Background())

	queue.Connect()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")

		registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		if err != nil {
			log.Fatalf("Could not fetch registered commands: %v", err)
		}

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
