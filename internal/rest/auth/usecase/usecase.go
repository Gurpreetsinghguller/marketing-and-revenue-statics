package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/db"
)

// AuthUseCase handles authentication business logic
type AuthUseCase struct {
	db db.PersistenceDB
}

// NewAuthUseCase creates a new auth usecase
func NewAuthUseCase(db db.PersistenceDB) *AuthUseCase {
	return &AuthUseCase{
		db: db,
	}
}

// RegisterRequest represents registration input
type RegisterRequest struct {
	Email    string
	Password string
	Name     string
	Role     string
}

// RegisterResponse represents registration output
type RegisterResponse struct {
	UserID string
	Token  string
}

// LoginRequest represents login input
type LoginRequest struct {
	Email    string
	Password string
}

// LoginResponse represents login output
type LoginResponse struct {
	Token string
	User  *domain.User
}

// Register creates a new user account
func (u *AuthUseCase) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, errors.New("email, password, and name are required")
	}

	// Check if email already exists
	key := fmt.Sprintf("user:email:%s", req.Email)
	exists, err := u.db.Read(ctx, key)
	if err == nil && exists != nil {
		return nil, errors.New("email already exists")
	}

	//TODO: Password Hashing and ID generation
	user := &domain.User{
		ID:       fmt.Sprintf("user_%d", len([]int{})),
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Role:     domain.Role(req.Role),
	}

	// Save user
	userKey := fmt.Sprintf("user:%s", user.ID)
	if err := u.db.Create(ctx, userKey, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Index by email for quick lookup
	if err := u.db.Create(ctx, key, user.ID); err != nil {
		return nil, fmt.Errorf("failed to index user: %w", err)
	}

	// Generate JWT token (simplified, use proper JWT in production)
	token := fmt.Sprintf("jwt_token_%s", user.ID)

	return &RegisterResponse{
		UserID: user.ID,
		Token:  token,
	}, nil
}

// Login authenticates a user
func (u *AuthUseCase) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	// Get user by email
	emailKey := fmt.Sprintf("user:email:%s", req.Email)
	userIDInterface, err := u.db.Read(ctx, emailKey)
	if err != nil || userIDInterface == nil {
		return nil, errors.New("invalid email or password")
	}

	userID, ok := userIDInterface.(string)
	if !ok {
		return nil, errors.New("invalid user data")
	}

	// Fetch user details
	userKey := fmt.Sprintf("user:%s", userID)
	userInterface, err := u.db.Read(ctx, userKey)
	if err != nil || userInterface == nil {
		return nil, errors.New("user not found")
	}

	user, ok := userInterface.(*domain.User)
	if !ok {
		return nil, errors.New("invalid user data format")
	}

	// Validate password (in production, use proper password hashing comparison)
	if user.Password != req.Password {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token (simplified)
	token := fmt.Sprintf("jwt_token_%s", user.ID)

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetUserByID retrieves user by ID
func (u *AuthUseCase) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	userKey := fmt.Sprintf("user:%s", userID)
	userInterface, err := u.db.Read(ctx, userKey)
	if err != nil || userInterface == nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	user, ok := userInterface.(*domain.User)
	if !ok {
		return nil, errors.New("invalid user data format")
	}

	return user, nil
}

// VerifyToken verifies JWT token and returns user ID (simplified)
func (u *AuthUseCase) VerifyToken(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", errors.New("token is required")
	}

	// Simplified token verification - in production use proper JWT library
	// For now, just extract userID from token format: "jwt_token_<userID>"
	if len(token) < 10 {
		return "", errors.New("invalid token format")
	}

	return "", nil // In production, properly validate JWT
}
