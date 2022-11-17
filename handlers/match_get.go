package handlers

import (
	"discord-test/models"
	"encoding/json"
	"log"
	"net/http"
)

func MatchLookup(profile models.LookupGet) []models.MatchDataResp {
	r, err := http.Get("https://api.lolqueue.com/match/" + profile.Server + "/" + profile.Username)
	if err != nil {
		log.Printf("fatal err 1 %v", err)
	}
	var response []models.MatchDataResp
	if r.StatusCode != 404 {
		err := json.NewDecoder(r.Body).Decode(&response)
		if err != nil {
			log.Printf("error decoding response into []models.Participants match_get.go \n%v", err)
		}
	}
	return response
}
