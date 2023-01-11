package main

import (
	"bytes"
	"discord-test/handlers"
	"discord-test/models"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var s *discordgo.Session
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", os.Getenv("DISCORD_TOKEN"), "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", false, "Remove all commands after shutting down or not")
)
var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "update",
			Description: "Update your profile",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "Your display name in League of Legends",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "server",
					Description: "select which server you play on",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "NA",
							Value: "NA",
						},
						{
							Name:  "EUNE",
							Value: "EUNE",
						},
						{
							Name:  "EUW",
							Value: "EUW",
						},
						{
							Name:  "LAN",
							Value: "LAN",
						},
						{
							Name:  "LAS",
							Value: "LAS",
						},
						{
							Name:  "OCE",
							Value: "OCE",
						},
						{
							Name:  "BR",
							Value: "BR",
						},
						{
							Name:  "JP",
							Value: "JP",
						},
						{
							Name:  "KR",
							Value: "KR",
						},
						{
							Name:  "TR",
							Value: "TR",
						},
						{
							Name:  "RU",
							Value: "RU",
						},
					},
					Required: true,
				},
			},
		},
		{
			Name:        "setup",
			Description: "Set up your profile to track your stats",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "Your display name in League of Legends",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "server",
					Description: "select which server you play on",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "NA",
							Value: "NA",
						},
						{
							Name:  "EUNE",
							Value: "EUNE",
						},
						{
							Name:  "EUW",
							Value: "EUW",
						},
						{
							Name:  "LAN",
							Value: "LAN",
						},
						{
							Name:  "LAS",
							Value: "LAS",
						},
						{
							Name:  "OCE",
							Value: "OCE",
						},
						{
							Name:  "BR",
							Value: "BR",
						},
						{
							Name:  "JP",
							Value: "JP",
						},
						{
							Name:  "KR",
							Value: "KR",
						},
						{
							Name:  "TR",
							Value: "TR",
						},
						{
							Name:  "RU",
							Value: "RU",
						},
					},
					Required: true,
				},
			},
		},
		{
			Name:        "profile",
			Description: "View your current profile",
		},
		{
			Name:        "retire",
			Description: "Remove your profile from the servers leaderboard",
		},
		{
			Name:        "lookup",
			Description: "Look up a users League of Legends profile",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "The desired users display name in League",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "server",
					Description: "The desired users server",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "NA",
							Value: "NA",
						},
						{
							Name:  "EUNE",
							Value: "EUNE",
						},
						{
							Name:  "EUW",
							Value: "EUW",
						},
						{
							Name:  "LAN",
							Value: "LAN",
						},
						{
							Name:  "LAS",
							Value: "LAS",
						},
						{
							Name:  "OCE",
							Value: "OCE",
						},
						{
							Name:  "BR",
							Value: "BR",
						},
						{
							Name:  "JP",
							Value: "JP",
						},
						{
							Name:  "KR",
							Value: "KR",
						},
						{
							Name:  "TR",
							Value: "TR",
						},
						{
							Name:  "RU",
							Value: "RU",
						},
					},
					Required: true,
				},
			},
		},
		{
			Name:        "match",
			Description: "Look up a users League of Legends match history",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "The desired users display name in League",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "server",
					Description: "The desired users server",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "NA",
							Value: "NA",
						},
						{
							Name:  "EUNE",
							Value: "EUNE",
						},
						{
							Name:  "EUW",
							Value: "EUW",
						},
						{
							Name:  "LAN",
							Value: "LAN",
						},
						{
							Name:  "LAS",
							Value: "LAS",
						},
						{
							Name:  "OCE",
							Value: "OCE",
						},
						{
							Name:  "BR",
							Value: "BR",
						},
						{
							Name:  "JP",
							Value: "JP",
						},
						{
							Name:  "KR",
							Value: "KR",
						},
						{
							Name:  "TR",
							Value: "TR",
						},
						{
							Name:  "RU",
							Value: "RU",
						},
					},
					Required: true,
				},
			},
		},
		{
			Name:        "lobby",
			Description: "Create a lobby",
		},
		{
			Name:        "close",
			Description: "Close your current lobby",
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"lobby": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			CreateLobby(s, i)
		},
		"close": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			CloseLobby(s, i)
		},
		"queue": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Queue(s, i)
		},
		"profile": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Profile(s, i)
		},
		"setup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Setup(s, i)
		},
		"lookup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Lookup(s, i)
		},
		"match": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Match(s, i)
		},
		"update": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Update(s, i)
		},
		"retire": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Delete(s, i)
		},
	}
)

