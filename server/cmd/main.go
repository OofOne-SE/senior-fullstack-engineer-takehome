package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"server/internal/controllers"
	"server/internal/db"
)

func main() {
	dbErr := godotenv.Load(".env")
	if dbErr != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()

	engine := gin.Default()
	engine.SetTrustedProxies(nil)

	// Set up routes from the weather controller
	weather.Routes(engine)

	engine.Run(":8080")
}
