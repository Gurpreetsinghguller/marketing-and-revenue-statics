package event

import "context"

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

	// Broker information
	Source    string `json:"source"`     // "mqtt", "kafka", "pubsub", "http", etc.
	MessageID string `json:"message_id"` // Broker's message ID for deduplication

	// Processing metadata (for event-time analytics)
	ProcessedAt string `json:"processed_at"` // When the event was processed
	Sequence    int64  `json:"sequence"`     // Order number for ordering guarantees
}

type Metadata struct {
	Amount float64 `json:"amount"` // Revenue amount (for conversions)
	Source string  `json:"source"` // e.g., "facebook"
	Device string  `json:"device"` // e.g., "mobile"
}

// EventHandler is the function signature for processing received events
type EventHandler func(ctx context.Context, event *Event) error

// EventBroker defines the interface for any event broker/message queue
type EventBroker interface {
	// Start begins consuming events from the broker
	// The handler function receives parsed events and should handle them
	Start(ctx context.Context, handler EventHandler) error

	// Close gracefully closes the broker connection
	Close() error

	// Publish sends an event to the broker (optional, for some implementations)
	Publish(ctx context.Context, event *Event) error

	// Health checks broker connectivity
	Health(ctx context.Context) error
}
