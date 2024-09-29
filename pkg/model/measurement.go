package model

import "time"

//This struct will represent a single climate measurement, linking countries, indicators, and values recorded at specific times (years).

type Measurement struct {
	ID          int       `json:"id"`
	CountryID   int       `json:"country_id"`
	IndicatorID int       `json:"indicator_id"`
	Year        int       `json:"year"`
	Value       float64   `json:"value"`
	Country     Country   `json:"country"`   // Nested struct for country details
	Indicator   Indicator `json:"indicator"` // Nested struct for indicator details
}

type UserMeasurement struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	CountryName   string    `json:"country_name"`
	IndicatorCode string    `json:"indicator_code"`
	Year          int       `json:"year"`
	Value         float64   `json:"value"`
	CreatedAt     time.Time `json:"created_at"`
}
