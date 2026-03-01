package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
)

// EventUseCase handles event tracking business logic
type EventUseCase struct {
	eventRepo domain.EventRepo
}

// NewEventUseCase creates a new event usecase
func NewEventUseCase(eventRepo domain.EventRepo) *EventUseCase {
	return &EventUseCase{
		eventRepo: eventRepo,
	}
}

func (e *EventUseCase) TrackEvent(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	if event == nil {
		return nil, errors.New("event is required")
	}

	if event.CampaignID == "" || event.EventType == "" {
		return nil, errors.New("campaign_id and event_type are required")
	}

	newEvent := &domain.Event{
		CampaignID: event.CampaignID,
		UserID:     event.UserID,
		EventType:  event.EventType,
		Timestamp:  event.Timestamp,
		Metadata:   event.Metadata,
	}

	if err := e.eventRepo.Create(newEvent); err != nil {
		return nil, fmt.Errorf("failed to track event: %w", err)
	}

	return newEvent, nil
}

// GetEventByID retrieves an event by ID
func (e *EventUseCase) GetEventByID(ctx context.Context, eventID string) (*domain.Event, error) {
	_ = ctx

	event, err := e.eventRepo.GetByID(eventID)
	if err != nil || event == nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}

	return event, nil
}

// GetEventsByCampaign retrieves all events for a campaign
func (e *EventUseCase) GetEventsByCampaign(ctx context.Context, campaignID string) ([]domain.Event, error) {
	_ = ctx

	events, err := e.eventRepo.GetByCampaignID(campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campaign events: %w", err)
	}

	return events, nil
}

// GetEventsByUser retrieves all events for a user
func (e *EventUseCase) GetEventsByUser(ctx context.Context, userID string) ([]domain.Event, error) {
	_ = ctx

	events, err := e.eventRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user events: %w", err)
	}

	return events, nil
}

// GetEventsByType retrieves events of a specific type
func (e *EventUseCase) GetEventsByType(ctx context.Context, eventType domain.EventType) ([]domain.Event, error) {
	_ = ctx

	events, err := e.eventRepo.GetByEventType(eventType)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
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
