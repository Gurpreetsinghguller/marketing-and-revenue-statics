package handler

// TrackEventRequest represents event tracking input
type TrackEventRequest struct {
	CampaignID string `json:"campaign_id" validate:"required,min=1"`
	EventType  string `json:"event_type" validate:"required,oneof=clicks impressions conversions"`
	// any public/external system can send this event that's why this userID is optional
	UserID    string                 `json:"user_id" validate:"omitempty,min=1"`
	Timestamp string                 `json:"timestamp" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Metadata  map[string]interface{} `json:"metadata" validate:"omitempty"`
}
