package variables

import (
	"discord-test/commands"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	Commands = []*discordgo.ApplicationCommand{
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
		{
			Name:        "leave",
			Description: "Leave your current position in queue",
		},
		{
			Name:        "pos",
			Description: "Check your current position in queue",
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"lobby": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			commands.CreateLobby(s, i)
		},
		"close": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			commands.CloseLobby(s, i)
		},
		"queue": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			commands.Queue(s, i)
		},
		"leave": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			commands.Leave(s, i)
		},
		"pos": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			//queue.Position(s, i)
			commands.Position()
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
