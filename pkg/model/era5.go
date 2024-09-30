package model

type ReanalysisRequest struct {
	Variable string `json:"variable"`
	Year     string `json:"year"`
	Month    string `json:"month"`
	Day      string `json:"day"`
	Time     string `json:"time"`
	Area     string `json:"area"`
}

type ReanalysisResponse struct {
	Message  string `json:"message"`
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	Error    string `json:"error,omitempty"`
}