func init() {
	flag.Parse()
}

func main() {

	DiscordConnect()
	DiscordAddHandlers(CommandHandlers)
	DiscordCreateSession()
	DiscordAddCommands(Commands)

	//remove commands
	DiscordRemoveCommands()

	log.Println("Gracefully shutting down.")
}

func DiscordConnect() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func DiscordAddCommands(commands []*discordgo.ApplicationCommand) {
	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()
	log.Println("Commands successfully added!")

	//event listener
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}

func DiscordAddHandlers(commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func DiscordCreateSession() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
}

func DiscordRemoveCommands() {
	if *RemoveCommands {
		log.Println("Removing commands...")

		registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		if err != nil {
			log.Fatalf("Could not fetch registered commands: %v", err)
		}

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}

//discord bot commands
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
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Your lobby will be deleted in 10 seconds @here."),
		},
	})
	if err != nil {
		log.Println("error respoinding to /close command")
	}
	time.Sleep(10 * time.Second)
	s.ChannelDelete(i.ChannelID)
}

func Match(s *discordgo.Session, i *discordgo.InteractionCreate) {
	matchHistory := make([]models.MatchDataResp, 20)
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
					Description: summoner.Tier + " " + summoner.Rank + " - (" + strconv.Itoa(summoner.Wins) + "W/" + strconv.Itoa(summoner.Losses) + "L)",
					Color:       0xffae00,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL:    "https://ddragon.leagueoflegends.com/cdn/12.20.1/img/profileicon/" + strconv.Itoa(summoner.ProfileIconId) + ".png",
						Width:  100,
						Height: 100},
					Fields: []*discordgo.MessageEmbedField{
						{Name: "\u200B", Value: "\u200B"},
						{Name: getResult(matchHistory[0]) + " - " + matchHistory[0].GameMode,
							Value: matchHistory[0].ChampionName + " " + strconv.Itoa(matchHistory[0].Kills) + "/" + strconv.Itoa(matchHistory[0].Deaths) + "/" + strconv.Itoa(matchHistory[0].Assists)},
						{Name: getResult(matchHistory[1]) + " - " + matchHistory[1].GameMode,
							Value: matchHistory[1].ChampionName + " " + strconv.Itoa(matchHistory[1].Kills) + "/" + strconv.Itoa(matchHistory[1].Deaths) + "/" + strconv.Itoa(matchHistory[1].Assists)},
						{Name: getResult(matchHistory[2]) + " - " + matchHistory[2].GameMode,
							Value: matchHistory[2].ChampionName + " " + strconv.Itoa(matchHistory[2].Kills) + "/" + strconv.Itoa(matchHistory[2].Deaths) + "/" + strconv.Itoa(matchHistory[2].Assists)},
						{Name: getResult(matchHistory[3]) + " - " + matchHistory[2].GameMode,
							Value: matchHistory[3].ChampionName + " " + strconv.Itoa(matchHistory[3].Kills) + "/" + strconv.Itoa(matchHistory[3].Deaths) + "/" + strconv.Itoa(matchHistory[3].Assists)},
						{Name: getResult(matchHistory[4]) + " - " + matchHistory[4].GameMode,
							Value: matchHistory[4].ChampionName + " " + strconv.Itoa(matchHistory[4].Kills) + "/" + strconv.Itoa(matchHistory[4].Deaths) + "/" + strconv.Itoa(matchHistory[4].Assists)},
						{Name: getResult(matchHistory[5]) + " - " + matchHistory[5].GameMode,
							Value: matchHistory[5].ChampionName + " " + strconv.Itoa(matchHistory[5].Kills) + "/" + strconv.Itoa(matchHistory[5].Deaths) + "/" + strconv.Itoa(matchHistory[5].Assists)},
						{Name: getResult(matchHistory[6]) + " - " + matchHistory[6].GameMode,
							Value: matchHistory[6].ChampionName + " " + strconv.Itoa(matchHistory[6].Kills) + "/" + strconv.Itoa(matchHistory[6].Deaths) + "/" + strconv.Itoa(matchHistory[6].Assists)},
						{Name: getResult(matchHistory[7]) + " - " + matchHistory[7].GameMode,
							Value: matchHistory[7].ChampionName + " " + strconv.Itoa(matchHistory[7].Kills) + "/" + strconv.Itoa(matchHistory[7].Deaths) + "/" + strconv.Itoa(matchHistory[7].Assists)},
						{Name: getResult(matchHistory[8]) + " - " + matchHistory[8].GameMode,
							Value: matchHistory[8].ChampionName + " " + strconv.Itoa(matchHistory[8].Kills) + "/" + strconv.Itoa(matchHistory[8].Deaths) + "/" + strconv.Itoa(matchHistory[8].Assists)},
						{Name: getResult(matchHistory[9]) + " - " + matchHistory[9].GameMode,
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

func Lookup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	matchHistory := make([]models.MatchDataResp, 10)
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
	matchHistory = handlers.MatchLookup(params)

	//round the champion points to the nearest 1000
	for i, v := range summoner.Champions {
		champions[i].ChampionPoints = math.Round(v.ChampionPoints / 1000)
	}

	//get win/loss for the user over their last 10 games
	var winLoss []int
	winLoss = calculateWinLoss(matchHistory)

	//respond to the initial lookup message
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
						{Name: champions[0].ChampionName, Value: strconv.Itoa(int(champions[0].ChampionPoints)) + "k\nMastery " + strconv.Itoa(champions[0].ChampionLevel) + "‎ ‎ ‎ ‎ ‎ ‎ ‎ ", Inline: true},
						{Name: champions[1].ChampionName, Value: strconv.Itoa(int(champions[1].ChampionPoints)) + "k\nMastery " + strconv.Itoa(champions[1].ChampionLevel) + "‎ ‎ ‎ ‎ ‎ ‎ ‎ ", Inline: true},
						{Name: champions[2].ChampionName, Value: strconv.Itoa(int(champions[2].ChampionPoints)) + "k\nMastery " + strconv.Itoa(champions[2].ChampionLevel) + "‎ ‎ ‎ ‎ ‎ ‎ ‎ ", Inline: true},
						{Name: champions[3].ChampionName, Value: strconv.Itoa(int(champions[3].ChampionPoints)) + "k\nMastery " + strconv.Itoa(champions[3].ChampionLevel) + "‎ ‎ ‎ ‎ ‎ ‎ ‎ ", Inline: true},
						{Name: champions[4].ChampionName, Value: strconv.Itoa(int(champions[4].ChampionPoints)) + "k\nMastery " + strconv.Itoa(champions[4].ChampionLevel) + "‎ ‎ ‎ ‎ ‎ ‎ ‎ ", Inline: true},
						{Name: champions[5].ChampionName, Value: strconv.Itoa(int(champions[5].ChampionPoints)) + "k\nMastery " + strconv.Itoa(champions[5].ChampionLevel) + "‎ ‎ ‎ ‎ ‎ ‎ ‎ ", Inline: true},
						{Name: "\u200B", Value: "\u200B"},
						{Name: "Recent Matches - (" + strconv.Itoa(winLoss[0]) + "W/" + strconv.Itoa(winLoss[1]) + "L)", Value: "\u200B"},
						{Name: getResult(matchHistory[0]) + " - " + matchHistory[0].GameMode,
							Value: matchHistory[0].ChampionName + "\n" + strconv.Itoa(matchHistory[0].Kills) + "/" + strconv.Itoa(matchHistory[0].Deaths) + "/" + strconv.Itoa(matchHistory[0].Assists) + "\n", Inline: true},
						{Name: getResult(matchHistory[1]) + " - " + matchHistory[1].GameMode,
							Value: matchHistory[1].ChampionName + "\n" + strconv.Itoa(matchHistory[1].Kills) + "/" + strconv.Itoa(matchHistory[1].Deaths) + "/" + strconv.Itoa(matchHistory[1].Assists) + "\n", Inline: true},
						{Name: getResult(matchHistory[2]) + " - " + matchHistory[2].GameMode,
							Value: matchHistory[2].ChampionName + "\n" + strconv.Itoa(matchHistory[2].Kills) + "/" + strconv.Itoa(matchHistory[2].Deaths) + "/" + strconv.Itoa(matchHistory[2].Assists) + "\n", Inline: true},
						{Name: getResult(matchHistory[3]) + " - " + matchHistory[2].GameMode,
							Value: matchHistory[3].ChampionName + "\n" + strconv.Itoa(matchHistory[3].Kills) + "/" + strconv.Itoa(matchHistory[3].Deaths) + "/" + strconv.Itoa(matchHistory[3].Assists) + "\n", Inline: true},
						{Name: getResult(matchHistory[4]) + " - " + matchHistory[4].GameMode,
							Value: matchHistory[4].ChampionName + "\n" + strconv.Itoa(matchHistory[4].Kills) + "/" + strconv.Itoa(matchHistory[4].Deaths) + "/" + strconv.Itoa(matchHistory[4].Assists) + "\n", Inline: true},
						{Name: getResult(matchHistory[5]) + " - " + matchHistory[5].GameMode,
							Value: matchHistory[5].ChampionName + "\n" + strconv.Itoa(matchHistory[5].Kills) + "/" + strconv.Itoa(matchHistory[5].Deaths) + "/" + strconv.Itoa(matchHistory[5].Assists) + "\n", Inline: true},
						{Name: getResult(matchHistory[6]) + " - " + matchHistory[6].GameMode,
							Value: matchHistory[6].ChampionName + "\n" + strconv.Itoa(matchHistory[6].Kills) + "/" + strconv.Itoa(matchHistory[6].Deaths) + "/" + strconv.Itoa(matchHistory[6].Assists) + "\n", Inline: true},
						{Name: getResult(matchHistory[7]) + " - " + matchHistory[7].GameMode,
							Value: matchHistory[7].ChampionName + "\n" + strconv.Itoa(matchHistory[7].Kills) + "/" + strconv.Itoa(matchHistory[7].Deaths) + "/" + strconv.Itoa(matchHistory[7].Assists) + "\n", Inline: true},
						{Name: getResult(matchHistory[8]) + " - " + matchHistory[8].GameMode,
							Value: matchHistory[8].ChampionName + "\n" + strconv.Itoa(matchHistory[8].Kills) + "/" + strconv.Itoa(matchHistory[8].Deaths) + "/" + strconv.Itoa(matchHistory[8].Assists) + "\n", Inline: true},
						{Name: getResult(matchHistory[9]) + " - " + matchHistory[9].GameMode,
							Value: matchHistory[9].ChampionName + "\n" + strconv.Itoa(matchHistory[9].Kills) + "/" + strconv.Itoa(matchHistory[9].Deaths) + "/" + strconv.Itoa(matchHistory[9].Assists) + "\n", Inline: true},
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
	if err != nil {
		log.Printf("error responding to /lookup request, %s", err)
	}
}

func Profile(s *discordgo.Session, i *discordgo.InteractionCreate) {
	discordID := i.Member.User.ID
	log.Println(i.Member.User.ID)
	r, err := http.Get("https://api.lolqueue.com/v1/user/" + discordID)
	if err != nil {
		log.Printf("Error calling user endpoing on riot API %v", err)
	}
	log.Println(r.Body)
	var response models.UserLookupResponse
	var respString string
	if r.StatusCode == 200 {
		err := json.NewDecoder(r.Body).Decode(&response)
		log.Println(response)
		if err != nil {
			log.Printf("error decoding response into &response for /profile command \n%v", err)
			respString = "There was an error retrieving your profile information\nIf this problem persists, please inform a server admin."
		}
		if respString == "" {
			respString = fmt.Sprintf("> Servers: %v\n> LoL Username: %v\n> Region: %v", strings.Join(response.Servers, ", "), response.RiotUsername, response.RiotServer)
		}
	} else {
		respString = "Unable to find your profile.\nHave you set up an profile yet?\nUse the /setup command to create a profile"
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf(respString),
		},
	})
}

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

