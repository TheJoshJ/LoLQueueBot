package commands

import (
	"discord-test/initializers"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mitchellh/mapstructure"
	"log"
)

type Command struct {
	Gamemode  string `json: "gamemode"`
	Primary   string `json: "primary"`
	Secondary string `json: "secondary"`
	Fill      string `json: "fill"`
}

func Queue(s *discordgo.Session, i *discordgo.InteractionCreate) {

	//process the command within the interaction
	var options = make(map[string]interface{})
	for _, option := range i.ApplicationCommandData().Options {
		options[option.Name] = option.Value
	}

	//convert it to match the Command struct
	var args Command
	mapstructure.Decode(options, &args)

	//check to see if the command is valid
	var allowed bool = true
	log.Printf("%#v", args)
	if args.Primary == args.Secondary && args.Fill == "no" {
		allowed = false
	}

	//respond to the command
	status := initializers.Instance.Get(i.Member.User.ID)
	if status != nil {
		initializers.Instance.Del()
	}

	cmdOptions := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(cmdOptions))
	for _, opt := range cmdOptions {
		optionMap[opt.Name] = opt
	}
	msgformat := "You have entered the queue for the following positions:\n"
	margs := make([]interface{}, 0, len(cmdOptions))

	if option, ok := optionMap["gamemode"]; ok {
		margs = append(margs, option.StringValue())
		msgformat += "> Gamemode: %s\n"
	}
	if option, ok := optionMap["primary"]; ok {
		margs = append(margs, option.StringValue())
		msgformat += "> Primary: %s\n"
	}
	if option, ok := optionMap["secondary"]; ok {
		margs = append(margs, option.StringValue())
		msgformat += "> Secondary: %s\n"
	}
	if opt, ok := optionMap["fill"]; ok {
		margs = append(margs, opt.StringValue())
		msgformat += "> Fill: %v\n"
	}
	if allowed {
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

	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("Please ensure that primary and secondary roles are different if you are not willing to fill."),
			},
		})
	}

	//add the queue to the database
	if status != nil {
		initializers.Instance.Del(i.Member.User.ID)
		for _, queue := range initializers.Queues {
			initializers.Instance.LRem(queue, -1, i.Member.User.ID)
		}
	}

	for _, queue := range initializers.Queues {
		if args.Gamemode == queue {
			initializers.Instance.RPush(queue, i.Member.User.ID)
		}
		if args.Gamemode == "normal" {
			jsonArgs, _ := json.Marshal(args)
			initializers.Instance.Set(i.Member.User.ID, jsonArgs, -1)
		}
	}
}
