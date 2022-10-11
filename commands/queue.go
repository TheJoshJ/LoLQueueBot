package commands

import (
	"discord-test/handlers"
	"discord-test/models"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mitchellh/mapstructure"
	"strings"
)

func Queue(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var cmdErr []string

	//process the command within the interaction
	var options = make(map[string]interface{})
	for _, option := range i.ApplicationCommandData().Options {
		options[option.Name] = option.Value
	}

	//convert it to match the Command struct
	var args models.Command
	mapstructure.Decode(options, &args)
	args.DiscordID = i.Member.User.ID

	//check to see if the command is valid
	if args.Primary == args.Secondary && args.Fill == "no" {
		cmdErr = append(cmdErr, "> Please ensure that primary and secondary roles are different if you are not willing to fill.")
	}

	//send the post request
	response := handlers.Queue(args)

	//check the response
	if response.StatusCode != 200 {
		cmdErr = append(cmdErr, "> Error posting the command to the API layer")
	}

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
					Content: fmt.Sprintf("You have queued up with the following arguments\n> Gamemode: %s\n> Primary: %s\n> Secondary: %s\n> Fill: %v\n", args.Gamemode, args.Primary, args.Secondary, args.Fill),
				},
			})

		}
	}
}
