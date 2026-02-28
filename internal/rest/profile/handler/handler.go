package handlers

import (
	"encoding/json"
	"net/http"
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

// GetProfileHandler retrieves user profile
func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO: Extract user ID from JWT token
	// TODO: Fetch profile from persistence

	profile := UserProfile{}
	json.NewEncoder(w).Encode(profile)
}

// UpdateProfileHandler updates user profile
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var profile UserProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	// TODO: Extract user ID from JWT token
	// TODO: Validate profile data
	// TODO: Update in persistence

	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}
