package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/db"
)

const eventPrefix = "events"

type eventDBMetadata struct {
	Amount float64 `json:"amount"`
	Source string  `json:"source"`
	Device string  `json:"device"`
}

type eventDBModel struct {
	ID         string          `json:"id"`
	CampaignID string          `json:"campaign_id"`
	UserID     string          `json:"user_id"`
	EventType  string          `json:"event_type"`
	Timestamp  string          `json:"timestamp"`
	Metadata   eventDBMetadata `json:"metadata"`
}

// EventRepository implements domain.EventRepo using JSON file storage
type EventRepository struct {
	storage db.PersistenceDB
}

// NewEventRepository creates a new event repository
func NewEventRepository(storage ...db.PersistenceDB) *EventRepository {
	selected := db.PersistenceDB(db.NewStorageMgr())
	if len(storage) > 0 && storage[0] != nil {
		selected = storage[0]
	}

	return &EventRepository{storage: selected}
}

// Create saves a new event
func (r *EventRepository) Create(event *domain.Event) error {
	// Generate ID if not provided
	if event.ID == "" {
		event.ID = db.GenerateID("evt")
	}

	return r.storage.Create(context.Background(), eventKey(event.ID), toEventDBModel(*event))
}

// GetByID retrieves an event by ID
func (r *EventRepository) GetByID(id string) (*domain.Event, error) {
	stored, err := r.storage.Read(context.Background(), eventKey(id))
	if err != nil {
		return nil, err
	}

	model, err := decodeEventDBModel(stored)
	if err != nil {
		return nil, err
	}

	entity := model.toDomain()
	return &entity, nil
}

func (r *EventRepository) getAll() ([]domain.Event, error) {
	stored, err := r.storage.List(context.Background(), eventPrefix)
	if err != nil {
		return nil, err
	}

	events := make([]domain.Event, 0, len(stored))
	for _, item := range stored {
		model, err := decodeEventDBModel(item)
		if err != nil {
			return nil, err
		}
		events = append(events, model.toDomain())
	}

	return events, nil
}

// GetByCampaignID retrieves all events for a campaign
func (r *EventRepository) GetByCampaignID(campaignID string) ([]domain.Event, error) {
	events, err := r.getAll()
	if err != nil {
		return nil, err
	}

	var result []domain.Event
	for i := range events {
		if events[i].CampaignID == campaignID {
			result = append(result, events[i])
		}
	}
	return result, nil
}

// GetByUserID retrieves all events for a user
func (r *EventRepository) GetByUserID(userID string) ([]domain.Event, error) {
	events, err := r.getAll()
	if err != nil {
		return nil, err
	}

	var result []domain.Event
	for i := range events {
		if events[i].UserID == userID {
			result = append(result, events[i])
		}
	}
	return result, nil
}

// GetByEventType retrieves events by type
func (r *EventRepository) GetByEventType(eventType domain.EventType) ([]domain.Event, error) {
	events, err := r.getAll()
	if err != nil {
		return nil, err
	}

	var result []domain.Event
	for i := range events {
		if events[i].EventType == eventType {
			result = append(result, events[i])
		}
	}
	return result, nil
}

// GetByDateRange retrieves events within a date range
func (r *EventRepository) GetByDateRange(startDate, endDate string) ([]domain.Event, error) {
	events, err := r.getAll()
	if err != nil {
		return nil, err
	}

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	var result []domain.Event
	for i := range events {
		eventTime, _ := time.Parse(time.RFC3339, events[i].Timestamp)
		if eventTime.After(start) && eventTime.Before(end.AddDate(0, 0, 1)) {
			result = append(result, events[i])
		}
	}
	return result, nil
}

// GetByCampaignAndDateRange retrieves events for a campaign in date range
func (r *EventRepository) GetByCampaignAndDateRange(campaignID, startDate, endDate string) ([]domain.Event, error) {
	events, err := r.getAll()
	if err != nil {
		return nil, err
	}

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	var result []domain.Event
	for i := range events {
		if events[i].CampaignID == campaignID {
			eventTime, _ := time.Parse(time.RFC3339, events[i].Timestamp)
			if eventTime.After(start) && eventTime.Before(end.AddDate(0, 0, 1)) {
				result = append(result, events[i])
			}
		}
	}
	return result, nil
}

