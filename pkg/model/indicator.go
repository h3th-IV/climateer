package model

// Indicators represent the different types of climate data we gather, such as temperature, precipitation, greenhouse gases, etc.
type Indicator struct {
	ID        int    `json:"id"`
	Indicator string `json:"indicator"` // Indicator name (e.g., "CO2 emissions")
	Unit      string `json:"unit"`      // Unit of measurement (e.g., "metric tons per capita")
}
