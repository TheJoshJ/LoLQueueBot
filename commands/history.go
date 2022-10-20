package commands

import (
	"discord-test/handlers"
	"discord-test/models"
	"github.com/bwmarrin/discordgo"
	"github.com/mitchellh/mapstructure"
	"log"
	"strconv"
)

func Match(s *discordgo.Session, i *discordgo.InteractionCreate) {
	matchHistory := make([]models.Participants, 10)
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
	matchHistory = handlers.MatchLookup(params)
	summoner = handlers.ProfileLookup(params)

	//respond to the initial lookup message
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       summoner.Username + "  -  " + strconv.Itoa(summoner.Level),
					Description: summoner.Tier + " " + summoner.Rank + " - (" + strconv.Itoa(summoner.Wins) + "/" + strconv.Itoa(summoner.Losses+summoner.Wins) + ")",
					Color:       0xffae00,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL:    "https://ddragon.leagueoflegends.com/cdn/12.20.1/img/profileicon/" + strconv.Itoa(summoner.ProfileIconId) + ".png",
						Width:  100,
						Height: 100},
					Fields: []*discordgo.MessageEmbedField{
						{Name: "\u200B", Value: "\u200B"},
						{Name: getResult(matchHistory[0]) + " - " + matchHistory[0].ChampionName + " - " + matchHistory[0].GameMode,
							Value: matchHistory[0].ChampionName + " " + strconv.Itoa(matchHistory[0].Kills) + "/" + strconv.Itoa(matchHistory[0].Deaths) + "/" + strconv.Itoa(matchHistory[0].Assists)},
						{Name: getResult(matchHistory[1]) + " - " + matchHistory[1].ChampionName + " - " + matchHistory[1].GameMode,
							Value: matchHistory[1].ChampionName + " " + strconv.Itoa(matchHistory[1].Kills) + "/" + strconv.Itoa(matchHistory[1].Deaths) + "/" + strconv.Itoa(matchHistory[1].Assists)},
						{Name: getResult(matchHistory[2]) + " - " + matchHistory[2].ChampionName + " - " + matchHistory[2].GameMode,
							Value: matchHistory[2].ChampionName + " " + strconv.Itoa(matchHistory[2].Kills) + "/" + strconv.Itoa(matchHistory[2].Deaths) + "/" + strconv.Itoa(matchHistory[2].Assists)},
						{Name: getResult(matchHistory[3]) + " - " + matchHistory[3].ChampionName + " - " + matchHistory[2].GameMode,
							Value: matchHistory[3].ChampionName + " " + strconv.Itoa(matchHistory[3].Kills) + "/" + strconv.Itoa(matchHistory[3].Deaths) + "/" + strconv.Itoa(matchHistory[3].Assists)},
						{Name: getResult(matchHistory[4]) + " - " + matchHistory[4].ChampionName + " - " + matchHistory[4].GameMode,
							Value: matchHistory[4].ChampionName + " " + strconv.Itoa(matchHistory[4].Kills) + "/" + strconv.Itoa(matchHistory[4].Deaths) + "/" + strconv.Itoa(matchHistory[4].Assists)},
						{Name: getResult(matchHistory[5]) + " - " + matchHistory[5].ChampionName + " - " + matchHistory[5].GameMode,
							Value: matchHistory[5].ChampionName + " " + strconv.Itoa(matchHistory[5].Kills) + "/" + strconv.Itoa(matchHistory[5].Deaths) + "/" + strconv.Itoa(matchHistory[5].Assists)},
						{Name: getResult(matchHistory[6]) + " - " + matchHistory[6].ChampionName + " - " + matchHistory[6].GameMode,
							Value: matchHistory[6].ChampionName + " " + strconv.Itoa(matchHistory[6].Kills) + "/" + strconv.Itoa(matchHistory[6].Deaths) + "/" + strconv.Itoa(matchHistory[6].Assists)},
						{Name: getResult(matchHistory[7]) + " - " + matchHistory[7].ChampionName + " - " + matchHistory[7].GameMode,
							Value: matchHistory[7].ChampionName + " " + strconv.Itoa(matchHistory[7].Kills) + "/" + strconv.Itoa(matchHistory[7].Deaths) + "/" + strconv.Itoa(matchHistory[7].Assists)},
						{Name: getResult(matchHistory[8]) + " - " + matchHistory[8].ChampionName + " - " + matchHistory[8].GameMode,
							Value: matchHistory[8].ChampionName + " " + strconv.Itoa(matchHistory[8].Kills) + "/" + strconv.Itoa(matchHistory[8].Deaths) + "/" + strconv.Itoa(matchHistory[8].Assists)},
						{Name: getResult(matchHistory[9]) + " - " + matchHistory[9].ChampionName + " - " + matchHistory[9].GameMode,
							Value: matchHistory[9].ChampionName + " " + strconv.Itoa(matchHistory[9].Kills) + "/" + strconv.Itoa(matchHistory[9].Deaths) + "/" + strconv.Itoa(matchHistory[9].Assists)},
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

func getResult(participants models.Participants) string {
	if participants.Win == true {
		return "Win"
	}
	return "Loss"
}
