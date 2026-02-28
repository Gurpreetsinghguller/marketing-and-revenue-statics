package handler

import (
	"context"
	"encoding/json"
	"net/http"

	auth_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/auth/usecase"
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

// AuthHandler handles authentication requests
type AuthHandler struct {
	usecase *auth_usecase.AuthUseCase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(uc *auth_usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		usecase: uc,
	}
}

// RegisterHandler handles user registration
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, `{"error": "Email, password, and name are required"}`, http.StatusBadRequest)
		return
	}

	// Call usecase
	usecaseReq := auth_usecase.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Role:     req.Role,
	}

	result, err := h.usecase.Register(context.Background(), usecaseReq)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"message": "User registered successfully",
		"user_id": result.UserID,
		"token":   result.Token,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// LoginHandler handles user login
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error": "Email and password are required"}`, http.StatusBadRequest)
		return
	}

	// Call usecase
	usecaseReq := auth_usecase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.usecase.Login(context.Background(), usecaseReq)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"message": "Login successful",
		"token":   result.Token,
		"user":    result.User,
	}

	json.NewEncoder(w).Encode(response)
}