// GetByFilter retrieves events with multiple filters
func (r *EventRepository) GetByFilter(filters map[string]interface{}) ([]domain.Event, error) {
	events, err := r.getAll()
	if err != nil {
		return nil, err
	}

	var result []domain.Event
	for i := range events {
		event := events[i]
		match := true

		if campaignID, ok := filters["campaign_id"].(string); ok && event.CampaignID != campaignID {
			match = false
		}
		if eventType, ok := filters["event_type"].(string); ok && string(event.EventType) != eventType {
			match = false
		}
		if userID, ok := filters["user_id"].(string); ok && event.UserID != userID {
			match = false
		}
		if source, ok := filters["source"].(string); ok && event.Metadata.Source != source {
			match = false
		}

		if match {
			result = append(result, event)
		}
	}
	return result, nil
}

// GetAll retrieves all events
func (r *EventRepository) GetAll() ([]domain.Event, error) {
	return r.getAll()
}

// Delete removes an event
func (r *EventRepository) Delete(id string) error {
	return r.storage.Delete(context.Background(), eventKey(id))
}

// GetAggregatedMetrics returns aggregated metrics for analytics
func (r *EventRepository) GetAggregatedMetrics(campaignID, startDate, endDate string) (map[string]interface{}, error) {
	events, _ := r.GetByCampaignAndDateRange(campaignID, startDate, endDate)

	impressions := int64(0)
	clicks := int64(0)
	conversions := int64(0)
	totalRevenue := float64(0)

	for _, event := range events {
		switch event.EventType {
		case "impressions":
			impressions++
		case "clicks":
			clicks++
		case "conversions":
			conversions++
			totalRevenue += event.Metadata.Amount
		}
	}

	return map[string]interface{}{
		"impressions":   impressions,
		"clicks":        clicks,
		"conversions":   conversions,
		"total_revenue": totalRevenue,
	}, nil
}

// GetEventCountByType returns event counts grouped by type
func (r *EventRepository) GetEventCountByType(campaignID, startDate, endDate string) (map[domain.EventType]int64, error) {
	events, _ := r.GetByCampaignAndDateRange(campaignID, startDate, endDate)

	counts := make(map[domain.EventType]int64)
	for _, event := range events {
		counts[event.EventType]++
	}
	return counts, nil
}

// GetTotalRevenue sums revenue from conversions
func (r *EventRepository) GetTotalRevenue(campaignID, startDate, endDate string) (float64, error) {
	events, _ := r.GetByCampaignAndDateRange(campaignID, startDate, endDate)

	totalRevenue := float64(0)
	for _, event := range events {
		if event.EventType == "conversions" {
			totalRevenue += event.Metadata.Amount
		}
	}
	return totalRevenue, nil
}

// GetEventsByChannel retrieves events grouped by source/channel
func (r *EventRepository) GetEventsByChannel(campaignID, startDate, endDate string) (map[string][]domain.Event, error) {
	events, _ := r.GetByCampaignAndDateRange(campaignID, startDate, endDate)

	channels := make(map[string][]domain.Event)
	for _, event := range events {
		channel := event.Metadata.Source
		channels[channel] = append(channels[channel], event)
	}
	return channels, nil
}

func toEventDBModel(e domain.Event) eventDBModel {
	return eventDBModel{
		ID:         e.ID,
		CampaignID: e.CampaignID,
		UserID:     e.UserID,
		EventType:  string(e.EventType),
		Timestamp:  e.Timestamp,
		Metadata: eventDBMetadata{
			Amount: e.Metadata.Amount,
			Source: e.Metadata.Source,
			Device: e.Metadata.Device,
		},
	}
}

func (m eventDBModel) toDomain() domain.Event {
	return domain.Event{
		ID:         m.ID,
		CampaignID: m.CampaignID,
		UserID:     m.UserID,
		EventType:  domain.EventType(m.EventType),
		Timestamp:  m.Timestamp,
		Metadata: domain.Metadata{
			Amount: m.Metadata.Amount,
			Source: m.Metadata.Source,
			Device: m.Metadata.Device,
		},
	}
}

func decodeEventDBModel(value interface{}) (eventDBModel, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return eventDBModel{}, err
	}
	var model eventDBModel
	if err := json.Unmarshal(b, &model); err != nil {
		return eventDBModel{}, err
	}
	return model, nil
}

func eventKey(id string) string {
	return eventPrefix + "/" + id
}
