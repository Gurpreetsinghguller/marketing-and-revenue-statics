package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/db"
)

// EventUseCase handles event tracking business logic
type EventUseCase struct {
	db db.PersistenceDB
}

// NewEventUseCase creates a new event usecase
func NewEventUseCase(db db.PersistenceDB) *EventUseCase {
	return &EventUseCase{
		db: db,
	}
}

// CreateEventRequest represents event creation input
type CreateEventRequest struct {
	CampaignID string
	UserID     string
	EventType  domain.EventType
	Timestamp  string
	Metadata   domain.Metadata
}

// TrackEvent creates a new event
func (e *EventUseCase) TrackEvent(ctx context.Context, req CreateEventRequest) (*domain.Event, error) {
	if req.CampaignID == "" || req.EventType == "" {
		return nil, errors.New("campaign_id and event_type are required")
	}

	event := &domain.Event{
		ID:         fmt.Sprintf("event_%d", len([]int{})), // Use UUID in production
		CampaignID: req.CampaignID,
		UserID:     req.UserID,
		EventType:  req.EventType,
		Timestamp:  req.Timestamp,
		Metadata:   req.Metadata,
	}

	// Save event
	key := fmt.Sprintf("event:%s", event.ID)
	if err := e.db.Create(ctx, key, event); err != nil {
		return nil, fmt.Errorf("failed to track event: %w", err)
	}

	// Index by campaign for quick lookup
	campaignEventKey := fmt.Sprintf("campaign:%s:events:%s", req.CampaignID, event.ID)
	_ = e.db.Create(ctx, campaignEventKey, event.ID)

	// Index by user if provided
	if req.UserID != "" {
		userEventKey := fmt.Sprintf("user:%s:events:%s", req.UserID, event.ID)
		_ = e.db.Create(ctx, userEventKey, event.ID)
	}

	return event, nil
}

// GetEventByID retrieves an event by ID
func (e *EventUseCase) GetEventByID(ctx context.Context, eventID string) (*domain.Event, error) {
	key := fmt.Sprintf("event:%s", eventID)
	eventInterface, err := e.db.Read(ctx, key)
	if err != nil || eventInterface == nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}

	event, ok := eventInterface.(*domain.Event)
	if !ok {
		return nil, errors.New("invalid event data format")
	}

	return event, nil
}

// GetEventsByCampaign retrieves all events for a campaign
func (e *EventUseCase) GetEventsByCampaign(ctx context.Context, campaignID string) ([]domain.Event, error) {
	events := []domain.Event{}
	results, err := e.db.List(ctx, fmt.Sprintf("campaign:%s:events:", campaignID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campaign events: %w", err)
	}

	for _, result := range results {
		if eventID, ok := result.(string); ok {
			event, err := e.GetEventByID(ctx, eventID)
			if err == nil {
				events = append(events, *event)
			}
		}
	}

	return events, nil
}

// GetEventsByUser retrieves all events for a user
func (e *EventUseCase) GetEventsByUser(ctx context.Context, userID string) ([]domain.Event, error) {
	events := []domain.Event{}
	results, err := e.db.List(ctx, fmt.Sprintf("user:%s:events:", userID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user events: %w", err)
	}

	for _, result := range results {
		if eventID, ok := result.(string); ok {
			event, err := e.GetEventByID(ctx, eventID)
			if err == nil {
				events = append(events, *event)
			}
		}
	}

	return events, nil
}

// GetEventsByType retrieves events of a specific type
func (e *EventUseCase) GetEventsByType(ctx context.Context, eventType domain.EventType) ([]domain.Event, error) {
	events := []domain.Event{}
	results, err := e.db.List(ctx, "event:")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}

	for _, result := range results {
		if event, ok := result.(*domain.Event); ok && event.EventType == eventType {
			events = append(events, *event)
		}
	}

	return events, nil
}

// GetAggregatedMetrics returns aggregated event metrics
func (e *EventUseCase) GetAggregatedMetrics(ctx context.Context, campaignID string) (map[string]interface{}, error) {
	events, err := e.GetEventsByCampaign(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	metrics := map[string]interface{}{
		"total_impressions": int64(0),
		"total_clicks":      int64(0),
		"total_conversions": int64(0),
		"total_revenue":     float64(0),
	}

	for _, event := range events {
		switch event.EventType {
		case domain.EventType("impressions"):
			metrics["total_impressions"] = metrics["total_impressions"].(int64) + 1
		case domain.EventType("clicks"):
			metrics["total_clicks"] = metrics["total_clicks"].(int64) + 1
		case domain.EventType("conversions"):
			metrics["total_conversions"] = metrics["total_conversions"].(int64) + 1
			metrics["total_revenue"] = metrics["total_revenue"].(float64) + event.Metadata.Amount
		}
	}

	return metrics, nil
}
