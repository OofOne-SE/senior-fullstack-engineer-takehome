package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"server/internal/database"
	"server/internal/weather"
)

func main() {
	dbErr := godotenv.Load(".env")
	if dbErr != nil {
		log.Fatal("Error loading .env file")
	}

	database.Init()

	engine := gin.Default()
	engine.SetTrustedProxies(nil)

	// Set up routes from the weather controller
	weather.Routes(engine)

	engine.Run(":8080")
}
