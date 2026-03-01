package campaign

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/db"
)

const campaignPrefix = "campaigns"

type campaignDBDateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type campaignDBModel struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Status      string              `json:"status"`
	DateRange   campaignDBDateRange `json:"date_range"`
	Budget      float64             `json:"budget"`
	Channel     string              `json:"channel"`
	CreatedBy   string              `json:"created_by"`
	IsPublic    bool                `json:"is_public"`
}

func toCampaignDBModel(c domain.Campaign) campaignDBModel {
	cc := &campaignDBModel{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Status:      string(c.Status),

		Budget:    c.Budget,
		Channel:   c.Channel,
		CreatedBy: c.CreatedBy,
		IsPublic:  c.IsPublic,
	}
	if c.DateRange.Start != nil {
		cc.DateRange.Start = *c.DateRange.Start
	}
	if c.DateRange.End != nil {
		cc.DateRange.End = *c.DateRange.End
	}
	return *cc
}

func (m campaignDBModel) toDomain() domain.Campaign {
	return domain.Campaign{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Status:      domain.CampaignStatus(m.Status),
		DateRange: domain.DateRange{
			Start: &m.DateRange.Start,
			End:   &m.DateRange.End,
		},
		Budget:    m.Budget,
		Channel:   m.Channel,
		CreatedBy: m.CreatedBy,
		IsPublic:  m.IsPublic,
	}
}

func decodeCampaignDBModel(value interface{}) (campaignDBModel, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return campaignDBModel{}, err
	}
	var model campaignDBModel
	if err := json.Unmarshal(b, &model); err != nil {
		return campaignDBModel{}, err
	}
	return model, nil
}

func campaignKey(id string) string {
	return campaignPrefix + "/" + id
}

// CampaignRepository implements domain.CampaignRepo using JSON file storage
type CampaignRepository struct {
	storage db.PersistenceDB
}

// NewCampaignRepository creates a new campaign repository
func NewCampaignRepository(storage ...db.PersistenceDB) *CampaignRepository {
	selected := db.PersistenceDB(db.NewStorageMgr())
	if len(storage) > 0 && storage[0] != nil {
		selected = storage[0]
	}
	return &CampaignRepository{storage: selected}
}

// Create saves a new campaign
func (r *CampaignRepository) Create(campaign *domain.Campaign) error {
	// Generate ID if not provided
	if campaign.ID == "" {
		campaign.ID = db.GenerateID("camp")
	}

	return r.storage.Create(context.Background(), campaignKey(campaign.ID), toCampaignDBModel(*campaign))
}

// GetByID retrieves a campaign by ID
func (r *CampaignRepository) GetByID(id string) (*domain.Campaign, error) {
	stored, err := r.storage.Read(context.Background(), campaignKey(id))
	if err != nil {
		return nil, err
	}

	model, err := decodeCampaignDBModel(stored)
	if err != nil {
		return nil, err
	}

	entity := model.toDomain()
	return &entity, nil
}

// GetAll retrieves all campaigns
func (r *CampaignRepository) GetAll() ([]domain.Campaign, error) {
	stored, err := r.storage.List(context.Background(), campaignPrefix)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Campaign, 0, len(stored))
	for _, item := range stored {
		model, err := decodeCampaignDBModel(item)
		if err != nil {
			return nil, err
		}
		result = append(result, model.toDomain())
	}

	return result, nil
}

// GetPublic retrieves all public campaigns
func (r *CampaignRepository) GetPublic() ([]domain.Campaign, error) {
	all, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	var result []domain.Campaign
	for i := range all {
		if all[i].IsPublic {
			result = append(result, all[i])
		}
	}
	return result, nil
}

// Update updates a campaign
func (r *CampaignRepository) Update(campaign *domain.Campaign) error {
	if _, err := r.storage.Read(context.Background(), campaignKey(campaign.ID)); err != nil {
		return err
	}

	return r.storage.Update(context.Background(), campaignKey(campaign.ID), toCampaignDBModel(*campaign))
}

// Delete removes a campaign
func (r *CampaignRepository) Delete(id string) error {
	return r.storage.Delete(context.Background(), campaignKey(id))
}

// Search searches campaigns by name or description
func (r *CampaignRepository) Search(query string) ([]domain.Campaign, error) {
	all, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	pattern, err := regexp.Compile("(?i)" + query)
	if err != nil {
		return nil, fmt.Errorf("invalid search pattern: %w", err)
	}

	var result []domain.Campaign
	for i := range all {
		if pattern.MatchString(all[i].Name) || pattern.MatchString(all[i].Description) {
			result = append(result, all[i])
		}
	}
	return result, nil
}

// GetWithFilters retrieves campaigns with multiple filters
func (r *CampaignRepository) GetWithFilters(filters map[string]interface{}) ([]domain.Campaign, error) {
	all, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	var result []domain.Campaign

	for i := range all {
		campaign := all[i]
		match := true

		// Apply filters
		if status, ok := filters["status"].(string); ok {
			if !strings.EqualFold(string(campaign.Status), status) {
				match = false
			}
		}
		if channel, ok := filters["channel"].(string); ok && campaign.Channel != channel {
			match = false
		}
		if createdBy, ok := filters["created_by"].(string); ok && campaign.CreatedBy != createdBy {
			match = false
		}
		if budget, ok := filters["budget"].(float64); ok {
			if math.Abs(campaign.Budget-budget) > 1e-9 {
				match = false
			}
		}
		if isPublic, ok := filters["is_public"].(bool); ok {
			if campaign.IsPublic != isPublic {
				match = false
			}
		}

		// Add more filter logic as needed

		if match {
			result = append(result, campaign)
		}
	}
	return result, nil
}
