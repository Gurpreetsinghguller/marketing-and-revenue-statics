package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	event_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/event/usecase"
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

// TrackEventRequest represents event tracking input
type TrackEventRequest struct {
	CampaignID string                 `json:"campaign_id"`
	EventType  string                 `json:"event_type"`
	UserID     string                 `json:"user_id"`
	Timestamp  string                 `json:"timestamp"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// EventHandler handles event requests
type EventHandler struct {
	usecase *event_usecase.EventUseCase
}

// NewEventHandler creates a new event handler
func NewEventHandler(uc *event_usecase.EventUseCase) *EventHandler {
	return &EventHandler{
		usecase: uc,
	}
}

// TrackEventHandler tracks user events (clicks, impressions, conversions)
func (h *EventHandler) TrackEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req TrackEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.CampaignID == "" || req.EventType == "" {
		http.Error(w, `{"error": "Campaign ID and event type are required"}`, http.StatusBadRequest)
		return
	}

	// Convert metadata
	metadata := domain.Metadata{}
	if req.Metadata != nil {
		if amount, ok := req.Metadata["amount"].(float64); ok {
			metadata.Amount = amount
		}
		if source, ok := req.Metadata["source"].(string); ok {
			metadata.Source = source
		}
		if device, ok := req.Metadata["device"].(string); ok {
			metadata.Device = device
		}
	}

	// Call usecase
	usecaseReq := event_usecase.CreateEventRequest{
		CampaignID: req.CampaignID,
		UserID:     req.UserID,
		EventType:  domain.EventType(req.EventType),
		Timestamp:  req.Timestamp,
		Metadata:   metadata,
	}

	event, err := h.usecase.TrackEvent(context.Background(), usecaseReq)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// GetEventsHandler retrieves event logs with filters
func (h *EventHandler) GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Extract query parameters
	campaignID := r.URL.Query().Get("campaign_id")
	eventType := r.URL.Query().Get("event_type")

	var events []domain.Event
	var err error

	// Apply filters
	if campaignID != "" {
		events, err = h.usecase.GetEventsByCampaign(context.Background(), campaignID)
	} else if eventType != "" {
		events, err = h.usecase.GetEventsByType(context.Background(), domain.EventType(eventType))
	} else {
		events, err = h.usecase.GetEventsByUser(context.Background(), userID)
	}

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}
