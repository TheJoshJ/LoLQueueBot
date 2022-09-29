package handlers

import "github.com/gin-gonic/gin"

type Command struct {
	Gamemode  string `json: "gamemode"`
	Primary   string `json: "primary"`
	Secondary string `json: "secondary"`
	Fill      string `json: "fill"`
}

func Queue(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}
