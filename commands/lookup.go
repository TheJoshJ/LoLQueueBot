package commands

import (
	"discord-test/handlers"
	"discord-test/models"
	"github.com/bwmarrin/discordgo"
	"github.com/mitchellh/mapstructure"
	"log"
	"math"
	"strconv"
)

func Lookup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var summoner models.LookupResponse
	var options = make(map[string]interface{})
	for _, option := range i.ApplicationCommandData().Options {
		options[option.Name] = option.Value
	}

	//convert it to match the profile struct
	var params models.LookupGet
	err := mapstructure.Decode(options, &params)
	if err != nil {
		log.Print(err)
	}

	//get the information from the API layer
	summoner = handlers.ProfileLookup(params)
	champions := summoner.Champions

	//round the champion points to the nearest 1000
	for i, v := range summoner.Champions {
		champions[i].ChampionPoints = math.Round(v.ChampionPoints / 1000)
	}

	//respond to the initial lookup message
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       summoner.Username + "  -  " + strconv.Itoa(summoner.Level),
					Description: summoner.Tier + " " + summoner.Rank + " - (" + strconv.Itoa(summoner.Wins) + "W/" + strconv.Itoa(summoner.Losses) + "L)",
					Color:       0xffae00,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL:    "https://ddragon.leagueoflegends.com/cdn/12.20.1/img/profileicon/" + strconv.Itoa(summoner.ProfileIconId) + ".png",
						Width:  100,
						Height: 100},
					Fields: []*discordgo.MessageEmbedField{
						{Name: "\u200B", Value: "\u200B"},
						{Name: champions[0].ChampionName, Value: strconv.Itoa(int(champions[0].ChampionPoints)) + "k\nMastery " + strconv.Itoa(champions[0].ChampionLevel) + "‎ ‎ ‎ ‎ ‎ ‎ ‎ ", Inline: true},
						{Name: champions[1].ChampionName, Value: strconv.Itoa(int(champions[1].ChampionPoints)) + "k\nMastery " + strconv.Itoa(champions[1].ChampionLevel) + "‎ ‎ ‎ ‎ ‎ ‎ ‎ ", Inline: true},
						{Name: champions[2].ChampionName, Value: strconv.Itoa(int(champions[2].ChampionPoints)) + "k\nMastery " + strconv.Itoa(champions[2].ChampionLevel) + "‎ ‎ ‎ ‎ ‎ ‎ ‎ ", Inline: true},
						{Name: champions[3].ChampionName, Value: strconv.Itoa(int(champions[3].ChampionPoints)) + "k\nMastery " + strconv.Itoa(champions[3].ChampionLevel) + "‎ ‎ ‎ ‎ ‎ ‎ ‎ ", Inline: true},
						{Name: champions[4].ChampionName, Value: strconv.Itoa(int(champions[4].ChampionPoints)) + "k\nMastery " + strconv.Itoa(champions[4].ChampionLevel) + "‎ ‎ ‎ ‎ ‎ ‎ ‎ ", Inline: true},
						{Name: champions[5].ChampionName, Value: strconv.Itoa(int(champions[5].ChampionPoints)) + "k\nMastery " + strconv.Itoa(champions[5].ChampionLevel) + "‎ ‎ ‎ ‎ ‎ ‎ ‎ ", Inline: true},
					},
					Image: &discordgo.MessageEmbedImage{
						URL:    "https://cdn.discordapp.com/attachments/1019324333098803340/1032455422474469407/LoLQueue.com_2.png",
						Width:  2000,
						Height: 80,
					},
				},
			},
		},
	})
}
