package models

type Profile struct {
	DiscordID string `json: "discordid"`
	Username  string `json: "username"`
	Server    string `json: "server"`
}
