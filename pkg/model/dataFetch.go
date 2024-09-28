package model

// Represents a request to fetch climate data (used for querying dynamic data from CDSAPI).
type DataFetchRequest struct {
	Variable string `json:"variable"` // e.g., "2m_temperature"
	Year     int    `json:"year"`     // Year for the data (e.g., 2021)
	Month    string `json:"month"`    // Month (e.g., "01" for January)
	Day      string `json:"day"`      // Day (e.g., "01" for 1st of the month)
	Time     string `json:"time"`     // Time (e.g., "12:00")
}

type DataFetchResponse struct {
	Message  string `json:"message"`
	FileName string `json:"file_name"` // Name of the file where data is stored
	FileURL  string `json:"file_url"`  // URL to download the file
}
