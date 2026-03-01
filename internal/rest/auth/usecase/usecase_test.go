package usecase_test

import (
	"context"
	"os"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	mocks "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/mocks"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/auth/usecase"
	tmock "github.com/stretchr/testify/mock"
)

func TestAuthUseCase_Register_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("JWT_SECRET")

	repo := mocks.NewUserRepo(t)

	// Expect EmailExists -> false
	repo.On("EmailExists", "a@b.com").Return(false)

	// When Create is called, populate ID and timestamps on the passed user
	repo.On("Create", tmock.AnythingOfType("*domain.User")).Run(func(args tmock.Arguments) {
		if u, ok := args.Get(0).(*domain.User); ok {
			u.ID = "user_test"
			u.CreatedAt = time.Now().UTC()
		}
	}).Return(nil)

	uc := usecase.NewAuthUseCase(repo)

	user := &domain.User{
		Email:    "a@b.com",
		Password: "secret",
		Name:     "Tester",
		Role:     domain.MarketerRole,
	}

	res, err := uc.Register(context.Background(), user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res == nil || res.UserID == "" || res.Token == "" {
		t.Fatalf("unexpected register response: %+v", res)
	}
}

func TestAuthUseCase_Register_EmailExists(t *testing.T) {
	repo := mocks.NewUserRepo(t)
	repo.On("EmailExists", "a@b.com").Return(true)

	uc := usecase.NewAuthUseCase(repo)

	user := &domain.User{Email: "a@b.com", Password: "secret", Name: "X"}
	_, err := uc.Register(context.Background(), user)
	if err == nil {
		t.Fatal("expected error when email exists")
	}
}

func TestAuthUseCase_Login_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("JWT_SECRET")

	// prepare hashed password that matches "secret"
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

	repo := mocks.NewUserRepo(t)
	repo.On("GetByEmail", "a@b.com").Return(&domain.User{ID: "u1", Email: "a@b.com", Password: string(hash), Role: domain.MarketerRole}, nil)

	uc := usecase.NewAuthUseCase(repo)
	cred := &domain.User{Email: "a@b.com", Password: "secret"}
	res, err := uc.Login(context.Background(), cred)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res == nil || res.Token == "" || res.User == nil || res.User.ID == "" {
		t.Fatalf("unexpected login response: %+v", res)
	}
}

func TestAuthUseCase_Login_InvalidCredentials(t *testing.T) {
	// return a user with a password that doesn't match
	hash, _ := bcrypt.GenerateFromPassword([]byte("other"), bcrypt.DefaultCost)

	repo := mocks.NewUserRepo(t)
	repo.On("GetByEmail", "a@b.com").Return(&domain.User{ID: "u1", Email: "a@b.com", Password: string(hash)}, nil)

	uc := usecase.NewAuthUseCase(repo)
	cred := &domain.User{Email: "a@b.com", Password: "secret"}
	_, err := uc.Login(context.Background(), cred)
	if err == nil {
		t.Fatal("expected error for invalid credentials")
	}
}
