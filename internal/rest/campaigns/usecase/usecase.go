package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/constant"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
)

// CampaignUseCase handles campaign business logic
type CampaignUseCase struct {
	campaignRepo domain.CampaignRepo
}

type CampaignUseCaseInterface interface {
	CreateCampaign(ctx context.Context, campaign *domain.Campaign, userID string) (*domain.Campaign, error)
	GetCampaignsWithFilters(ctx context.Context, userID, userRole, status string, budget *float64, isPublic *bool) ([]domain.Campaign, error)
	GetCampaignByID(ctx context.Context, campaignID string) (*domain.Campaign, error)
	GetAllCampaigns(ctx context.Context) ([]domain.Campaign, error)
	GetPublicCampaigns(ctx context.Context) ([]domain.Campaign, error)
	UpdateCampaign(ctx context.Context, campaignID string, updatedCampaign *domain.Campaign, userID, userRole string) (*domain.Campaign, error)
	DeleteCampaignWithRole(ctx context.Context, campaignID, userID, userRole string) error
	SearchCampaigns(ctx context.Context, query string) ([]domain.Campaign, error)
	PatchCampaignStatus(ctx context.Context, campaignID, userID, userRole, status string) (*domain.Campaign, error)
	EndCampaign(ctx context.Context, campaignID, userID, userRole string) (*domain.Campaign, error)
}

// NewCampaignUseCase creates a new campaign usecase
func NewCampaignUseCase(campaignRepo domain.CampaignRepo) *CampaignUseCase {
	return &CampaignUseCase{
		campaignRepo: campaignRepo,
	}
}

// CreateCampaign creates a new campaign
func (c *CampaignUseCase) CreateCampaign(ctx context.Context, campaign *domain.Campaign, userID string) (*domain.Campaign, error) {
	if campaign == nil {
		return nil, errors.New("campaign is required")
	}
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	campaign.CreatedBy = userID

	if err := c.campaignRepo.Create(campaign); err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	return campaign, nil
}

// GetCampaignsWithFilters retrieves campaigns for authenticated users with optional filters.
func (c *CampaignUseCase) GetCampaignsWithFilters(ctx context.Context, userID, userRole, status string, budget *float64, isPublic *bool) ([]domain.Campaign, error) {
	_ = ctx

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user id is required")
	}
	if strings.TrimSpace(userRole) == "" {
		return nil, errors.New("user role is required")
	}

	filters := map[string]interface{}{
		"request_user_id":   strings.TrimSpace(userID),
		"request_user_role": strings.ToLower(strings.TrimSpace(userRole)),
		"created_by":        strings.TrimSpace(userID),
	}

	if status != "" {
		filters["status"] = strings.ToLower(strings.TrimSpace(status))
	}
	if budget != nil {
		filters["budget"] = *budget
	}
	if isPublic != nil {
		filters["is_public"] = *isPublic
	}

	campaigns, err := c.campaignRepo.GetWithFilters(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campaigns with filters: %w", err)
	}

	return campaigns, nil
}

// GetCampaignByID retrieves a campaign by ID
func (c *CampaignUseCase) GetCampaignByID(ctx context.Context, campaignID string) (*domain.Campaign, error) {
	_ = ctx

	campaign, err := c.campaignRepo.GetByID(campaignID)
	if err != nil || campaign == nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	return campaign, nil
}

// GetAllCampaigns retrieves all campaigns
func (c *CampaignUseCase) GetAllCampaigns(ctx context.Context) ([]domain.Campaign, error) {
	_ = ctx

	campaigns, err := c.campaignRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campaigns: %w", err)
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

// UpdateCampaign updates an existing campaign
func (c *CampaignUseCase) UpdateCampaign(ctx context.Context, campaignID string, updatedCampaign *domain.Campaign, userID, userRole string) (*domain.Campaign, error) {
	if updatedCampaign == nil {
		return nil, errors.New("campaign is required")
	}

	// Fetch existing campaign
	campaign, err := c.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if updatedCampaign.Name != "" {
		campaign.Name = updatedCampaign.Name
	}
	if updatedCampaign.Description != "" {
		campaign.Description = updatedCampaign.Description
	}
	if updatedCampaign.Status != "" {
		campaign.Status = updatedCampaign.Status
	}

	campaign.Name = updatedCampaign.Name
	campaign.Description = updatedCampaign.Description
	campaign.Status = updatedCampaign.Status
	campaign.Budget = updatedCampaign.Budget
	campaign.IsPublic = updatedCampaign.IsPublic

	// Save updated campaign
	if err := c.campaignRepo.Update(campaign); err != nil {
		return nil, fmt.Errorf("failed to update campaign: %w", err)
	}

	return campaign, nil
}

// DeleteCampaignWithRole removes a campaign with admin override
func (c *CampaignUseCase) DeleteCampaignWithRole(ctx context.Context, campaignID, userID, userRole string) error {
	// Fetch existing campaign
	campaign, err := c.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return err
	}
	if campaign.CreatedBy != userID && strings.ToLower(userRole) != constant.RoleAdmin {
		return errors.New("user not authorized to delete this campaign")
	}

	if err := c.campaignRepo.Delete(campaignID); err != nil {
		return fmt.Errorf("failed to delete campaign: %w", err)
	}

	return nil
}

// SearchCampaigns searches campaigns by name or description
func (c *CampaignUseCase) SearchCampaigns(ctx context.Context, query string) ([]domain.Campaign, error) {
	_ = ctx

	return c.campaignRepo.Search(query)
}

// PatchCampaignStatus updates the status of a campaign.
func (c *CampaignUseCase) PatchCampaignStatus(ctx context.Context, campaignID, userID, userRole, status string) (*domain.Campaign, error) {
	campaign, err := c.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	campaign.Status = domain.CampaignStatus(status)
	if err := c.campaignRepo.Update(campaign); err != nil {
		return nil, fmt.Errorf("failed to update campaign status: %w", err)
	}

	return campaign, nil
}

// EndCampaign marks an active/paused campaign as inactive and sets end date to now.
func (c *CampaignUseCase) EndCampaign(ctx context.Context, campaignID, userID, userRole string) (*domain.Campaign, error) {
	if strings.TrimSpace(campaignID) == "" {
		return nil, errors.New("campaign id is required")
	}
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user id is required")
	}

	campaign, err := c.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	role := strings.ToLower(strings.TrimSpace(userRole))
	switch role {
	case "admin":
		// admin can end any campaign
	case "marketer":
		if campaign.CreatedBy != userID {
			return nil, errors.New("user not authorized to end this campaign")
		}
	default:
		return nil, errors.New("user not authorized to end this campaign")
	}

	if campaign.Status != domain.CampaignStatusActive && campaign.Status != domain.CampaignStatusPaused {
		return nil, errors.New("only active or paused campaigns can be ended")
	}

	now := time.Now().UTC()
	campaign.DateRange.End = &now
	campaign.Status = domain.CampaignStatusInactive
	campaign.UpdatedAt = now

	if err := c.campaignRepo.Update(campaign); err != nil {
		return nil, fmt.Errorf("failed to end campaign: %w", err)
	}

	return campaign, nil
}
