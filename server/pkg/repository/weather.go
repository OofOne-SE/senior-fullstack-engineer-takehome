package repository

import (
	"context"
	"server/pkg/database"
	"server/pkg/model"
	"time"
)

func InsertWeather(ts time.Time, temp, hum float64) (model.WeatherRecord, error) {
	_, err := database.DB.Exec(context.Background(), `
        INSERT INTO weather (timestamp, temperature, humidity)
        VALUES ($1, $2, $3)
    `, ts, temp, hum)
	return model.WeatherRecord{Timestamp: ts, Temperature: temp, Humidity: hum}, err
}

func GetWeatherByDate(ts time.Time) (*model.WeatherRecord, error) {
	row := database.DB.QueryRow(context.Background(), `
        SELECT temperature, humidity FROM weather WHERE timestamp = $1
    `, ts)

	var temp, hum float64
	if err := row.Scan(&temp, &hum); err != nil {
		return nil, err
	}

	return &model.WeatherRecord{Timestamp: ts, Temperature: temp, Humidity: hum}, nil
}

func GetWeatherInRange(start, end time.Time) ([]model.WeatherRecord, error) {
	rows, err := database.DB.Query(context.Background(), `
        SELECT timestamp, temperature, humidity
        FROM weather
        WHERE timestamp >= $1 AND timestamp < $2
        ORDER BY timestamp ASC
    `, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.WeatherRecord
	for rows.Next() {
		var r model.WeatherRecord
		if err := rows.Scan(&r.Timestamp, &r.Temperature, &r.Humidity); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}
