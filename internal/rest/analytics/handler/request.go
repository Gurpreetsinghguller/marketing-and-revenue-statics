package handler

type DailyReportRequest struct {
	Date string `json:"date"`
}

type WeeklyReportRequest struct {
	WeekStart string `json:"week_start"`
}

type MonthlyReportRequest struct {
	Month string `json:"month"`
}

type CampaignStatsRequest struct {
	CampaignID string `json:"campaign_id"`
}
