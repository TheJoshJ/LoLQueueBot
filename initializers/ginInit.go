package initializers

import (
	"discord-test/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	r = gin.Default()
)

func CreateGinConnection() {
	handlers.CreateGinHandlers(r)

	ginErr := r.Run(":8080")

	if ginErr != nil {
		log.Printf("Error connecting to gin services %v", ginErr)
	}
}
