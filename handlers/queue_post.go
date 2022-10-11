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
