package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/tokenmachine"
	"golang.org/x/crypto/bcrypt"
)

// AuthUseCase handles authentication business logic
type AuthUseCase struct {
	userRepo     domain.UserRepo
	tokenMachine tokenmachine.TokenMachine
}

type AuthUseCaseInterface interface {
	Register(ctx context.Context, user *domain.User) (*RegisterResponse, error)
	Login(ctx context.Context, credentials *domain.User) (*LoginResponse, error)
}

func NewAuthUseCase(userRepo domain.UserRepo, tokenMachine tokenmachine.TokenMachine) AuthUseCaseInterface {
	return &AuthUseCase{
		userRepo:     userRepo,
		tokenMachine: tokenMachine,
	}
}

type RegisterResponse struct {
	UserID string
	Token  string
}

type LoginResponse struct {
	Token string
	User  *domain.User
}

func (u *AuthUseCase) Register(ctx context.Context, user *domain.User) (*RegisterResponse, error) {
	_ = ctx

	if user == nil {
		return nil, errors.New("user is required")
	}

	if user.Email == "" || user.Password == "" || user.Name == "" {
		return nil, errors.New("email, password, and name are required")
	}

	if u.userRepo.EmailExists(user.Email) {
		return nil, errors.New("email already exists")
	}

	userID, err := generateUserID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate user id: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := &domain.User{
		ID:       userID,
		Email:    user.Email,
		Password: string(hashedPassword),
		Name:     user.Name,
		Role:     user.Role,
	}

	if err := u.userRepo.Create(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := u.tokenMachine.GenerateToken(newUser.ID, newUser.Role.String())
	if err != nil {
		if rollbackErr := u.userRepo.Delete(newUser.ID); rollbackErr != nil {
			return nil, fmt.Errorf("failed to generate auth token: %w (rollback failed: %v)", err, rollbackErr)
		}
		return nil, fmt.Errorf("failed to generate auth token: %w", err)
	}

	return &RegisterResponse{
		UserID: newUser.ID,
		Token:  token,
	}, nil
}

// Login authenticates a user
func (u *AuthUseCase) Login(ctx context.Context, credentials *domain.User) (*LoginResponse, error) {
	_ = ctx

	if credentials == nil {
		return nil, errors.New("credentials are required")
	}

	if credentials.Email == "" || credentials.Password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := u.userRepo.GetByEmail(credentials.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := u.tokenMachine.GenerateToken(user.ID, user.Role.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func generateUserID() (string, error) {
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return "", err
	}

	return fmt.Sprintf("user_%s", hex.EncodeToString(idBytes)), nil
}
