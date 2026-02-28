package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// Event represents a user interaction event
type Event struct {
	ID         string                 `json:"id"`
	CampaignID string                 `json:"campaign_id"`
	EventType  string                 `json:"event_type"` // impression, click, conversion
	UserID     string                 `json:"user_id"`
	Timestamp  time.Time              `json:"timestamp"`
	MetaData   map[string]interface{} `json:"metadata"`
}

// TrackEventHandler tracks user events (clicks, impressions, conversions)
func TrackEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	// TODO: Validate event type (impression, click, conversion)
	// TODO: Validate campaign exists
	// TODO: Apply rate limiting
	// TODO: Save to EventLog (persistence)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Event tracked successfully"})
}

// GetEventsHandler retrieves event logs with filters
func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO: Extract query params (campaign_id, event_type, date_range)
	// TODO: Check authorization
	// TODO: Fetch from EventLog

	events := []Event{}
	json.NewEncoder(w).Encode(events)
}
