package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthHandler responds with service health status.
type HealthHandler struct{}

// NewHealthHandler creates a new health handler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// GetHealthHandler returns API health status.
func (h *HealthHandler) GetHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := map[string]interface{}{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(resp)
}
