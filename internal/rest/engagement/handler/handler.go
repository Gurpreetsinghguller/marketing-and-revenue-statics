package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	engagement_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/engagement/usecase"
)

// EngagementHandler handles engagement requests
type EngagementHandler struct {
	usecase *engagement_usecase.EngagementUseCase
}

// NewEngagementHandler creates a new engagement handler
func NewEngagementHandler(uc *engagement_usecase.EngagementUseCase) *EngagementHandler {
	return &EngagementHandler{
		usecase: uc,
	}
}

// GetUserEngagementHandler retrieves user's engagement across all campaigns
func (h *EngagementHandler) GetUserEngagementHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract user_id from URL
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		http.Error(w, `{"error": "User ID is required"}`, http.StatusBadRequest)
		return
	}

	// Check authorization
	authenticatedUserID := r.Header.Get("X-User-ID")
	if authenticatedUserID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Call usecase
	engagement, err := h.usecase.GetUserEngagement(context.Background(), userID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(engagement)
}

// GetCampaignFunnelHandler retrieves funnel/drop-off data for a campaign
func (h *EngagementHandler) GetCampaignFunnelHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract campaign_id from URL
	vars := mux.Vars(r)
	campaignID := vars["campaign_id"]
	if campaignID == "" {
		http.Error(w, `{"error": "Campaign ID is required"}`, http.StatusBadRequest)
		return
	}

	// Check authorization
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Call usecase
	funnel, err := h.usecase.GetCampaignFunnel(context.Background(), campaignID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(funnel)
}

// GetUserCampaignEngagementHandler retrieves user's engagement with specific campaign
func (h *EngagementHandler) GetUserCampaignEngagementHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract user_id and campaign_id from URL
	vars := mux.Vars(r)
	userID := vars["user_id"]
	campaignID := vars["campaign_id"]

	if userID == "" || campaignID == "" {
		http.Error(w, `{"error": "User ID and Campaign ID are required"}`, http.StatusBadRequest)
		return
	}

	// Check authorization
	authenticatedUserID := r.Header.Get("X-User-ID")
	if authenticatedUserID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Call usecase
	engagement, err := h.usecase.GetUserCampaignEngagement(context.Background(), userID, campaignID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(engagement)
}
