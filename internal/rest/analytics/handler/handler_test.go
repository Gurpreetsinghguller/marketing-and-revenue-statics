package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/mocks"
	mockpkg "github.com/stretchr/testify/mock"
)

func TestGetAnalyticsReportHandler_success(t *testing.T) {
	mock := mocks.NewAnalyticsUseCaseInterface(t)
	mock.On("GetAnalyticsReport", mockpkg.Anything).
		Return(map[string]interface{}{"total_campaigns": 2}, nil)

	h := NewAnalyticsHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/analytics", nil)
	req.Header.Set("X-User-ID", "u1")
	rr := httptest.NewRecorder()

	h.GetAnalyticsReportHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rr.Code)
	}
}

func TestGetDailyReportHandler_success(t *testing.T) {
	mock := mocks.NewAnalyticsUseCaseInterface(t)
	mock.On("GetDailyReport", mockpkg.Anything, "2026-03-01").
		Return(map[string]interface{}{"date": "2026-03-01", "total_events": 10}, nil)

	h := NewAnalyticsHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/analytics/daily?date=2026-03-01", nil)
	req.Header.Set("X-User-ID", "u1")
	rr := httptest.NewRecorder()

	h.GetDailyReportHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rr.Code)
	}
}

func TestGetPublicStatsHandler_success(t *testing.T) {
	mock := mocks.NewAnalyticsUseCaseInterface(t)
	mock.On("GetPublicStats", mockpkg.Anything).
		Return(map[string]interface{}{"total_campaigns": 1}, nil)

	h := NewAnalyticsHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/analytics/public", nil)
	rr := httptest.NewRecorder()

	h.GetPublicStatsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rr.Code)
	}
}
