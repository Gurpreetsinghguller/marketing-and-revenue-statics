package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
)

// AnalyticsUseCase handles analytics and reporting business logic
type AnalyticsUseCase struct {
	campaignRepo domain.CampaignRepo
	eventRepo    domain.EventRepo
}

// NewAnalyticsUseCase creates a new analytics usecase
func NewAnalyticsUseCase(campaignRepo domain.CampaignRepo, eventRepo domain.EventRepo) *AnalyticsUseCase {
	return &AnalyticsUseCase{
		campaignRepo: campaignRepo,
		eventRepo:    eventRepo,
	}
}

// CampaignStats represents campaign statistics
type CampaignStats struct {
	CampaignID       string  `json:"campaign_id"`
	CampaignName     string  `json:"campaign_name"`
	TotalImpressions int64   `json:"total_impressions"`
	TotalClicks      int64   `json:"total_clicks"`
	TotalConversions int64   `json:"total_conversions"`
	TotalRevenue     float64 `json:"total_revenue"`
	CTR              float64 `json:"ctr"`
	ConversionRate   float64 `json:"conversion_rate"`
}

// GetPublicStats returns public statistics
func (a *AnalyticsUseCase) GetPublicStats(ctx context.Context) (map[string]interface{}, error) {
	_ = ctx

	// Fetch all public campaigns
	campaigns, err := a.campaignRepo.GetPublic()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campaigns: %w", err)
	}

	var publicCampaigns []interface{}
	for _, campaign := range campaigns {
		publicCampaigns = append(publicCampaigns, campaign)
	}

	stats := map[string]interface{}{
		"total_campaigns": len(publicCampaigns),
		"campaigns":       publicCampaigns,
	}

	return stats, nil
}

// GetCampaignStats returns statistics for a specific campaign
func (a *AnalyticsUseCase) GetCampaignStats(ctx context.Context, campaignID string) (*CampaignStats, error) {
	_ = ctx

	// Fetch campaign
	campaign, err := a.campaignRepo.GetByID(campaignID)
	if err != nil || campaign == nil {
		return nil, errors.New("campaign not found")
	}

	// Fetch events for campaign
	events, err := a.eventRepo.GetByCampaignID(campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campaign events: %w", err)
	}

	stats := &CampaignStats{
		CampaignID:   campaignID,
		CampaignName: campaign.Name,
	}

	for _, event := range events {
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
	_ = ctx

	campaigns, err := a.campaignRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campaigns: %w", err)
	}

	report := map[string]interface{}{
		"total_campaigns": len(campaigns),
		"campaigns":       []interface{}{},
	}

	campaignsList := report["campaigns"].([]interface{})
	for _, campaign := range campaigns {
		campaignsList = append(campaignsList, campaign)
	}
	report["campaigns"] = campaignsList

	return report, nil
}

// GetDailyReport returns daily analytics report (simplified)
func (a *AnalyticsUseCase) GetDailyReport(ctx context.Context, date string) (map[string]interface{}, error) {
	_ = ctx

	day, err := normalizeDailyDate(date)
	if err != nil {
		return nil, err
	}

	events, err := a.eventRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}

	aggregate, totalEvents, totalRevenue := aggregateEvents(events, func(t time.Time) bool {
		return t.UTC().Format("2006-01-02") == day
	})

	return map[string]interface{}{
		"date":          day,
		"total_events":  totalEvents,
		"total_revenue": totalRevenue,
		"top_campaigns": a.buildTopCampaigns(aggregate),
	}, nil
}

// GetWeeklyReport returns weekly analytics report (simplified)
func (a *AnalyticsUseCase) GetWeeklyReport(ctx context.Context, weekStart string) (map[string]interface{}, error) {
	_ = ctx

	weekStartDate, weekEndDate, err := normalizeWeeklyRange(weekStart)
	if err != nil {
		return nil, err
	}

	events, err := a.eventRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}

	aggregate, totalEvents, totalRevenue := aggregateEvents(events, func(t time.Time) bool {
		utc := t.UTC()
		return !utc.Before(weekStartDate) && utc.Before(weekEndDate)
	})

	return map[string]interface{}{
		"week_start":    weekStartDate.Format("2006-01-02"),
		"week_end":      weekEndDate.Add(-time.Nanosecond).Format("2006-01-02"),
		"total_events":  totalEvents,
		"total_revenue": totalRevenue,
		"top_campaigns": a.buildTopCampaigns(aggregate),
	}, nil
}

