package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	mocks "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/mocks"
	mockpkg "github.com/stretchr/testify/mock"
)

func TestGetUserEngagementHandler_success(t *testing.T) {
	mock := mocks.NewEngagementUseCaseInterface(t)
	mock.On("GetUserEngagement", mockpkg.Anything, "u1").
		Return(&domain.UserEngagement{UserID: "u1", TotalInteractions: 5}, nil)

	h := NewEngagementHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/engagement/users/u1", nil)
	req = mux.SetURLVars(req, map[string]string{"user_id": "u1"})
	req.Header.Set("X-User-ID", "u1")
	rr := httptest.NewRecorder()

	h.GetUserEngagementHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rr.Code)
	}
}

func TestGetCampaignFunnelHandler_success(t *testing.T) {
	mock := mocks.NewEngagementUseCaseInterface(t)
	mock.On("GetCampaignFunnel", mockpkg.Anything, "c1").
		Return(&domain.CampaignFunnel{CampaignID: "c1", Impressions: 100, Clicks: 10}, nil)

	h := NewEngagementHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/engagement/campaigns/c1/funnel", nil)
	req = mux.SetURLVars(req, map[string]string{"campaign_id": "c1"})
	req.Header.Set("X-User-ID", "u1")
	rr := httptest.NewRecorder()

	h.GetCampaignFunnelHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rr.Code)
	}
}