func Setup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var res []string

	//process the command within the interaction
	var options = make(map[string]interface{})
	for _, option := range i.ApplicationCommandData().Options {
		options[option.Name] = option.Value
	}

	//convert it to match the profile struct
	var args models.SetupCmd
	err := mapstructure.Decode(options, &args)
	if err != nil {
		log.Fatal(err)
	}
	g, _ := s.Guild(i.GuildID)
	profile := models.Profile{
		RiotUsername:      args.Username,
		RiotServer:        args.Server,
		DiscordID:         i.Member.User.ID,
		DiscordUsername:   i.Member.User.Username,
		DiscordServerID:   i.GuildID,
		DiscordServerName: g.Name,
	}

	log.Println(profile)

	//send the post request
	data, err := json.Marshal(profile)
	if err != nil {
		log.Println("error marshalling profile data.")
	}

	resp, err := http.Post("https://api.lolqueue.com/v1/user", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("error sending post request to lolqueue.com when creating a new user \n%v", err)
	}

	response := resp.StatusCode

	//check the response
	if response == 208 {
		res = append(res, "User already exists! To update your information, use '/update' instead!")
	} else if response == 500 {
		res = append(res, "Error posting the command to the API layer, contact an admin if this issue persists.")
	} else if response == 404 {
		res = append(res, "Invalid Riot Username - please double check that your username is correct and try again.")
	} else if response == 409 {
		res = append(res, "You've added your discord account to this server!\n> Your LoL account doesn't match the existing account we have a record of though, so it was not updated.\n> Please use /update to update your LoL account username or /profile to view your current profile.")
	} else if response != 201 {
		res = append(res, "Unknown response. Please contact an admin with what you did and how to recreate it.")
	} else {
		res = append(res, fmt.Sprintf("You have updated your profile!\n> IGN: %s\n> Server: %s\n> Discord Username: %s\n Use /update to update any incorrect information!", profile.RiotUsername, profile.RiotServer, profile.DiscordUsername))
	}

	//convert the slice of err to a string
	resString := strings.Join(res[:], ">\n")
	log.Println(resString)

	//reply accordingly
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf(resString),
		},
	})
}

