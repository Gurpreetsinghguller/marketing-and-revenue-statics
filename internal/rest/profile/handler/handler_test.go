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

func TestGetProfileHandler_success(t *testing.T) {
	mock := mocks.NewProfileUseCaseInterface(t)
	mock.On("GetProfile", mockpkg.Anything, "u1").Return(&domain.User{ID: "u1", Name: "John", Email: "john@example.com"}, nil)

	h := NewProfileHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.Header.Set("X-User-ID", "u1")
	rr := httptest.NewRecorder()

	h.GetProfileHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rr.Code)
	}

	var resp UserProfileResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if resp.ID != "u1" || resp.Name != "John" {
		t.Errorf("unexpected profile: %+v", resp)
	}
}

func TestUpdateProfileHandler_success(t *testing.T) {
	mock := mocks.NewProfileUseCaseInterface(t)
	mock.On("UpdateProfile", mockpkg.Anything, "u1", mockpkg.Anything).Return(&domain.User{ID: "u1", Name: "Jane", Email: "jane@example.com"}, nil)

	h := NewProfileHandler(mock)

	reqBody := map[string]string{"name": "Jane"}
	b, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewReader(b))
	req.Header.Set("X-User-ID", "u1")
	rr := httptest.NewRecorder()

	h.UpdateProfileHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rr.Code)
	}

	var resp UserProfileResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if resp.Name != "Jane" {
		t.Errorf("unexpected updated profile: %+v", resp)
	}
}
