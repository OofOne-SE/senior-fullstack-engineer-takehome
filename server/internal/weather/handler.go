package weather

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const dateFormat = "2006-01-02"

func Routes(router *gin.Engine) {
	router.POST("/weather", handlePostWeather)
	router.GET("/weather/day", handleGetWeatherByDate)
	router.GET("/weather/range", handleGetWeatherRange)
}

func handlePostWeather(c *gin.Context) {
	var req WeatherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	ts, err := time.Parse(dateFormat, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	if err := InsertWeather(ts, req.Temperature, req.Humidity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Data inserted"})
}

func handleGetWeatherByDate(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date is required"})
		return
	}

	ts, err := time.Parse(dateFormat, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	record, err := GetWeatherByDate(ts)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found for the given date"})
		return
	}

	c.JSON(http.StatusOK, record)
}

func handleGetWeatherRange(c *gin.Context) {
	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start and end dates are required"})
		return
	}

	start, err := time.Parse(dateFormat, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	end, err := time.Parse(dateFormat, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	end = end.Add(24 * time.Hour) // make inclusive

	results, err := GetWeatherInRange(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
