package domain

// Metrics represents calculated campaign metrics
type Metrics struct {
	CTR            float64 `json:"ctr"`             // Click Through Rate: (clicks / impressions) * 100
	CPC            float64 `json:"cpc"`             // Cost Per Click: spend / clicks
	ROI            float64 `json:"roi"`             // Return on Investment: ((revenue - spend) / spend) * 100
	ConversionRate float64 `json:"conversion_rate"` // (conversions / clicks) * 100
}

// ChannelBreakdown represents metrics grouped by channel/source
type ChannelBreakdown struct {
	Channel        string  `json:"channel"`
	Impressions    int64   `json:"impressions"`
	Clicks         int64   `json:"clicks"`
	Conversions    int64   `json:"conversions"`
	Revenue        float64 `json:"revenue"`
	CTR            float64 `json:"ctr"`
	ConversionRate float64 `json:"conversion_rate"`
}

// DeviceBreakdown represents metrics grouped by device type
type DeviceBreakdown struct {
	Device         string  `json:"device"`
	Impressions    int64   `json:"impressions"`
	Clicks         int64   `json:"clicks"`
	Conversions    int64   `json:"conversions"`
	Revenue        float64 `json:"revenue"`
	CTR            float64 `json:"ctr"`
	ConversionRate float64 `json:"conversion_rate"`
}

// Report represents aggregated analytics data
type Report struct {
	// Period information
	CampaignID   string `json:"campaign_id"`
	CampaignName string `json:"campaign_name"`
	Period       string `json:"period"` // e.g., "2026-02-01 to 2026-02-28"
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	ReportType   string `json:"report_type"` // "custom", "daily", "weekly", "monthly"

	// Aggregated counts
	TotalImpressions int64 `json:"total_impressions"`
	TotalClicks      int64 `json:"total_clicks"`
	TotalConversions int64 `json:"total_conversions"`

	// Financial metrics
	TotalRevenue float64 `json:"total_revenue"`
	TotalSpend   float64 `json:"total_spend"`

	// Calculated metrics
	Metrics Metrics `json:"metrics"`

	// Breakdown by channel/source
	ByChannel []ChannelBreakdown `json:"by_channel"`

	// Breakdown by device
	ByDevice []DeviceBreakdown `json:"by_device"`
}

// ReportRepo defines database operations for Report entity
type ReportRepo interface {
	// Create saves a pre-aggregated report
	Create(report *Report) error

	// GetByID retrieves a report by ID
	GetByID(id string) (*Report, error)

	// GetByCampaignAndDateRange retrieves report for a campaign and date range
	GetByCampaignAndDateRange(campaignID, startDate, endDate string) (*Report, error)

	// GetDaily retrieves daily aggregated report for a campaign
	GetDaily(campaignID, date string) (*Report, error)

	// GetWeekly retrieves weekly aggregated report for a campaign
	GetWeekly(campaignID, year, week string) (*Report, error)

	// GetMonthly retrieves monthly aggregated report for a campaign
	GetMonthly(campaignID, year, month string) (*Report, error)

	// GetAll retrieves all reports
	GetAll() ([]Report, error)
}
