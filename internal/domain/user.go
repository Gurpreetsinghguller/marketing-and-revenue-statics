package domain

import "time"

type Role string

const (
	AdminRole    Role = "admin"
	MarketerRole Role = "marketer"
	AnalystRole  Role = "analyst"
)

// User represents a user entity with profile and auth info
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose in API responses
	Name      string    `json:"name"`
	Role      Role      `json:"role"` // Admin, Marketer, Analyst
	Bio       string    `json:"bio"`
	Phone     string    `json:"phone"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserInteractionCampaign struct {
	CampaignID   string `json:"campaign_id"`
	CampaignName string `json:"campaign_name"`
	Duration     int64  `json:"duration"`     // Time spent on campaign pages in seconds
	Interactions int64  `json:"interactions"` // Total interactions (clicks + impressions + conversions)
}

// UserRepo defines database operations for User entity
type UserRepo interface {
	// Create saves a new user to the database
	Create(user *User) error

	// GetByID retrieves a user by ID
	GetByID(id string) (*User, error)

	// GetByEmail retrieves a user by email (for login)
	GetByEmail(email string) (*User, error)

	// Update updates an existing user
	Update(user *User) error

	// Delete removes a user from the database
	Delete(id string) error

	// GetAll retrieves all users (for admin purposes)
	GetAll() ([]User, error)

	// GetByRole retrieves all users with a specific role
	GetByRole(role string) ([]User, error)

	// EmailExists checks if an email already exists
	EmailExists(email string) bool
}
