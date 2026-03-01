package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	campaign_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/campaigns/usecase"
)

type CampaignHandler struct {
	usecase *campaign_usecase.CampaignUseCase
}

var requestValidator = validator.New()

func NewCampaignHandler(uc *campaign_usecase.CampaignUseCase) *CampaignHandler {
	return &CampaignHandler{usecase: uc}
}

func (h *CampaignHandler) CreateCampaignHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
		return
	}

	var req CreateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Channel = strings.TrimSpace(req.Channel)

	if err := requestValidator.Struct(req); err != nil {
		http.Error(w, `{"error": "`+formatValidationError(err)+`"}`, http.StatusBadRequest)
		return
	}

	start := time.Now().UTC()
	campaign := &domain.Campaign{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID,
		Status:      domain.CampaignStatusActive,
		DateRange: domain.DateRange{
			Start: &start,
			End:   &time.Time{},
		},
		Budget:    req.Budget,
		Channel:   req.Channel,
		IsPublic:  req.IsPublic,
		UpdatedAt: time.Now(),
	}

	created, err := h.usecase.CreateCampaign(context.Background(), campaign, userID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *CampaignHandler) GetCampaignsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := strings.TrimSpace(r.Header.Get("X-User-ID"))
	userRole := strings.TrimSpace(r.Header.Get("X-User-Role"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	budgetQuery := strings.TrimSpace(r.URL.Query().Get("budget"))
	isPublicQuery := strings.TrimSpace(r.URL.Query().Get("is_public"))

	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	if userRole == "" {
		http.Error(w, `{"error": "User role not found"}`, http.StatusForbidden)
		return
	}

	var budget *float64
	if budgetQuery != "" {
		parsedBudget, err := strconv.ParseFloat(budgetQuery, 64)
		if err != nil {
			http.Error(w, `{"error": "Invalid budget query value"}`, http.StatusBadRequest)
			return
		}
		budget = &parsedBudget
	}

	var isPublic *bool
	if isPublicQuery != "" {
		parsedIsPublic, err := strconv.ParseBool(isPublicQuery)
		if err != nil {
			http.Error(w, `{"error": "Invalid is_public query value"}`, http.StatusBadRequest)
			return
		}
		isPublic = &parsedIsPublic
	}

	campaigns, err := h.usecase.GetCampaignsWithFilters(
		context.Background(),
		userID,
		userRole,
		status,
		budget,
		isPublic,
	)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(campaigns)
}

func (h *CampaignHandler) UpdateCampaignHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("X-User-ID")
	userRole := strings.TrimSpace(r.Header.Get("X-User-Role"))
	if userID == "" {
		http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
		return
	}

	campaignID := mux.Vars(r)["id"]
	if campaignID == "" {
		http.Error(w, `{"error": "Campaign ID is required"}`, http.StatusBadRequest)
		return
	}

	var req UpdateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	if req.Name != nil {
		trimmedName := strings.TrimSpace(*req.Name)
		req.Name = &trimmedName
	}
	if req.Status != nil {
		normalizedStatus := strings.ToLower(strings.TrimSpace(*req.Status))
		req.Status = &normalizedStatus
	}

	if err := requestValidator.Struct(req); err != nil {
		http.Error(w, `{"error": "`+formatValidationError(err)+`"}`, http.StatusBadRequest)
		return
	}

	updatedCampaign := &domain.Campaign{
		Name:      strings.TrimSpace(*req.Name),
		Status:    domain.CampaignStatus(strings.ToLower(strings.TrimSpace(*req.Status))),
		UpdatedAt: time.Now(),
	}

	if req.Description != nil {
		updatedCampaign.Description = *req.Description
	}
	if req.Budget != nil {
		updatedCampaign.Budget = *req.Budget
	}
	if req.IsPublic != nil {
		updatedCampaign.IsPublic = *req.IsPublic
	}

	campaign, err := h.usecase.UpdateCampaign(context.Background(), campaignID, updatedCampaign, userID, userRole)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(campaign)
}

