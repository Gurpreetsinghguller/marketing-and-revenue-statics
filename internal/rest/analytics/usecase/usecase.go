package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/db"
)

// AnalyticsUseCase handles analytics and reporting business logic
type AnalyticsUseCase struct {
	db db.PersistenceDB
}

// NewAnalyticsUseCase creates a new analytics usecase
func NewAnalyticsUseCase(db db.PersistenceDB) *AnalyticsUseCase {
	return &AnalyticsUseCase{
		db: db,
	}
}

// CampaignStats represents campaign statistics
type CampaignStats struct {
	CampaignID       string
	CampaignName     string
	TotalImpressions int64
	TotalClicks      int64
	TotalConversions int64
	TotalRevenue     float64
	CTR              float64 // Click-through rate
	ConversionRate   float64
}

// GetPublicStats returns public statistics
func (a *AnalyticsUseCase) GetPublicStats(ctx context.Context) (map[string]interface{}, error) {
	// Fetch all public campaigns
	campaigns, err := a.db.List(ctx, "campaign:")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campaigns: %w", err)
	}

	var publicCampaigns []interface{}
	for _, campaign := range campaigns {
		if c, ok := campaign.(*domain.Campaign); ok && c.IsPublic {
			publicCampaigns = append(publicCampaigns, c)
		}
	}

	stats := map[string]interface{}{
		"total_campaigns": len(publicCampaigns),
		"campaigns":       publicCampaigns,
	}

	return stats, nil
}

// GetCampaignStats returns statistics for a specific campaign
func (a *AnalyticsUseCase) GetCampaignStats(ctx context.Context, campaignID string) (*CampaignStats, error) {
	// Fetch campaign
	campaignKey := fmt.Sprintf("campaign:%s", campaignID)
	campaignInterface, err := a.db.Read(ctx, campaignKey)
	if err != nil || campaignInterface == nil {
		return nil, errors.New("campaign not found")
	}

	campaign, ok := campaignInterface.(*domain.Campaign)
	if !ok {
		return nil, errors.New("invalid campaign data format")
	}

	// Fetch events for campaign
	events, err := a.db.List(ctx, fmt.Sprintf("campaign:%s:events:", campaignID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campaign events: %w", err)
	}

	stats := &CampaignStats{
		CampaignID:   campaignID,
		CampaignName: campaign.Name,
	}

	for _, result := range events {
		if eventID, ok := result.(string); ok {
			eventKey := fmt.Sprintf("event:%s", eventID)
			eventInterface, err := a.db.Read(ctx, eventKey)
			if err != nil || eventInterface == nil {
				continue
			}

			event, ok := eventInterface.(*domain.Event)
			if !ok {
				continue
			}

			switch event.EventType {
			case domain.EventType("impressions"):
				stats.TotalImpressions++
			case domain.EventType("clicks"):
				stats.TotalClicks++
			case domain.EventType("conversions"):
				stats.TotalConversions++
				stats.TotalRevenue += event.Metadata.Amount
			}
		}
	}

	// Calculate rates
	if stats.TotalImpressions > 0 {
		stats.CTR = float64(stats.TotalClicks) / float64(stats.TotalImpressions)
	}
	if stats.TotalClicks > 0 {
		stats.ConversionRate = float64(stats.TotalConversions) / float64(stats.TotalClicks)
	}

	return stats, nil
}

// GetAnalyticsReport returns general analytics report
func (a *AnalyticsUseCase) GetAnalyticsReport(ctx context.Context) (map[string]interface{}, error) {
	campaigns, err := a.db.List(ctx, "campaign:")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campaigns: %w", err)
	}

	report := map[string]interface{}{
		"total_campaigns": len(campaigns),
		"campaigns":       []interface{}{},
	}

	campaignsList := report["campaigns"].([]interface{})
	for _, c := range campaigns {
		if campaign, ok := c.(*domain.Campaign); ok {
			campaignsList = append(campaignsList, campaign)
		}
	}
	report["campaigns"] = campaignsList

	return report, nil
}

// GetDailyReport returns daily analytics report (simplified)
func (a *AnalyticsUseCase) GetDailyReport(ctx context.Context, date string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"date":          date,
		"total_events":  0,
		"total_revenue": 0.0,
		"top_campaigns": []interface{}{},
	}, nil
}

// GetWeeklyReport returns weekly analytics report (simplified)
func (a *AnalyticsUseCase) GetWeeklyReport(ctx context.Context, weekStart string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"week_start":    weekStart,
		"total_events":  0,
		"total_revenue": 0.0,
		"top_campaigns": []interface{}{},
	}, nil
}

// GetMonthlyReport returns monthly analytics report (simplified)
func (a *AnalyticsUseCase) GetMonthlyReport(ctx context.Context, month string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"month":         month,
		"total_events":  0,
		"total_revenue": 0.0,
		"top_campaigns": []interface{}{},
	}, nil
}
