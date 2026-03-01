package usecase

import (
	"context"
	"testing"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	mocks "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/mocks"
	tmock "github.com/stretchr/testify/mock"
)

func TestEventUseCase_TrackEvent_Success(t *testing.T) {
	repo := mocks.NewEventRepo(t)
	repo.On("Create", tmock.Anything).Return(nil)

	uc := NewEventUseCase(repo)
	event := &domain.Event{CampaignID: "c1", EventType: domain.EventType("clicks"), UserID: "u1"}
	result, err := uc.TrackEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil || result.CampaignID != "c1" {
		t.Errorf("unexpected event: %+v", result)
	}
}

func TestEventUseCase_GetEventsByUser_Success(t *testing.T) {
	repo := mocks.NewEventRepo(t)
	events := []domain.Event{{ID: "e1", UserID: "u1"}}
	repo.On("GetByUserID", "u1").Return(events, nil)

	uc := NewEventUseCase(repo)
	results, err := uc.GetEventsByUser(context.Background(), "u1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 event, got %d", len(results))
	}
}
