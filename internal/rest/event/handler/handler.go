package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	event_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/event/usecase"
	"github.com/go-playground/validator/v10"
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

// EventHandler handles event requests
type EventHandler struct {
	usecase *event_usecase.EventUseCase
}

var eventRequestValidator = validator.New()

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

	req.CampaignID = strings.TrimSpace(req.CampaignID)
	req.EventType = strings.ToLower(strings.TrimSpace(req.EventType))
	req.UserID = strings.TrimSpace(req.UserID)
	req.Timestamp = strings.TrimSpace(req.Timestamp)

	if err := eventRequestValidator.Struct(req); err != nil {
		http.Error(w, `{"error": "`+formatEventValidationError(err)+`"}`, http.StatusBadRequest)
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
	eventInput := &domain.Event{
		CampaignID: req.CampaignID,
		UserID:     req.UserID,
		EventType:  domain.EventType(req.EventType),
		Timestamp:  req.Timestamp,
		Metadata:   metadata,
	}

	event, err := h.usecase.TrackEvent(context.Background(), eventInput)
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

func formatEventValidationError(err error) string {
	if err == nil {
		return "invalid request payload"
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) && len(validationErrors) > 0 {
		field := strings.ToLower(validationErrors[0].Field())
		tag := validationErrors[0].Tag()

		switch tag {
		case "required":
			return field + " is required"
		case "oneof":
			return field + " has invalid value"
		case "datetime":
			return field + " must be RFC3339 format"
		case "min":
			return field + " cannot be empty"
		default:
			return field + " is invalid"
		}
	}

	return "invalid request payload"
}
