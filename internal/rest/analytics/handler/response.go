package handler

import "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"

// ===== Report Responses =====

// ReportResponse represents a single analytics report
type ReportResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    *domain.Report `json:"data"`
}

// ReportListResponse represents a list of reports
type ReportListResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    []domain.Report `json:"data"`
	Count   int             `json:"count"`
}

// ===== Metrics Responses =====

// MetricsResponse represents campaign metrics for a period
type MetricsResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    *domain.Metrics `json:"data"`
}

// ===== Breakdown Responses =====

// ChannelBreakdownResponse represents metrics grouped by channel
type ChannelBreakdownResponse struct {
	Status  string                    `json:"status"`
	Message string                    `json:"message"`
	Data    []domain.ChannelBreakdown `json:"data"`
}

// DeviceBreakdownResponse represents metrics grouped by device
type DeviceBreakdownResponse struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    []domain.DeviceBreakdown `json:"data"`
}

// ===== Public Stats Response =====

// PublicStatsResponse represents anonymized public campaign stats
type PublicStatsResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    *PublicStats `json:"data"`
}

// PublicStats represents anonymized campaign statistics
type PublicStats struct {
	CampaignName  string                    `json:"campaign_name"`
	Period        string                    `json:"period"`
	TotalClicks   int64                     `json:"total_clicks"`
	TotalRevenue  float64                   `json:"total_revenue"`
	AvgConversion float64                   `json:"avg_conversion_rate"`
	TopChannels   []domain.ChannelBreakdown `json:"top_channels"`
}

// ===== Aggregated Summary Response =====

// SummaryResponse represents high-level analytics summary
type SummaryResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Data    *Summary `json:"data"`
}

// Summary contains aggregated summary metrics
type Summary struct {
	TotalCampaigns   int64           `json:"total_campaigns"`
	TotalImpressions int64           `json:"total_impressions"`
	TotalClicks      int64           `json:"total_clicks"`
	TotalConversions int64           `json:"total_conversions"`
	TotalRevenue     float64         `json:"total_revenue"`
	TotalSpend       float64         `json:"total_spend"`
	AverageMetrics   *domain.Metrics `json:"average_metrics"`
	TopCampaigns     []domain.Report `json:"top_campaigns"`
}

// ===== Error Response =====

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// ===== Success Response (Generic) =====

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