// GetMonthlyReport returns monthly analytics report (simplified)
func (a *AnalyticsUseCase) GetMonthlyReport(ctx context.Context, month string) (map[string]interface{}, error) {
	_ = ctx

	monthLabel, monthStart, monthEnd, err := normalizeMonthlyRange(month)
	if err != nil {
		return nil, err
	}

	events, err := a.eventRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}

	aggregate, totalEvents, totalRevenue := aggregateEvents(events, func(t time.Time) bool {
		utc := t.UTC()
		return !utc.Before(monthStart) && utc.Before(monthEnd)
	})

	return map[string]interface{}{
		"month":         monthLabel,
		"total_events":  totalEvents,
		"total_revenue": totalRevenue,
		"top_campaigns": a.buildTopCampaigns(aggregate),
	}, nil
}

type campaignAggregate struct {
	CampaignID   string
	Events       int64
	Impressions  int64
	Clicks       int64
	Conversions  int64
	TotalRevenue float64
}

func aggregateEvents(events []domain.Event, include func(time.Time) bool) (map[string]*campaignAggregate, int64, float64) {
	aggregate := make(map[string]*campaignAggregate)
	totalEvents := int64(0)
	totalRevenue := float64(0)

	for _, event := range events {
		eventTime, ok := parseEventTimestamp(event.Timestamp)
		if !ok {
			continue
		}
		if !include(eventTime) {
			continue
		}

		totalEvents++

		agg, exists := aggregate[event.CampaignID]
		if !exists {
			agg = &campaignAggregate{CampaignID: event.CampaignID}
			aggregate[event.CampaignID] = agg
		}

		agg.Events++
		switch event.EventType {
		case domain.EventType("impressions"):
			agg.Impressions++
		case domain.EventType("clicks"):
			agg.Clicks++
		case domain.EventType("conversions"):
			agg.Conversions++
			agg.TotalRevenue += event.Metadata.Amount
			totalRevenue += event.Metadata.Amount
		}
	}

	return aggregate, totalEvents, totalRevenue
}

func parseEventTimestamp(value string) (time.Time, bool) {
	if value == "" {
		return time.Time{}, false
	}
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed, true
	}
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed, true
	}
	return time.Time{}, false
}

func normalizeDailyDate(date string) (string, error) {
	if date == "" || date == "today" {
		return time.Now().UTC().Format("2006-01-02"), nil
	}
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}
	return parsed.UTC().Format("2006-01-02"), nil
}

func normalizeWeeklyRange(weekStart string) (time.Time, time.Time, error) {
	if weekStart == "" || weekStart == "this_week" {
		now := time.Now().UTC()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		return start, start.AddDate(0, 0, 7), nil
	}

	parsed, err := time.Parse("2006-01-02", weekStart)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid week_start format, expected YYYY-MM-DD")
	}
	start := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
	return start, start.AddDate(0, 0, 7), nil
}

func normalizeMonthlyRange(month string) (string, time.Time, time.Time, error) {
	if month == "" || month == "this_month" {
		now := time.Now().UTC()
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0)
		return start.Format("2006-01"), start, end, nil
	}

	parsed, err := time.Parse("2006-01", month)
	if err != nil {
		return "", time.Time{}, time.Time{}, fmt.Errorf("invalid month format, expected YYYY-MM")
	}
	start := time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	return start.Format("2006-01"), start, end, nil
}

func (a *AnalyticsUseCase) buildTopCampaigns(aggregate map[string]*campaignAggregate) []interface{} {
	if len(aggregate) == 0 {
		return []interface{}{}
	}

	campaignNameMap := make(map[string]string)
	if campaigns, err := a.campaignRepo.GetAll(); err == nil {
		for _, campaign := range campaigns {
			campaignNameMap[campaign.ID] = campaign.Name
		}
	}

	items := make([]*campaignAggregate, 0, len(aggregate))
	for _, agg := range aggregate {
		items = append(items, agg)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Events == items[j].Events {
			return items[i].TotalRevenue > items[j].TotalRevenue
		}
		return items[i].Events > items[j].Events
	})

	limit := 5
	if len(items) < limit {
		limit = len(items)
	}

	topCampaigns := make([]interface{}, 0, limit)
	for i := 0; i < limit; i++ {
		item := items[i]
		topCampaigns = append(topCampaigns, map[string]interface{}{
			"campaign_id":     item.CampaignID,
			"campaign_name":   campaignNameMap[item.CampaignID],
			"events":          item.Events,
			"impressions":     item.Impressions,
			"clicks":          item.Clicks,
			"conversions":     item.Conversions,
			"total_revenue":   item.TotalRevenue,
			"conversion_rate": conversionRate(item.Clicks, item.Conversions),
		})
	}

	return topCampaigns
}

func conversionRate(clicks, conversions int64) float64 {
	if clicks == 0 {
		return 0
	}
	return float64(conversions) / float64(clicks)
}
