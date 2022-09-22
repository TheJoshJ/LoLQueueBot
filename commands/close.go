package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

func CloseLobby(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Your lobby will be deleted in 10 seconds @here."),
		},
	})
	time.Sleep(10 * time.Second)
	s.ChannelDelete(i.ChannelID)
}
