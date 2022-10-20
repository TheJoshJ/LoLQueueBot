package commands

import (
	"discord-test/initializers"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func Leave(s *discordgo.Session, i *discordgo.InteractionCreate) {
	initializers.Instance.Dele(i.Member.User.ID)
	for _, queue := range initializers.Queues {
		initializers.Instance.LRem(queue, -1, i.Member.User.ID)
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("You have left your position in the queue."),
		},
	})
}
