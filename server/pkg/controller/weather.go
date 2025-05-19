package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/pkg/model"
	"server/pkg/repository"
	"server/pkg/websocket"
	"time"
)

const dateFormat = "2006-01-02"

func PostWeather(c *gin.Context) {
	var req model.WeatherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	ts, err := time.Parse(dateFormat, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	data, err := repository.InsertWeather(ts, req.Temperature, req.Humidity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data"})
		return
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}

	websocket.HubInstance().SendUpdate(string(jsonBytes))
	c.JSON(http.StatusCreated, gin.H{"message": "Data inserted"})
}

func GetWeatherByDate(c *gin.Context) {
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

	record, err := repository.GetWeatherByDate(ts)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found for the given date"})
		return
	}

	c.JSON(http.StatusOK, record)
}

func GetWeatherByRange(c *gin.Context) {
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

	results, err := repository.GetWeatherInRange(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
