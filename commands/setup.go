package commands

import (
	"discord-test/handlers"
	"discord-test/models"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mitchellh/mapstructure"
	"log"
	"strings"
)

func Setup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var cmdErr []string

	//process the command within the interaction
	var options = make(map[string]interface{})
	for _, option := range i.ApplicationCommandData().Options {
		options[option.Name] = option.Value
	}

	//convert it to match the profile struct
	var profile models.Profile
	err := mapstructure.Decode(options, &profile)
	if err != nil {
		log.Fatal(err)
	}
	profile.DiscordID = i.Member.User.ID

	//send the post request
	handlers.Setup(profile)

	//check the response
	//if response.StatusCode != 200 {
	//	cmdErr = append(cmdErr, "Error posting the command to the API layer")
	//}

	//convert the slice of err to a string
	cmdString := strings.Join(cmdErr[:], ">\n")

	//reply accordingly
	if cmdErr != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf(cmdString),
			},
		})
	} else {
		{
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf("You have updated your profile!"),
				},
			})

		}
	}
}
