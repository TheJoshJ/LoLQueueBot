package queue

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

type NormalTeam struct {
	top, jg, mid, bot, sup string
}
type AramTeam struct {
	one, two, three, four, five string
}

//func MatchNormal() {
//
//}

func MatchAram() {
	//if instance.LLen("aram").Val() >= 5 {
	queue, err := instance.Get("aram").Bytes()
	if err != nil {
		log.Printf("there was an error: %s", err)
	} else {
		for idx, uid := range queue {
			log.Println(idx, uid)
		}
	}
}

//}

//func MatchSpecial() {
//
//}

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
