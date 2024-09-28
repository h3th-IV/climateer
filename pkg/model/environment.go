package model

// This represents sea level rise data for specific regions or countries.
type SeaLevelData struct {
	ID        int     `json:"id"`
	CountryID int     `json:"country_id"`
	Year      int     `json:"year"`
	SeaLevel  float64 `json:"sea_level"` // Sea level in meters or cm
	Country   Country `json:"country"`
}
