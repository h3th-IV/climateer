package model

// This model will represent countries, as they are related to climate measurements.
type Country struct {
	ID          int    `json:"id"`
	CountryCode string `json:"country_code"` // ISO code (e.g., "USA", "NGA")
	CountryName string `json:"country_name"`
}
