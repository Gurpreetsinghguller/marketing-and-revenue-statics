package handler

// TODO: request validation
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

// TODO: request validation
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
