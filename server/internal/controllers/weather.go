package weather

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/db"
	"time"
)

type WeatherRequest struct {
	Date        string  `json:"date"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

func Routes(engine *gin.Engine) {
	engine.POST("/weather", func(context *gin.Context) {
		var req WeatherRequest

		if err := context.ShouldBindJSON(&req); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		ts, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}

		err = InsertWeather(ts, req.Temperature, req.Humidity)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data"})
			return
		}

		context.JSON(http.StatusAccepted, gin.H{"message": "Data inserted"})
	})
}

func InsertWeather(timestamp time.Time, temperature, humidity float64) error {
	_, err := db.DB.Exec(context.Background(), `
        INSERT INTO weather (timestamp, temperature, humidity)
        VALUES ($1, $2, $3)
    `, timestamp, temperature, humidity)

	return err
}
