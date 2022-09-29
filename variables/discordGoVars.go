package variables

import (
	"discord-test/commands"
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
			Name:        "find",
			Description: "Find a player to complete your team.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "gamemode",
					Description: "select which gamemode you will be playing",
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
					Name:        "position",
					Description: "Select the position of the player you need.",
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
			commands.Position()
		},
		"find": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			commands.Find()
		},
		"setup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			commands.Setup(s, i)
		},
	}
)
