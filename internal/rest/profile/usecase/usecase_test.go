package usecase

import (
	"context"
	"testing"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	mocks "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/mocks"
	tmock "github.com/stretchr/testify/mock"
)

func TestProfileUseCase_GetProfile_Success(t *testing.T) {
	repo := mocks.NewUserRepo(t)
	repo.On("GetByID", "u1").Return(&domain.User{ID: "u1", Name: "John", Email: "john@example.com"}, nil)

	uc := NewProfileUseCase(repo)
	user, err := uc.GetProfile(context.Background(), "u1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil || user.ID != "u1" || user.Name != "John" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestProfileUseCase_UpdateProfile_Success(t *testing.T) {
	repo := mocks.NewUserRepo(t)
	repo.On("GetByID", "u1").Return(&domain.User{ID: "u1", Name: "John", Email: "john@example.com"}, nil)
	repo.On("Update", tmock.MatchedBy(func(u *domain.User) bool {
		return u.ID == "u1" && u.Name == "Jane"
	})).Return(nil)

	uc := NewProfileUseCase(repo)
	updates := &domain.User{Name: "Jane"}
	user, err := uc.UpdateProfile(context.Background(), "u1", updates)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil || user.Name != "Jane" {
		t.Errorf("unexpected user after update: %+v", user)
	}
}
