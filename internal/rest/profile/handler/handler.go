package handler

import (
	"context"
	"encoding/json"
	"net/http"

	profile_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/profile/usecase"
)

// UserProfile represents user profile information
type UserProfile struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Bio     string `json:"bio"`
	Phone   string `json:"phone"`
	Picture string `json:"picture"`
}

// UpdateProfileRequest represents profile update input
type UpdateProfileRequest struct {
	Name    string `json:"name"`
	Bio     string `json:"bio"`
	Phone   string `json:"phone"`
	Picture string `json:"picture"`
}

// ProfileHandler handles profile requests
type ProfileHandler struct {
	usecase *profile_usecase.ProfileUseCase
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(uc *profile_usecase.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{
		usecase: uc,
	}
}

// GetProfileHandler retrieves user profile
func (h *ProfileHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract user ID from context
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Call usecase
	user, err := h.usecase.GetProfile(context.Background(), userID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	profile := UserProfile{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Role:    string(user.Role),
		Bio:     user.Bio,
		Phone:   user.Phone,
		Picture: user.Picture,
	}

	json.NewEncoder(w).Encode(profile)
}

// UpdateProfileHandler updates user profile
func (h *ProfileHandler) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract user ID from context
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	// Call usecase
	usecaseReq := profile_usecase.UpdateProfileRequest{
		Name:    req.Name,
		Bio:     req.Bio,
		Phone:   req.Phone,
		Picture: req.Picture,
	}

	user, err := h.usecase.UpdateProfile(context.Background(), userID, usecaseReq)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	profile := UserProfile{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Role:    string(user.Role),
		Bio:     user.Bio,
		Phone:   user.Phone,
		Picture: user.Picture,
	}

	json.NewEncoder(w).Encode(profile)
}
