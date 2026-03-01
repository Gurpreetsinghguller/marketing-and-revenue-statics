package usecase

import (
	"context"
	"fmt"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
)

// EngagementUseCase handles user engagement and behavioral data
type EngagementUseCase struct {
	eventRepo domain.EventRepo
}

// NewEngagementUseCase creates a new engagement usecase
func NewEngagementUseCase(eventRepo domain.EventRepo) *EngagementUseCase {
	return &EngagementUseCase{
		eventRepo: eventRepo,
	}
}

// UserEngagement represents user engagement metrics
type UserEngagement struct {
	UserID            string
	TotalInteractions int64
	TotalDuration     int64
	CampaignsEngaged  int64
	AverageEngagement float64
	TopCampaigns      []interface{}
}

// CampaignFunnel represents campaign funnel data
type CampaignFunnel struct {
	CampaignID     string
	Impressions    int64
	Clicks         int64
	Conversions    int64
	ConversionRate float64
	DropoffRate    float64
}

// GetUserEngagement retrieves user engagement metrics
func (e *EngagementUseCase) GetUserEngagement(ctx context.Context, userID string) (*UserEngagement, error) {
	_ = ctx

	// Fetch all events for the user
	userEvents, err := e.eventRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user events: %w", err)
	}

	engagement := &UserEngagement{
		UserID:       userID,
		TopCampaigns: []interface{}{},
	}

	campaignMap := make(map[string]int64)

	for _, event := range userEvents {
		campaignMap[event.CampaignID]++
		engagement.TotalInteractions++
	}

	engagement.CampaignsEngaged = int64(len(campaignMap))
	if engagement.TotalInteractions > 0 && engagement.CampaignsEngaged > 0 {
		engagement.AverageEngagement = float64(engagement.TotalInteractions) / float64(engagement.CampaignsEngaged)
	}

	return engagement, nil
}

// GetUserCampaignEngagement retrieves engagement for a specific user-campaign combination
func (e *EngagementUseCase) GetUserCampaignEngagement(ctx context.Context, userID, campaignID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"user_id":      userID,
		"campaign_id":  campaignID,
		"interactions": 0,
		"duration":     0,
		"event_types":  []interface{}{},
	}, nil
}

// GetCampaignFunnel retrieves campaign funnel data
func (e *EngagementUseCase) GetCampaignFunnel(ctx context.Context, campaignID string) (*CampaignFunnel, error) {
	_ = ctx

	funnel := &CampaignFunnel{
		CampaignID: campaignID,
	}

	// Fetch campaign events
	campaignEvents, err := e.eventRepo.GetByCampaignID(campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campaign events: %w", err)
	}

	for _, event := range campaignEvents {
		switch event.EventType {
		case domain.EventType("impressions"):
			funnel.Impressions++
		case domain.EventType("clicks"):
			funnel.Clicks++
		case domain.EventType("conversions"):
			funnel.Conversions++
		}
	}

	// Calculate rates
	if funnel.Impressions > 0 {
		funnel.DropoffRate = float64(funnel.Impressions-funnel.Clicks) / float64(funnel.Impressions)
		funnel.ConversionRate = float64(funnel.Conversions) / float64(funnel.Impressions)
	}

	return funnel, nil
}
