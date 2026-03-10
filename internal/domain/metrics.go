package domain

import "time"

// CampaignMetrics represents real-time aggregated metrics for a campaign
type CampaignMetrics struct {
	CampaignID     string    `json:"campaign_id"`
	Date           string    `json:"date"` // YYYY-MM-DD
	Hour           string    `json:"hour"` // YYYY-MM-DD HH:00:00 (optional, for hourly metrics)
	Impressions    int64     `json:"impressions"`
	Clicks         int64     `json:"clicks"`
	Conversions    int64     `json:"conversions"`
	TotalRevenue   float64   `json:"total_revenue"`
	CTR            float64   `json:"ctr"`             // Clicks / Impressions
	ConversionRate float64   `json:"conversion_rate"` // Conversions / Clicks
	AvgRevenue     float64   `json:"avg_revenue"`     // Revenue per conversion
	LastUpdated    time.Time `json:"last_updated"`
}

// ChannelMetrics groups metrics by source/channel
type ChannelMetrics struct {
	Channel    string          `json:"channel"` // e.g., "facebook", "google"
	CampaignID string          `json:"campaign_id"`
	Date       string          `json:"date"`
	Metrics    CampaignMetrics `json:"metrics"`
}

// DeviceMetrics groups metrics by device type
type DeviceMetrics struct {
	Device     string          `json:"device"` // e.g., "mobile", "desktop"
	CampaignID string          `json:"campaign_id"`
	Date       string          `json:"date"`
	Metrics    CampaignMetrics `json:"metrics"`
}

// MetricsRepo handles persistent storage and retrieval of pre-aggregated metrics
type MetricsRepo interface {
	// SaveMetrics persists computed metrics
	SaveMetrics(metrics *CampaignMetrics) error

	// GetMetrics retrieves metrics for a campaign-date
	GetMetrics(campaignID, date string) (*CampaignMetrics, error)

	// GetMetricsByHour retrieves hourly metrics
	GetMetricsByHour(campaignID, date, hour string) (*CampaignMetrics, error)

	// GetMetricsRange retrieves metrics across a date range
	GetMetricsRange(campaignID, startDate, endDate string) ([]CampaignMetrics, error)

	// GetMetricsForAllCampaigns retrieves metrics for all campaigns on a date
	GetMetricsForAllCampaigns(date string) ([]CampaignMetrics, error)

	// GetChannelMetrics retrieves metrics grouped by channel
	GetChannelMetrics(campaignID, date string) ([]ChannelMetrics, error)

	// GetDeviceMetrics retrieves metrics grouped by device
	GetDeviceMetrics(campaignID, date string) ([]DeviceMetrics, error)
}
