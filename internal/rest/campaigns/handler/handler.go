package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	campaign_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/campaigns/usecase"
)

// Campaign represents a marketing campaign
type Campaign struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	Channel   string  `json:"channel"`
	Budget    float64 `json:"budget"`
	CreatedBy string  `json:"created_by"`
	Public    bool    `json:"public"`
}

// CreateCampaignRequest represents campaign creation input
type CreateCampaignRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	DateRange   domain.DateRange `json:"date_range"`
	Budget      float64          `json:"budget"`
	Channel     string           `json:"channel"`
	IsPublic    bool             `json:"is_public"`
}

// CampaignHandler handles HTTP requests for campaign operations
type CampaignHandler struct {
	usecase *campaign_usecase.CampaignUseCase
}

// NewCampaignHandler creates and returns a new CampaignHandler
func NewCampaignHandler(uc *campaign_usecase.CampaignUseCase) *CampaignHandler {
	return &CampaignHandler{usecase: uc}
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(r *http.Request) string {
	// In production, extract from JWT token in middleware
	return r.Header.Get("X-User-ID") // Simplified for demo
}

// GetCampaignsHandler retrieves campaigns with filters
func (h *CampaignHandler) GetCampaignsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Extract query parameters
	status := r.URL.Query().Get("status")
	channel := r.URL.Query().Get("channel")

	var campaigns []domain.Campaign
	var err error

	// Apply filters
	if status != "" {
		campaigns, err = h.usecase.GetCampaignsByStatus(context.Background(), status)
	} else if channel != "" {
		campaigns, err = h.usecase.GetCampaignsByChannel(context.Background(), channel)
	} else {
		campaigns, err = h.usecase.GetAllCampaigns(context.Background())
	}

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(campaigns)
}

// CreateCampaignHandler creates a new campaign
func (h *CampaignHandler) CreateCampaignHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	var req CreateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		http.Error(w, `{"error": "Campaign name is required"}`, http.StatusBadRequest)
		return
	}

	// Call usecase
	usecaseReq := campaign_usecase.CreateCampaignRequest{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		DateRange:   req.DateRange,
		Budget:      req.Budget,
		Channel:     req.Channel,
		IsPublic:    req.IsPublic,
	}

	campaign, err := h.usecase.CreateCampaign(context.Background(), usecaseReq, userID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(campaign)
}

// UpdateCampaignHandler updates campaign details
func (h *CampaignHandler) UpdateCampaignHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignID := vars["id"]
	if campaignID == "" {
		http.Error(w, `{"error": "Campaign ID is required"}`, http.StatusBadRequest)
		return
	}

	var req CreateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	// Call usecase
	usecaseReq := campaign_usecase.CreateCampaignRequest{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		DateRange:   req.DateRange,
		Budget:      req.Budget,
		Channel:     req.Channel,
		IsPublic:    req.IsPublic,
	}

	campaign, err := h.usecase.UpdateCampaign(context.Background(), campaignID, usecaseReq, userID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(campaign)
}

// DeleteCampaignHandler deletes a campaign
func (h *CampaignHandler) DeleteCampaignHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignID := vars["id"]
	if campaignID == "" {
		http.Error(w, `{"error": "Campaign ID is required"}`, http.StatusBadRequest)
		return
	}

	// Call usecase
	err := h.usecase.DeleteCampaign(context.Background(), campaignID, userID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Campaign deleted successfully"})
}

// GetCampaignPreviewHandler retrieves public campaign preview
func (h *CampaignHandler) GetCampaignPreviewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract campaign ID from URL
	vars := mux.Vars(r)
	campaignID := vars["id"]
	if campaignID == "" {
		http.Error(w, `{"error": "Campaign ID is required"}`, http.StatusBadRequest)
		return
	}

	// Fetch campaign
	campaign, err := h.usecase.GetCampaignByID(context.Background(), campaignID)
	if err != nil {
		http.Error(w, `{"error": "Campaign not found"}`, http.StatusNotFound)
		return
	}

	// Check if public
	if !campaign.IsPublic {
		http.Error(w, `{"error": "Campaign is not public"}`, http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(campaign)
}

// SearchCampaignsHandler searches campaigns by name/description
func (h *CampaignHandler) SearchCampaignsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := GetUserIDFromContext(r)
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Extract search query
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, `{"error": "Search query is required"}`, http.StatusBadRequest)
		return
	}

	// Call usecase
	campaigns, err := h.usecase.SearchCampaigns(context.Background(), query)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(campaigns)
}
