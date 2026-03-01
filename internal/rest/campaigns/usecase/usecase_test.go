package usecase

import (
	"context"
	"testing"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	mocks "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/mocks"
	tmock "github.com/stretchr/testify/mock"
)

func TestCampaignUseCase_CreateCampaign_Success(t *testing.T) {
	repo := mocks.NewCampaignRepo(t)
	repo.On("Create", tmock.Anything).Return(nil)

	uc := NewCampaignUseCase(repo)

	campaign := &domain.Campaign{
		Name:        "Summer Sale",
		Description: "Summer promo",
		Status:      domain.CampaignStatusActive,
	}

	result, err := uc.CreateCampaign(context.Background(), campaign, "user1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil || result.CreatedBy != "user1" {
		t.Errorf("unexpected campaign: %+v", result)
	}
}

func TestCampaignUseCase_GetCampaignByID_Success(t *testing.T) {
	repo := mocks.NewCampaignRepo(t)
	campaign := &domain.Campaign{ID: "c1", Name: "Summer Sale"}
	repo.On("GetByID", "c1").Return(campaign, nil)

	uc := NewCampaignUseCase(repo)
	result, err := uc.GetCampaignByID(context.Background(), "c1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil || result.ID != "c1" {
		t.Errorf("unexpected campaign: %+v", result)
	}
}

func TestCampaignUseCase_GetAllCampaigns_Success(t *testing.T) {
	repo := mocks.NewCampaignRepo(t)
	campaigns := []domain.Campaign{
		{ID: "c1", Name: "Campaign 1"},
		{ID: "c2", Name: "Campaign 2"},
	}
	repo.On("GetAll").Return(campaigns, nil)

	uc := NewCampaignUseCase(repo)
	results, err := uc.GetAllCampaigns(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 campaigns, got %d", len(results))
	}
}

func TestCampaignUseCase_SearchCampaigns_Success(t *testing.T) {
	repo := mocks.NewCampaignRepo(t)
	campaigns := []domain.Campaign{{ID: "c1", Name: "Summer Sale"}}
	repo.On("Search", "Summer").Return(campaigns, nil)

	uc := NewCampaignUseCase(repo)
	results, err := uc.SearchCampaigns(context.Background(), "Summer")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 1 || results[0].Name != "Summer Sale" {
		t.Errorf("unexpected search results: %+v", results)
	}
}

func TestCampaignUseCase_UpdateCampaign_Success(t *testing.T) {
	repo := mocks.NewCampaignRepo(t)
	repo.On("GetByID", "c1").Return(&domain.Campaign{ID: "c1", Name: "Old Name"}, nil)
	repo.On("Update", tmock.Anything).Return(nil)

	uc := NewCampaignUseCase(repo)
	updates := &domain.Campaign{Name: "New Name"}
	result, err := uc.UpdateCampaign(context.Background(), "c1", updates, "user1", "marketer")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Name != "New Name" {
		t.Errorf("expected updated name, got %s", result.Name)
	}
}

func TestCampaignUseCase_EndCampaign_Success(t *testing.T) {
	repo := mocks.NewCampaignRepo(t)
	repo.On("GetByID", "c1").Return(&domain.Campaign{ID: "c1", Status: domain.CampaignStatusActive}, nil)
	repo.On("Update", tmock.Anything).Return(nil)

	uc := NewCampaignUseCase(repo)
	result, err := uc.EndCampaign(context.Background(), "c1", "user1", "admin")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Status != domain.CampaignStatusInactive {
		t.Errorf("expected status to be inactive, got %v", result.Status)
	}
	if result.DateRange.End == nil || result.DateRange.End.IsZero() {
		t.Error("expected end date to be set")
	}
}
