package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	mocks "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/mocks"
	mockpkg "github.com/stretchr/testify/mock"
)

func TestTrackEventHandler_success(t *testing.T) {
	mock := mocks.NewEventUseCaseInterface(t)
	mock.On("TrackEvent", mockpkg.Anything, mockpkg.Anything).
		Return(&domain.Event{ID: "e1", CampaignID: "c1"}, nil)

	h := NewEventHandler(mock)

	reqBody := TrackEventRequest{
		CampaignID: "c1",
		EventType:  "clicks",
		UserID:     "u1",
		Timestamp:  "2026-03-01T12:00:00Z",
	}
	b, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.TrackEventHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d got %d", http.StatusCreated, rr.Code)
	}
}
