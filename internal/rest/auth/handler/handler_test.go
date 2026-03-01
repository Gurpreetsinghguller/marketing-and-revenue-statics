package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	mocks "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/mocks"
	auth_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/auth/usecase"
	mockpkg "github.com/stretchr/testify/mock"
)

func TestRegisterHandler_success(t *testing.T) {
	mock := mocks.NewAuthUseCaseInterface(t)
	mock.On("Register", mockpkg.Anything, mockpkg.Anything).Return(&auth_usecase.RegisterResponse{UserID: "user_123", Token: "tok"}, nil)
	h := NewAuthHandler(mock)

	reqBody := map[string]string{
		"email":    "foo@example.com",
		"password": "secret",
		"name":     "Foo Bar",
		"role":     "marketer",
	}
	b, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.RegisterHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d got %d", http.StatusCreated, rr.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if resp["user_id"] != "user_123" {
		t.Errorf("unexpected user_id: %v", resp["user_id"])
	}
	if resp["token"] != "tok" {
		t.Errorf("unexpected token: %v", resp["token"])
	}
}

func TestRegisterHandler_failure(t *testing.T) {
	mock := mocks.NewAuthUseCaseInterface(t)
	mock.On("Register", mockpkg.Anything, mockpkg.Anything).Return(nil, errors.New("email already exists"))
	h := NewAuthHandler(mock)

	reqBody := map[string]string{
		"email":    "foo@example.com",
		"password": "secret",
		"name":     "Foo Bar",
		"role":     "marketer",
	}
	b, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.RegisterHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestLoginHandler_success(t *testing.T) {
	mock := mocks.NewAuthUseCaseInterface(t)
	mock.On("Login", mockpkg.Anything, mockpkg.Anything).Return(&auth_usecase.LoginResponse{Token: "tok", User: &domain.User{ID: "user_123"}}, nil)
	h := NewAuthHandler(mock)

	reqBody := map[string]string{
		"email":    "foo@example.com",
		"password": "secret",
	}
	b, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.LoginHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rr.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if resp["token"] != "tok" {
		t.Errorf("unexpected token: %v", resp["token"])
	}
}

func TestLoginHandler_failure(t *testing.T) {
	mock := mocks.NewAuthUseCaseInterface(t)
	mock.On("Login", mockpkg.Anything, mockpkg.Anything).Return(nil, errors.New("invalid email or password"))
	h := NewAuthHandler(mock)

	reqBody := map[string]string{
		"email":    "foo@example.com",
		"password": "wrong",
	}
	b, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.LoginHandler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d got %d", http.StatusUnauthorized, rr.Code)
	}
}
