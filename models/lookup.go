package models

type LookupGet struct {
	Username string `json:"username"`
	Server   string `json:"server"`
}

type LookupResponse struct {
	Username      string            `json:"username"`
	Tier          string            `json:"tier"`
	Rank          string            `json:"rank"`
	Level         int               `json:"level"`
	ProfileIconId int               `json:"profileIconId"`
	Champions     []ChampionMastery `json:"champions"`
	Wins          int               `json:"wins"`
	Losses        int               `json:"losses"`
}

type ChampionMastery struct {
	ChampionName                 string
	ChampionId                   int     `json:"championId"`
	ChampionLevel                int     `json:"championLevel"`
	ChampionPoints               float64 `json:"championPoints"`
	LastPlayTime                 int64   `json:"lastPlayTime"`
	ChampionPointsSinceLastLevel int     `json:"championPointsSinceLastLevel"`
	ChampionPointsUntilNextLevel int     `json:"championPointsUntilNextLevel"`
	ChestGranted                 bool    `json:"chestGranted"`
	TokensEarned                 int     `json:"tokensEarned"`
	SummonerId                   string  `json:"summonerId"`
}
