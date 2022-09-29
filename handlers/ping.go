package handlers

import (
	"log"
	"net/http"
)

func Ping() {
	client := http.DefaultClient

	response, err := client.Post("https://api.lolqueue.com/ping", "application/json", nil)
	if err != nil {
		log.Fatalf("Post failed %v", err)
	}
	log.Println(response)
}
