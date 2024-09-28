package model

//This model represents the various economic or industrial sectors (for emissions data).

type Sector struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`        // Sector name (e.g., "Industry", "Transport")
	Description string `json:"description"` // Description of the sector
}

// This represents different sources of energy (e.g., renewable, fossil fuels) tied to energy consumption data.
type EnergySource struct {
	ID        int    `json:"id"`
	Source    string `json:"source"`    // Energy source (e.g., "Solar", "Wind", "Coal")
	Renewable bool   `json:"renewable"` // Whether the source is renewable (true/false)
}