func Update(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var options = make(map[string]interface{})
	for _, option := range i.ApplicationCommandData().Options {
		options[option.Name] = option.Value
	}

	//convert it to match the profile struct
	var args models.SetupCmd
	err := mapstructure.Decode(options, &args)
	if err != nil {
		log.Fatal(err)
	}
	update := models.ProfileUpdate{
		RiotUsername:    args.Username,
		RiotServer:      args.Server,
		DiscordID:       i.Member.User.ID,
		DiscordServerID: i.GuildID,
	}

	body, _ := json.Marshal(&update)

	req, err := http.NewRequest("PUT", "api.lolqueue.com/v1/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var response string
	if resp.StatusCode == 400 || resp.StatusCode == 500 {
		response = "Unable to execute command. If the issue persists, please alert an admin."
	} else if resp.StatusCode == 404 {
		response = "A riot username could not be found on that server. Please update and try again."
	} else if resp.StatusCode == 409 {
		response = "You have not yet created a profile! First do /setup to create a profile."
	} else if resp.StatusCode == 200 {
		response = fmt.Sprintf("You have updated your profile!\n> IGN: %s\n> Server: %s\nUse /update to update any incorrect information!", update.RiotUsername, update.RiotServer)
	} else {
		response = "Unable to execute command. If the issue persists, please alert an admin."
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf(response),
		},
	})
}

func Delete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string
	var URL string

	URL = fmt.Sprintf("api.lolqueue.com/v1/user/%v?serverid=%v", i.Member.User.ID, i.GuildID)
	req, err := http.NewRequest("DELETE", URL, nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		response = "Uh oh! Our servers seem to not want to let you retire at this moment!\nIf this issue persists, please contact an admin!"
	} else {
		response = fmt.Sprintf("%v... we are sad to see you go but nevertheless, we have processed your retirement and you are free to leave.\nHave fun in the retirement home.", i.User.Username)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: response,
		},
	})
}

//support functions
func getResult(participants models.MatchDataResp) string {
	if participants.Win == true {
		return "Win"
	}
	return "Loss"
}
func calculateWinLoss(matchHistory []models.MatchDataResp) []int {
	var win = 0
	var loss = 0
	winLoss := make([]int, 2)

	for _, match := range matchHistory {
		if match.Win == true {
			win++
		} else {
			loss++
		}

	}
	winLoss[0] = win
	winLoss[1] = loss

	return winLoss
}
