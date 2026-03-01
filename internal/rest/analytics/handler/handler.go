package handler

import (
	"context"
	"encoding/json"
	"net/http"

	analytics_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/analytics/usecase"
	"github.com/gorilla/mux"
)

type AnalyticsHandler struct {
	usecase analytics_usecase.AnalyticsUseCaseInterface
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(uc analytics_usecase.AnalyticsUseCaseInterface) *AnalyticsHandler {
	return &AnalyticsHandler{
		usecase: uc,
	}
}

// GetAnalyticsReportHandler retrieves aggregated analytics report
func (h *AnalyticsHandler) GetAnalyticsReportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Call usecase
	report, err := h.usecase.GetAnalyticsReport(context.Background())
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(report)
}

// GetDailyReportHandler retrieves daily performance summary
func (h *AnalyticsHandler) GetDailyReportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Extract date from query params
	date := r.URL.Query().Get("date")
	if date == "" {
		date = "today"
	}

	// Call usecase
	report, err := h.usecase.GetDailyReport(r.Context(), date)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(report)
}

// GetWeeklyReportHandler retrieves weekly performance summary
func (h *AnalyticsHandler) GetWeeklyReportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Extract week from query params
	weekStart := r.URL.Query().Get("week_start")
	if weekStart == "" {
		weekStart = "this_week"
	}

	// Call usecase
	report, err := h.usecase.GetWeeklyReport(context.Background(), weekStart)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(report)
}

// GetMonthlyReportHandler retrieves monthly performance summary
func (h *AnalyticsHandler) GetMonthlyReportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Extract month from query params
	month := r.URL.Query().Get("month")
	if month == "" {
		month = "this_month"
	}

	// Call usecase
	report, err := h.usecase.GetMonthlyReport(context.Background(), month)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(report)
}

// GetPublicStatsHandler retrieves anonymized public campaign stats
func (h *AnalyticsHandler) GetPublicStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Call usecase (no auth required for public stats)
	stats, err := h.usecase.GetPublicStats(context.Background())
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

// GetCampaignStatsHandler retrieves campaign statistics by campaign ID
func (h *AnalyticsHandler) GetCampaignStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	campaignID := mux.Vars(r)["campaign_id"]
	if campaignID == "" {
		http.Error(w, `{"error": "campaign_id is required"}`, http.StatusBadRequest)
		return
	}

	stats, err := h.usecase.GetCampaignStats(context.Background(), campaignID)
	if err != nil {
		if err.Error() == "campaign not found" {
			http.Error(w, `{"error": "campaign not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}
