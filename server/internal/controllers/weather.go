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
	engine.GET("/weather/day", func(c *gin.Context) {
		date := c.Query("date")
		if date == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Date is required"})
			return
		}
		ts, err := time.Parse("2006-01-02", date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		rows, err := db.DB.Query(context.Background(), `
			SELECT temperature, humidity
			FROM weather
			WHERE timestamp = $1
		`, ts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
			return
		}
		defer rows.Close()
		if rows.Next() {
			var temperature, humidity float64
			if err := rows.Scan(&temperature, &humidity); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"date":        date,
				"temperature": temperature,
				"humidity":    humidity,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "No data found for the given date"})
		}
		// Get weather for a date range
	})

	engine.GET("/weather/range", func(c *gin.Context) {
		start := c.Query("start")
		end := c.Query("end")

		if start == "" || end == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Start and end dates are required"})
			return
		}

		startDate, err := time.Parse("2006-01-02", start)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}

		endDate, err := time.Parse("2006-01-02", end)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}

		// Make endDate inclusive by adding one full day
		endDate = endDate.Add(24 * time.Hour)

		rows, err := db.DB.Query(context.Background(), `
		SELECT timestamp, temperature, humidity
		FROM weather
		WHERE timestamp >= $1 AND timestamp < $2
		ORDER BY timestamp ASC
	`, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
			return
		}
		defer rows.Close()

		var results []map[string]interface{}
		for rows.Next() {
			var ts time.Time
			var temp, hum float64
			if err := rows.Scan(&ts, &temp, &hum); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read query result"})
				return
			}
			results = append(results, gin.H{
				"timestamp":   ts,
				"temperature": temp,
				"humidity":    hum,
			})
		}

		c.JSON(http.StatusOK, gin.H{"results": results})
	})

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
