package event

import (
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/db"
)

// EventRepository implements domain.EventRepo using JSON file storage
type EventRepository struct {
	storage *db.StorageMgr
	events  []domain.Event
}

// NewEventRepository creates a new event repository
func NewEventRepository() *EventRepository {
	repo := &EventRepository{
		storage: db.NewStorageMgr(),
		events:  []domain.Event{},
	}
	// Load existing events from file
	repo.storage.ReadJSON(EventsFile, &repo.events)
	return repo
}

// Create saves a new event
func (r *EventRepository) Create(event *domain.Event) error {
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()

	// Generate ID if not provided
	if event.ID == "" {
		event.ID = "evt_" + generateRandomString(12)
	}

	r.events = append(r.events, *event)
	return r.storage.WriteJSON(EventsFile, r.events)
}

// GetByID retrieves an event by ID
func (r *EventRepository) GetByID(id string) (*domain.Event, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	for i := range r.events {
		if r.events[i].ID == id {
			return &r.events[i], nil
		}
	}
	return nil, ErrNotFound
}

// GetByCampaignID retrieves all events for a campaign
func (r *EventRepository) GetByCampaignID(campaignID string) ([]domain.Event, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	var result []domain.Event
	for i := range r.events {
		if r.events[i].CampaignID == campaignID {
			result = append(result, r.events[i])
		}
	}
	return result, nil
}

// GetByUserID retrieves all events for a user
func (r *EventRepository) GetByUserID(userID string) ([]domain.Event, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	var result []domain.Event
	for i := range r.events {
		if r.events[i].UserID == userID {
			result = append(result, r.events[i])
		}
	}
	return result, nil
}

// GetByEventType retrieves events by type
func (r *EventRepository) GetByEventType(eventType domain.EventType) ([]domain.Event, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	var result []domain.Event
	for i := range r.events {
		if r.events[i].EventType == eventType {
			result = append(result, r.events[i])
		}
	}
	return result, nil
}

// GetByDateRange retrieves events within a date range
func (r *EventRepository) GetByDateRange(startDate, endDate string) ([]domain.Event, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	var result []domain.Event
	for i := range r.events {
		eventTime, _ := time.Parse(time.RFC3339, r.events[i].Timestamp)
		if eventTime.After(start) && eventTime.Before(end.AddDate(0, 0, 1)) {
			result = append(result, r.events[i])
		}
	}
	return result, nil
}

// GetByCampaignAndDateRange retrieves events for a campaign in date range
func (r *EventRepository) GetByCampaignAndDateRange(campaignID, startDate, endDate string) ([]domain.Event, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	var result []domain.Event
	for i := range r.events {
		if r.events[i].CampaignID == campaignID {
			eventTime, _ := time.Parse(time.RFC3339, r.events[i].Timestamp)
			if eventTime.After(start) && eventTime.Before(end.AddDate(0, 0, 1)) {
				result = append(result, r.events[i])
			}
		}
	}
	return result, nil
}

// GetByFilter retrieves events with multiple filters
func (r *EventRepository) GetByFilter(filters map[string]interface{}) ([]domain.Event, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	var result []domain.Event
	for i := range r.events {
		event := r.events[i]
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
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	return r.events, nil
}

// Delete removes an event
func (r *EventRepository) Delete(id string) error {
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()

	for i := range r.events {
		if r.events[i].ID == id {
			r.events = append(r.events[:i], r.events[i+1:]...)
			return r.storage.WriteJSON(EventsFile, r.events)
		}
	}
	return ErrNotFound
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
