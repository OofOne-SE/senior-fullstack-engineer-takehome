package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"server/pkg/controller"
	"server/pkg/database"
)

// @title           Swagger Example API
func main() {
	dbErr := godotenv.Load(".env")
	if dbErr != nil {
		log.Fatal("Error loading .env file")
	}

	database.Init()

	r := gin.Default()
	r.SetTrustedProxies(nil)

	v1 := r.Group("/api/v1")
	{
		weather := v1.Group("/weather")
		{
			weather.POST("/", controller.PostWeather)
			weather.GET("/range", controller.GetWeatherByRange)
			weather.GET("/day", controller.GetWeatherByDate)
		}
	}

	r.Run(":8080")
}
