package handler

import "time"

// Event represents a user interaction event response.
type Event struct {
	ID         string                 `json:"id"`
	CampaignID string                 `json:"campaign_id"`
	EventType  string                 `json:"event_type"`
	UserID     string                 `json:"user_id"`
	Timestamp  time.Time              `json:"timestamp"`
	MetaData   map[string]interface{} `json:"metadata"`
}

type TrackEventResponse = Event

type EventsListResponse = []Event
