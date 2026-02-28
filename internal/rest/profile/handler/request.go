package handler

// UpdateProfileRequest represents profile update input
type UpdateProfileRequest struct {
	Name    string `json:"name"`
	Bio     string `json:"bio"`
	Phone   string `json:"phone"`
	Picture string `json:"picture"`
}
