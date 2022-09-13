package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"os/signal"
	"time"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "bot-token", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
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
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "main",
					Description: "Select your primary position",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "secondary",
					Description: "Select your secondary position",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "fill",
					Description: "Would you like to fill? true/false",
					Required:    true,
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
			Name: "lobby",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Create a lobby",
		},
		{
			Name: "close",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Close your current lobby",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"lobby": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf("The lobby has been created."),
				},
			})
			newChan, _ := s.GuildChannelCreate(i.GuildID, "Test", 2)
			s.ChannelMessageSend(newChan.ID, fmt.Sprintf("%s", i.Member.Mention()))
		},
		"close": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Your lobby will be deleted in 10 seconds @here."),
				},
			})
			time.Sleep(10 * time.Second)
			s.ChannelDelete(i.ChannelID)
		},
		"queue": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			margs := make([]interface{}, 0, len(options))
			msgformat := "You have entered the queue for the following positions:\n"

			// Get the value from the option map.
			// When the option exists, ok = true
			if option, ok := optionMap["gamemode"]; ok {
				margs = append(margs, option.StringValue())
				msgformat += "> Position 1: %s\n"
			}
			if option, ok := optionMap["main"]; ok {
				margs = append(margs, option.StringValue())
				msgformat += "> Position 1: %s\n"
			}
			if option, ok := optionMap["secondary"]; ok {
				margs = append(margs, option.StringValue())
				msgformat += "> Position 2: %s\n"
			}
			if opt, ok := optionMap["fill"]; ok {
				margs = append(margs, opt.BoolValue())
				msgformat += "> Fill: %v\n"
			}

			queueErr := QueueAdd()
			if queueErr != "nil" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					// Ignore type for now, they will be discussed in "responses"
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: fmt.Sprintf("Unable to add you to the queue at this time. Try again later."),
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					// Ignore type for now, they will be discussed in "responses"
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags: discordgo.MessageFlagsEphemeral,
						Content: fmt.Sprintf(
							msgformat,
							margs...,
						),
					},
				})
			}

			//log.Println("command received - waiting ten seconds")
			//time.Sleep(10 * time.Second)
			//log.Println("10 seconds is up")
		},

		//setup command to allow users to set up their profile from the slash command
		//Need to find out where the user information is stored within the command to pass
		//along to the db and store the inforamtion there.
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
				// Ignore type for now, they will be discussed in "responses"
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

	////establish connection to the Queue DB
	log.Println("Attempting to establish connection to the database...")
	conn, conErr := pgx.Connect(context.Background(), os.Getenv("postgresql://postgres:QnUrdUL30lj4lzincW9R@containers-us-west-86.railway.app:6308/railway"))
	if conErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", conErr)
		os.Exit(1)
	}
	if conErr == nil {
		fmt.Fprintf(os.Stderr, "Connection established!")
	}
	defer conn.Close(context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
