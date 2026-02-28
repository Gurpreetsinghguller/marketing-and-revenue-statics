package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/db"
)

// CampaignUseCase handles campaign business logic
type CampaignUseCase struct {
	db db.PersistenceDB
}

// NewCampaignUseCase creates a new campaign usecase
func NewCampaignUseCase(db db.PersistenceDB) *CampaignUseCase {
	return &CampaignUseCase{
		db: db,
	}
}

// CreateCampaignRequest represents campaign creation input
type CreateCampaignRequest struct {
	Name        string
	Description string
	Status      string
	DateRange   domain.DateRange
	Budget      float64
	Channel     string
	CreatedBy   string
	IsPublic    bool
}

// GetCampaignsRequest represents filtering options
type GetCampaignsRequest struct {
	Status    string
	Channel   string
	CreatedBy string
	IsPublic  *bool
}

// CreateCampaign creates a new campaign
func (c *CampaignUseCase) CreateCampaign(ctx context.Context, req CreateCampaignRequest, userID string) (*domain.Campaign, error) {
	if req.Name == "" {
		return nil, errors.New("campaign name is required")
	}

	campaign := &domain.Campaign{
		ID:          fmt.Sprintf("campaign_%s_%d", userID, len([]int{})), // Use UUID in production
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		DateRange:   req.DateRange,
		Budget:      req.Budget,
		Channel:     req.Channel,
		CreatedBy:   userID,
		IsPublic:    req.IsPublic,
	}

	// Save campaign
	key := fmt.Sprintf("campaign:%s", campaign.ID)
	if err := c.db.Create(ctx, key, campaign); err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	// Index by created user for quick lookup
	userCampaignKey := fmt.Sprintf("user:%s:campaigns:%s", userID, campaign.ID)
	if err := c.db.Create(ctx, userCampaignKey, campaign.ID); err != nil {
		return nil, fmt.Errorf("failed to index campaign: %w", err)
	}

	return campaign, nil
}

// GetCampaignByID retrieves a campaign by ID
func (c *CampaignUseCase) GetCampaignByID(ctx context.Context, campaignID string) (*domain.Campaign, error) {
	key := fmt.Sprintf("campaign:%s", campaignID)
	campaignInterface, err := c.db.Read(ctx, key)
	if err != nil || campaignInterface == nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	campaign, ok := campaignInterface.(*domain.Campaign)
	if !ok {
		return nil, errors.New("invalid campaign data format")
	}

	return campaign, nil
}

// GetAllCampaigns retrieves all campaigns
func (c *CampaignUseCase) GetAllCampaigns(ctx context.Context) ([]domain.Campaign, error) {
	campaigns := []domain.Campaign{}
	results, err := c.db.List(ctx, "campaign:")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campaigns: %w", err)
	}

	for _, result := range results {
		if campaign, ok := result.(*domain.Campaign); ok {
			campaigns = append(campaigns, *campaign)
		}
	}

	return campaigns, nil
}

// GetPublicCampaigns retrieves all public campaigns
func (c *CampaignUseCase) GetPublicCampaigns(ctx context.Context) ([]domain.Campaign, error) {
	allCampaigns, err := c.GetAllCampaigns(ctx)
	if err != nil {
		return nil, err
	}

	var publicCampaigns []domain.Campaign
	for _, campaign := range allCampaigns {
		if campaign.IsPublic {
			publicCampaigns = append(publicCampaigns, campaign)
		}
	}

	return publicCampaigns, nil
}

// GetCampaignsByUser retrieves campaigns created by a specific user
func (c *CampaignUseCase) GetCampaignsByUser(ctx context.Context, userID string) ([]domain.Campaign, error) {
	campaigns := []domain.Campaign{}
	results, err := c.db.List(ctx, fmt.Sprintf("user:%s:campaigns:", userID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user campaigns: %w", err)
	}

	for _, result := range results {
		if campaignID, ok := result.(string); ok {
			campaign, err := c.GetCampaignByID(ctx, campaignID)
			if err == nil {
				campaigns = append(campaigns, *campaign)
			}
		}
	}

	return campaigns, nil
}

// GetCampaignsByStatus retrieves campaigns with a specific status
func (c *CampaignUseCase) GetCampaignsByStatus(ctx context.Context, status string) ([]domain.Campaign, error) {
	allCampaigns, err := c.GetAllCampaigns(ctx)
	if err != nil {
		return nil, err
	}

	var filteredCampaigns []domain.Campaign
	for _, campaign := range allCampaigns {
		if campaign.Status == status {
			filteredCampaigns = append(filteredCampaigns, campaign)
		}
	}

	return filteredCampaigns, nil
}

// GetCampaignsByChannel retrieves campaigns for a specific channel
func (c *CampaignUseCase) GetCampaignsByChannel(ctx context.Context, channel string) ([]domain.Campaign, error) {
	allCampaigns, err := c.GetAllCampaigns(ctx)
	if err != nil {
		return nil, err
	}

	var filteredCampaigns []domain.Campaign
	for _, campaign := range allCampaigns {
		if campaign.Channel == channel {
			filteredCampaigns = append(filteredCampaigns, campaign)
		}
	}

	return filteredCampaigns, nil
}

// UpdateCampaign updates an existing campaign
func (c *CampaignUseCase) UpdateCampaign(ctx context.Context, campaignID string, req CreateCampaignRequest, userID string) (*domain.Campaign, error) {
	// Fetch existing campaign
	campaign, err := c.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if campaign.CreatedBy != userID {
		return nil, errors.New("unauthorized: only creator can update campaign")
	}

	// Update fields
	campaign.Name = req.Name
	campaign.Description = req.Description
	campaign.Status = req.Status
	campaign.DateRange = req.DateRange
	campaign.Budget = req.Budget
	campaign.Channel = req.Channel
	campaign.IsPublic = req.IsPublic

	// Save updated campaign
	key := fmt.Sprintf("campaign:%s", campaignID)
	if err := c.db.Update(ctx, key, campaign); err != nil {
		return nil, fmt.Errorf("failed to update campaign: %w", err)
	}

	return campaign, nil
}

// DeleteCampaign removes a campaign
func (c *CampaignUseCase) DeleteCampaign(ctx context.Context, campaignID string, userID string) error {
	// Fetch existing campaign
	campaign, err := c.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return err
	}

	// Verify ownership
	if campaign.CreatedBy != userID {
		return errors.New("unauthorized: only creator can delete campaign")
	}

	// Delete campaign
	key := fmt.Sprintf("campaign:%s", campaignID)
	if err := c.db.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete campaign: %w", err)
	}

	// Remove user index
	userCampaignKey := fmt.Sprintf("user:%s:campaigns:%s", userID, campaignID)
	_ = c.db.Delete(ctx, userCampaignKey)

	return nil
}

// SearchCampaigns searches campaigns by name or description
func (c *CampaignUseCase) SearchCampaigns(ctx context.Context, query string) ([]domain.Campaign, error) {
	allCampaigns, err := c.GetAllCampaigns(ctx)
	if err != nil {
		return nil, err
	}

	var results []domain.Campaign
	for _, campaign := range allCampaigns {
		if contains(campaign.Name, query) || contains(campaign.Description, query) {
			results = append(results, campaign)
		}
	}

	return results, nil
}

// Helper function
func contains(str, substr string) bool {
	return len(substr) > 0 && len(str) >= len(substr) && str[:len(substr)] == substr
}
