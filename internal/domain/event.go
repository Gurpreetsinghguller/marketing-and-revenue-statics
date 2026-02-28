package domain

type EventType string

const (
	clicks      EventType = "clicks"
	impressions EventType = "impressions"
	conversions EventType = "conversions"
)

type Event struct {
	ID         string    `json:"id"`
	CampaignID string    `json:"campaign_id"`
	UserID     string    `json:"user_id"`
	EventType  EventType `json:"event_type"`
	Timestamp  string    `json:"timestamp"`
	Metadata   Metadata  `json:"metadata"`
}
type Metadata struct {
	Amount float64 `json:"amount"` // Revenue amount (for conversions)
	Source string  `json:"source"` // e.g., "facebook"
	Device string  `json:"device"` // e.g., "mobile"
}

type EventRepo interface {
	// Create saves a new event to the EventLog
	Create(event *Event) error

	// GetByID retrieves a single event by ID
	GetByID(id string) (*Event, error)

	// GetByCampaignID retrieves all events for a specific campaign
	GetByCampaignID(campaignID string) ([]Event, error)

	// GetByUserID retrieves all events for a specific user
	GetByUserID(userID string) ([]Event, error)

	// GetByEventType retrieves all events of a specific type (impression, click, conversion)
	GetByEventType(eventType EventType) ([]Event, error)

	// GetByDateRange retrieves events within a date range
	GetByDateRange(startDate, endDate string) ([]Event, error)

	// GetByCampaignAndDateRange retrieves events for a campaign within a date range
	GetByCampaignAndDateRange(campaignID, startDate, endDate string) ([]Event, error)

	// GetByFilter retrieves events with multiple filters
	// Filters: campaign_id, event_type, date_range, channel (source), user_id
	GetByFilter(filters map[string]interface{}) ([]Event, error)

	// GetAll retrieves all events from EventLog
	GetAll() ([]Event, error)

	// Delete removes an event from EventLog
	Delete(id string) error

	// GetAggregatedMetrics returns event counts and sums for analytics
	// Returns: total impressions, clicks, conversions, total revenue
	GetAggregatedMetrics(campaignID, startDate, endDate string) (map[string]interface{}, error)

	// GetEventCountByType counts events grouped by event_type for a campaign and date range
	GetEventCountByType(campaignID, startDate, endDate string) (map[EventType]int64, error)

	// GetTotalRevenue sums revenue from conversion events for a campaign and date range
	GetTotalRevenue(campaignID, startDate, endDate string) (float64, error)

	// GetEventsByChannel retrieves events grouped by source/channel
	GetEventsByChannel(campaignID, startDate, endDate string) (map[string][]Event, error)
}
