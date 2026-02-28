package handler

type UserInteractionCampaign struct {
	CampaignID   string  `json:"campaign_id"`
	CampaignName string  `json:"campaign_name"`
	Duration     int64   `json:"duration"`      // ← Sum of all metadata.duration
	Interactions int64   `json:"interactions"`  // ← Count of events
	FunnelStages []Stage `json:"funnel_stages"` // ← Stages visited
}

type Stage struct {
	Page  string `json:"page"`
	Count int64  `json:"drop_off_count"` // Users who stopped here
}
