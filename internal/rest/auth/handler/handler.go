package handler

import (
	"encoding/json"
	"net/http"
)

// RegisterRequest represents user registration payload
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

// LoginRequest represents user login payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents auth response with token
type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	// TODO: Validate email format
	// TODO: Hash password
	// TODO: Save user to persistence
	// TODO: Generate JWT token

	response := map[string]interface{}{
		"message": "User registered successfully",
		"token":   "jwt_token_here",
	}
	json.NewEncoder(w).Encode(response)
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	// TODO: Validate credentials
	// TODO: Generate JWT token

	response := map[string]interface{}{
		"message": "Login successful",
		"token":   "jwt_token_here",
	}
	json.NewEncoder(w).Encode(response)
}