func (h *CampaignHandler) DeleteCampaignHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("X-User-ID")
	userRole := strings.TrimSpace(r.Header.Get("X-User-Role"))
	if userID == "" {
		http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
		return
	}

	campaignID := mux.Vars(r)["id"]
	if campaignID == "" {
		http.Error(w, `{"error": "Campaign ID is required"}`, http.StatusBadRequest)
		return
	}

	if err := h.usecase.DeleteCampaignWithRole(context.Background(), campaignID, userID, userRole); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Campaign deleted successfully"})
}

func (h *CampaignHandler) GetCampaignPreviewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	campaignID := mux.Vars(r)["id"]
	if campaignID == "" {
		http.Error(w, `{"error": "Campaign ID is required"}`, http.StatusBadRequest)
		return
	}

	campaign, err := h.usecase.GetCampaignByID(context.Background(), campaignID)
	if err != nil {
		http.Error(w, `{"error": "Campaign not found"}`, http.StatusNotFound)
		return
	}

	if !campaign.IsPublic {
		http.Error(w, `{"error": "Campaign is not public"}`, http.StatusForbidden)
		return
	}

	response := CampaignPreviewResponse{
		Name:        campaign.Name,
		Description: campaign.Description,
		Status:      string(campaign.Status),
	}
	json.NewEncoder(w).Encode(response)
}

func (h *CampaignHandler) PatchCampaignStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Header.Get("X-User-ID")
	userRole := strings.TrimSpace(r.Header.Get("X-User-Role"))
	if userID == "" {
		http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
		return
	}

	campaignID := mux.Vars(r)["id"]
	if campaignID == "" {
		http.Error(w, `{"error": "Campaign ID is required"}`, http.StatusBadRequest)
		return
	}

	var req PatchCampaignStatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	req.Status = strings.ToLower(strings.TrimSpace(req.Status))

	if err := requestValidator.Struct(req); err != nil {
		http.Error(w, `{"error": "`+formatValidationError(err)+`"}`, http.StatusBadRequest)
		return
	}

	campaign, err := h.usecase.PatchCampaignStatus(context.Background(), campaignID, userID, userRole, req.Status)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(campaign)
}

func (h *CampaignHandler) EndCampaignHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := strings.TrimSpace(r.Header.Get("X-User-ID"))
	userRole := strings.TrimSpace(r.Header.Get("X-User-Role"))
	if userID == "" {
		http.Error(w, `{"error": "User not authorized"}`, http.StatusUnauthorized)
		return
	}

	campaignID := mux.Vars(r)["id"]
	if campaignID == "" {
		http.Error(w, `{"error": "Campaign ID is required"}`, http.StatusBadRequest)
		return
	}

	campaign, err := h.usecase.EndCampaign(context.Background(), campaignID, userID, userRole)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "not found"):
			http.Error(w, `{"error": "`+errMsg+`"}`, http.StatusNotFound)
		case strings.Contains(errMsg, "not authorized"):
			http.Error(w, `{"error": "`+errMsg+`"}`, http.StatusForbidden)
		case strings.Contains(errMsg, "only active or paused"):
			http.Error(w, `{"error": "`+errMsg+`"}`, http.StatusBadRequest)
		default:
			http.Error(w, `{"error": "`+errMsg+`"}`, http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(campaign)
}

func (h *CampaignHandler) SearchCampaignsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := strings.TrimSpace(r.Header.Get("X-User-ID"))
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, `{"error": "Search query is required"}`, http.StatusBadRequest)
		return
	}

	campaigns, err := h.usecase.SearchCampaigns(context.Background(), query)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(strings.ToLower(errMsg), "invalid search pattern") {
			http.Error(w, `{"error": "`+errMsg+`"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"error": "`+errMsg+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(campaigns)
}

func formatValidationError(err error) string {
	if err == nil {
		return "invalid request payload"
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) && len(validationErrors) > 0 {
		field := strings.ToLower(validationErrors[0].Field())
		tag := validationErrors[0].Tag()

		switch tag {
		case "required":
			return field + " is required"
		case "oneof":
			return field + " has invalid value"
		case "min":
			return field + " cannot be empty"
		case "gte":
			return field + " must be greater than or equal to 0"
		case "gtefield":
			return field + " must be greater than or equal to start"
		default:
			return field + " is invalid"
		}
	}

	return "invalid request payload"
}
