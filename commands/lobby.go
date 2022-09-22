package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func CreateLobby(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("The lobby has been created."),
		},
	})
	newChan, _ := s.GuildChannelCreate(i.GuildID, "Test", 2)
	s.ChannelMessageSend(newChan.ID, fmt.Sprintf("%s", i.Member.Mention()))
}
