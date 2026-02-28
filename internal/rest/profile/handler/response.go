package handler

// UserProfileResponse represents user profile information
type UserProfileResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Bio     string `json:"bio"`
	Phone   string `json:"phone"`
	Picture string `json:"picture"`
}
