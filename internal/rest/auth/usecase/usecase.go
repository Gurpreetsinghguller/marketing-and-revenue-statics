package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const tokenTTL = 24 * time.Hour

type tokenClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// AuthUseCase handles authentication business logic
type AuthUseCase struct {
	userRepo domain.UserRepo
}

func NewAuthUseCase(userRepo domain.UserRepo) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
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

	token, err := generateJWT(newUser.ID, newUser.Role)
	if err != nil {
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

	token, err := generateJWT(user.ID, user.Role)
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

func generateJWT(userID string, role domain.Role) (string, error) {
	secret := loadJWTSecret()
	if secret == "" {
		return "", errors.New("jwt secret is missing")
	}

	now := time.Now().UTC()
	claims := tokenClaims{
		Role: string(role),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(secret))
}

func loadJWTSecret() string {
	if secret := strings.TrimSpace(os.Getenv("JWT_SECRET")); secret != "" {
		return secret
	}

	data, err := os.ReadFile("shared/secret")
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}
