package usecase

import (
	"context"
	"testing"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	mocks "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/mocks"
)

func TestEngagementUseCase_GetUserEngagement_Success(t *testing.T) {
	repo := mocks.NewEventRepo(t)
	events := []domain.Event{{UserID: "u1", CampaignID: "c1"}, {UserID: "u1", CampaignID: "c2"}}
	repo.On("GetByUserID", "u1").Return(events, nil)

	uc := NewEngagementUseCase(repo)
	engagement, err := uc.GetUserEngagement(context.Background(), "u1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if engagement == nil || engagement.TotalInteractions != 2 {
		t.Errorf("unexpected engagement: %+v", engagement)
	}
}

func TestEngagementUseCase_GetCampaignFunnel_Success(t *testing.T) {
	repo := mocks.NewEventRepo(t)
	events := []domain.Event{
		{EventType: domain.EventType("impressions")},
		{EventType: domain.EventType("clicks")},
		{EventType: domain.EventType("conversions")},
	}
	repo.On("GetByCampaignID", "c1").Return(events, nil)

	uc := NewEngagementUseCase(repo)
	funnel, err := uc.GetCampaignFunnel(context.Background(), "c1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if funnel == nil || funnel.Impressions != 1 {
		t.Errorf("unexpected funnel: %+v", funnel)
	}
}
