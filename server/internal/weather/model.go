package weather

import "time"

type WeatherRequest struct {
	Date        string  `json:"date"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

type WeatherRecord struct {
	Timestamp   time.Time `json:"timestamp"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
}
