package handler

// AuthResponse represents auth response with token
type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
