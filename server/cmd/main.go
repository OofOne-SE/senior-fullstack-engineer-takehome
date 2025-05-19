package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"server/pkg/controller"
	"server/pkg/database"
	"server/pkg/websocket"
)

func main() {
	dbErr := godotenv.Load(".env")
	if dbErr != nil {
		log.Fatal("Error loading .env file")
	}

	database.Init()

	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.LoadHTMLGlob("resources/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/ws", func(c *gin.Context) {
		websocket.HubInstance().HandleConnections(c.Writer, c.Request)
	})

	go websocket.HubInstance().Run()

	api := r.Group("/api")

	{
		v1 := api.Group("/v1")
		{
			weather := v1.Group("/weather")

			weather.POST("", controller.PostWeather)
			weather.GET("/range", controller.GetWeatherByRange)
			weather.GET("/day", controller.GetWeatherByDate)
		}
	}

	r.Run(":8080")
}
