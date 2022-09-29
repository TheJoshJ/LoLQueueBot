package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func Setup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	
	margs := make([]interface{}, 0, len(options))
	msgformat := "You have updated your profile! " +
		"Here is the information that you entered:\n"

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
}
