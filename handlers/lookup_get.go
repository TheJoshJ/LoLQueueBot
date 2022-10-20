package handlers

import (
	"discord-test/models"
	"encoding/json"
	"log"
	"net/http"
)

func ProfileLookup(profile models.LookupGet) models.LookupResponse {
	r, err := http.Get("https://api.lolqueue.com/lookup/" + profile.Server + "/" + profile.Username)
	if err != nil {
		log.Printf("fatal err 1 %v", err)
	}
	var response models.LookupResponse
	if r.StatusCode != 404 {
		err := json.NewDecoder(r.Body).Decode(&response)
		if err != nil {
			log.Printf("error decoding response into rankedArray \n%v", err)
		}
	}
	return response
}
