package domain

import "time"

type CampaignStatus string

const (
	CampaignStatusActive    CampaignStatus = "active"
	CampaignStatusPaused    CampaignStatus = "paused"
	CampaignStatusCompleted CampaignStatus = "completed"
	CampaignStatusInactive  CampaignStatus = "inactive"
)

// from the requirements, we need a Campaign entity with fields like ID, Name, Description, Status, DateRange, Channel, CreatedBy, IsPublic
type Campaign struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      CampaignStatus `json:"status"` // Active, Paused, Completed
	DateRange   DateRange      `json:"date_range"`
	Budget      float64        `json:"budget"`
	Channel     string         `json:"channel"` // Email, Social Media, etc.
	CreatedBy   string         `json:"created_by"`
	IsPublic    bool           `json:"is_public"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type DateRange struct {
	Start *time.Time `json:"start"`
	End   *time.Time `json:"end"`
}

// This will be implemented by persistence layer (e.g., database) to manage campaign data
type CampaignRepo interface {
	// Create saves a new campaign
	Create(campaign *Campaign) error

	// GetByID retrieves a campaign by ID
	GetByID(id string) (*Campaign, error)

	// GetAll retrieves all campaigns
	GetAll() ([]Campaign, error)

	// GetPublic retrieves all public campaigns (accessible without auth)
	GetPublic() ([]Campaign, error)

	// Update updates an existing campaign
	Update(campaign *Campaign) error

	// Delete removes a campaign
	Delete(id string) error

	// Search searches campaigns by name or description
	Search(query string) ([]Campaign, error)

	// GetWithFilters retrieves campaigns with multiple filters applied
	// Filters: status, channel, date_range, etc.
	GetWithFilters(filters map[string]interface{}) ([]Campaign, error)
}
