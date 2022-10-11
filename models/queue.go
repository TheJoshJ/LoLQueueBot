package models

type Command struct {
	DiscordID string `json: "id"`
	Gamemode  string `json: "gamemode"`
	Primary   string `json: "primary"`
	Secondary string `json: "secondary"`
	Fill      string `json: "fill"`
}
