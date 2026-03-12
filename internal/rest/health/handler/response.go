package handler

type HealthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}
