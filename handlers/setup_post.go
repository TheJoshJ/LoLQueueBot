package handlers

import (
	"discord-test/models"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func Setup(profile models.Profile) {

	data := url.Values{}
	data.Add("server", profile.Server)
	data.Add("username", profile.Username)
	data.Add("discordid", profile.DiscordID)

	resp, err := http.PostForm("http://api.lolqueue.com/user", data)
	if err != nil {
		log.Printf("fatal err 1 %v", err)
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("fatal err 3 %v", err)
	}

	log.Println(string(response))

}
