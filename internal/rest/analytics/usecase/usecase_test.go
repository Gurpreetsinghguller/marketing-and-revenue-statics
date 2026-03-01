package usecase

import (
	"context"
	"testing"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	mocks "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/mocks"
)

func TestAnalyticsUseCase_GetPublicStats_Success(t *testing.T) {
	campaignRepo := mocks.NewCampaignRepo(t)
	eventRepo := mocks.NewEventRepo(t)

	campaigns := []domain.Campaign{{ID: "c1", Name: "Campaign 1", IsPublic: true}}
	campaignRepo.On("GetPublic").Return(campaigns, nil)

	uc := NewAnalyticsUseCase(campaignRepo, eventRepo)
	stats, err := uc.GetPublicStats(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if stats == nil || stats["total_campaigns"] != 1 {
		t.Errorf("unexpected stats: %+v", stats)
	}
}

func TestAnalyticsUseCase_GetCampaignStats_Success(t *testing.T) {
	campaignRepo := mocks.NewCampaignRepo(t)
	eventRepo := mocks.NewEventRepo(t)

	campaign := &domain.Campaign{ID: "c1", Name: "Campaign 1"}
	campaignRepo.On("GetByID", "c1").Return(campaign, nil)

	events := []domain.Event{
		{EventType: domain.EventType("impressions")},
		{EventType: domain.EventType("clicks")},
		{EventType: domain.EventType("conversions"), Metadata: domain.Metadata{Amount: 100}},
	}
	eventRepo.On("GetByCampaignID", "c1").Return(events, nil)

	uc := NewAnalyticsUseCase(campaignRepo, eventRepo)
	stats, err := uc.GetCampaignStats(context.Background(), "c1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if stats == nil || stats.TotalImpressions != 1 {
		t.Errorf("unexpected stats: %+v", stats)
	}
}

func TestAnalyticsUseCase_GetDailyReport_Success(t *testing.T) {
	campaignRepo := mocks.NewCampaignRepo(t)
	eventRepo := mocks.NewEventRepo(t)

	eventRepo.On("GetAll").Return([]domain.Event{}, nil)

	uc := NewAnalyticsUseCase(campaignRepo, eventRepo)
	report, err := uc.GetDailyReport(context.Background(), "2026-03-01")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if report == nil || report["date"] == nil {
		t.Errorf("unexpected report: %+v", report)
	}
}

func TestAnalyticsUseCase_GetWeeklyReport_Success(t *testing.T) {
	campaignRepo := mocks.NewCampaignRepo(t)
	eventRepo := mocks.NewEventRepo(t)

	eventRepo.On("GetAll").Return([]domain.Event{}, nil)

	uc := NewAnalyticsUseCase(campaignRepo, eventRepo)
	report, err := uc.GetWeeklyReport(context.Background(), "2026-03-01")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if report == nil {
		t.Error("expected report")
	}
}

func TestAnalyticsUseCase_GetMonthlyReport_Success(t *testing.T) {
	campaignRepo := mocks.NewCampaignRepo(t)
	eventRepo := mocks.NewEventRepo(t)

	eventRepo.On("GetAll").Return([]domain.Event{}, nil)

	uc := NewAnalyticsUseCase(campaignRepo, eventRepo)
	report, err := uc.GetMonthlyReport(context.Background(), "2026-03")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if report == nil {
		t.Error("expected report")
	}
}
