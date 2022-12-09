package handlers

import (
	"bytes"
	"discord-test/models"
	"encoding/json"
	"log"
	"net/http"
)

var (
	client = http.DefaultClient
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

func MatchLookup(profile models.LookupGet) []models.MatchDataResp {
	r, err := http.Get("https://api.lolqueue.com/match/" + profile.Server + "/" + profile.Username)
	if err != nil {
		log.Printf("fatal err 1 %v", err)
	}
	var response []models.MatchDataResp
	if r.StatusCode != 404 {
		err := json.NewDecoder(r.Body).Decode(&response)
		if err != nil {
			log.Printf("error decoding response into []Participants match_get.go \n%v", err)
		}
	}
	return response
}

func Ping() {
	var output map[string]string

	r, err := client.Get("https://api.lolqueue.com/ping")
	if err != nil {
		log.Fatalf("Post failed %v", err)
	}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&output)
	log.Println(output["message"])
}

func Queue(args models.Command) *http.Response {
	margs, err := json.Marshal(args)
	if err != nil {
		log.Fatalf("Failed to marshal command args\n %v\n %v", args, err)
	}

	response, err := client.Post("https://api.lolqueue.com/queue", "application/json", bytes.NewBuffer(margs))
	if err != nil {
		log.Fatalf("Post failed %v", err)
	}

	log.Println(response.StatusCode)

	var res map[string]interface{}

	json.NewDecoder(response.Body).Decode(&res)
	log.Println(res["message"])

	return response
}

func Setup(profile models.Profile) int {

	log.Printf("this is profile from within the setup function\n%#v", profile)
	data, err := json.Marshal(profile)
	if err != nil {
		log.Println("error marshalling profile data.")
	}

	resp, err := http.Post("https://api.lolqueue.com/user", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("fatal err 1 %v", err)
	}

	return resp.StatusCode
}
