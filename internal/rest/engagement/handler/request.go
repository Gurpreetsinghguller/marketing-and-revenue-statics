package handler

type UserEngagementRequest struct {
	UserID string `json:"user_id"`
}

type CampaignFunnelRequest struct {
	CampaignID string `json:"campaign_id"`
}

type UserCampaignEngagementRequest struct {
	UserID     string `json:"user_id"`
	CampaignID string `json:"campaign_id"`
}
