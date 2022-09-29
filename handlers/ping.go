package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func Ping() {
	var output map[string]string

	client := http.DefaultClient

	r, err := client.Get("https://api.lolqueue.com/ping")
	if err != nil {
		log.Fatalf("Post failed %v", err)
	}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&output)
	log.Println(output["message"])
}
