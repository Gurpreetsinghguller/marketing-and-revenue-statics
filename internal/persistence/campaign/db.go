package campaign

import (
	"strings"
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/db"
)

// CampaignRepository implements domain.CampaignRepo using JSON file storage
type CampaignRepository struct {
	storage   *db.StorageMgr
	campaigns []domain.Campaign
}

// NewCampaignRepository creates a new campaign repository
func NewCampaignRepository() *CampaignRepository {
	repo := &CampaignRepository{
		storage:   db.NewStorageMgr(),
		campaigns: []domain.Campaign{},
	}
	// Load existing campaigns from file
	repo.storage.ReadJSON(db.CampaignsFile, &repo.campaigns)
	return repo
}

// Create saves a new campaign
func (r *CampaignRepository) Create(campaign *domain.Campaign) error {
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()

	// Generate ID if not provided
	if campaign.ID == "" {
		campaign.ID = "camp_" + generateRandomString(12)
	}

	r.campaigns = append(r.campaigns, *campaign)
	return r.storage.WriteJSON(CampaignsFile, r.campaigns)
}

// GetByID retrieves a campaign by ID
func (r *CampaignRepository) GetByID(id string) (*domain.Campaign, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	for i := range r.campaigns {
		if r.campaigns[i].ID == id {
			return &r.campaigns[i], nil
		}
	}
	return nil, ErrNotFound
}

// GetAll retrieves all campaigns
func (r *CampaignRepository) GetAll() ([]domain.Campaign, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	return r.campaigns, nil
}

// GetByStatus retrieves campaigns by status
func (r *CampaignRepository) GetByStatus(status string) ([]domain.Campaign, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	var result []domain.Campaign
	for i := range r.campaigns {
		if r.campaigns[i].Status == status {
			result = append(result, r.campaigns[i])
		}
	}
	return result, nil
}

// GetByCreatedBy retrieves campaigns created by a user
func (r *CampaignRepository) GetByCreatedBy(userID string) ([]domain.Campaign, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	var result []domain.Campaign
	for i := range r.campaigns {
		if r.campaigns[i].CreatedBy == userID {
			result = append(result, r.campaigns[i])
		}
	}
	return result, nil
}

// GetByChannel retrieves campaigns by channel
func (r *CampaignRepository) GetByChannel(channel string) ([]domain.Campaign, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	var result []domain.Campaign
	for i := range r.campaigns {
		if r.campaigns[i].Channel == channel {
			result = append(result, r.campaigns[i])
		}
	}
	return result, nil
}

// GetPublic retrieves all public campaigns
func (r *CampaignRepository) GetPublic() ([]domain.Campaign, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	var result []domain.Campaign
	for i := range r.campaigns {
		if r.campaigns[i].IsPublic {
			result = append(result, r.campaigns[i])
		}
	}
	return result, nil
}

// GetByDateRange retrieves campaigns within a date range
func (r *CampaignRepository) GetByDateRange(startDate, endDate string) ([]domain.Campaign, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	var result []domain.Campaign
	for i := range r.campaigns {
		campaignStart, _ := time.Parse("2006-01-02", r.campaigns[i].DateRange.Start)
		campaignEnd, _ := time.Parse("2006-01-02", r.campaigns[i].DateRange.End)

		// Check overlap
		if campaignStart.Before(end) && campaignEnd.After(start) {
			result = append(result, r.campaigns[i])
		}
	}
	return result, nil
}

// Update updates a campaign
func (r *CampaignRepository) Update(campaign *domain.Campaign) error {
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()

	for i := range r.campaigns {
		if r.campaigns[i].ID == campaign.ID {
			r.campaigns[i] = *campaign
			return r.storage.WriteJSON(CampaignsFile, r.campaigns)
		}
	}
	return ErrNotFound
}

// Delete removes a campaign
func (r *CampaignRepository) Delete(id string) error {
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()

	for i := range r.campaigns {
		if r.campaigns[i].ID == id {
			r.campaigns = append(r.campaigns[:i], r.campaigns[i+1:]...)
			return r.storage.WriteJSON(CampaignsFile, r.campaigns)
		}
	}
	return ErrNotFound
}

// Search searches campaigns by name or description
func (r *CampaignRepository) Search(query string) ([]domain.Campaign, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	query = strings.ToLower(query)
	var result []domain.Campaign
	for i := range r.campaigns {
		if strings.Contains(strings.ToLower(r.campaigns[i].Name), query) ||
			strings.Contains(strings.ToLower(r.campaigns[i].Description), query) {
			result = append(result, r.campaigns[i])
		}
	}
	return result, nil
}

// GetWithFilters retrieves campaigns with multiple filters
func (r *CampaignRepository) GetWithFilters(filters map[string]interface{}) ([]domain.Campaign, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	var result []domain.Campaign

	for i := range r.campaigns {
		campaign := r.campaigns[i]
		match := true

		// Apply filters
		if status, ok := filters["status"].(string); ok && campaign.Status != status {
			match = false
		}
		if channel, ok := filters["channel"].(string); ok && campaign.Channel != channel {
			match = false
		}
		if createdBy, ok := filters["created_by"].(string); ok && campaign.CreatedBy != createdBy {
			match = false
		}

		// Add more filter logic as needed

		if match {
			result = append(result, campaign)
		}
	}
	return result, nil
}
