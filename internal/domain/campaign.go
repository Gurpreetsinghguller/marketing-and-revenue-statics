package domain

// from the requirements, we need a Campaign entity with fields like ID, Name, Description, Status, DateRange, Channel, CreatedBy, IsPublic
type Campaign struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // Active, Paused, Completed
	DateRange   DateRange `json:"date_range"`
	Budget      float64   `json:"budget"`
	Channel     string    `json:"channel"` // Email, Social Media, etc.
	CreatedBy   string    `json:"created_by"`
	IsPublic    bool      `json:"is_public"`
}

type DateRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// This will be implemented by persistence layer (e.g., database) to manage campaign data
type CampaignRepo interface {
	// Create saves a new campaign
	Create(campaign *Campaign) error

	// GetByID retrieves a campaign by ID
	GetByID(id string) (*Campaign, error)

	// GetAll retrieves all campaigns
	GetAll() ([]Campaign, error)

	// GetByStatus retrieves campaigns filtered by status (Active, Paused, Completed)
	GetByStatus(status string) ([]Campaign, error)

	// GetByCreatedBy retrieves campaigns created by a specific marketer
	GetByCreatedBy(userID string) ([]Campaign, error)

	// GetByChannel retrieves campaigns filtered by channel
	GetByChannel(channel string) ([]Campaign, error)

	// GetPublic retrieves all public campaigns (accessible without auth)
	GetPublic() ([]Campaign, error)

	// GetByDateRange retrieves campaigns within a date range
	GetByDateRange(startDate, endDate string) ([]Campaign, error)

	// Update updates an existing campaign
	Update(campaign *Campaign) error

	// Delete removes a campaign
	Delete(id string) error

	// Search searches campaigns by name or description
	Search(query string) ([]Campaign, error)

	// GetWithFilters retrieves campaigns with multiple filters applied
	// Filters: status, channel, created_by, date_range, etc.
	GetWithFilters(filters map[string]interface{}) ([]Campaign, error)
}
