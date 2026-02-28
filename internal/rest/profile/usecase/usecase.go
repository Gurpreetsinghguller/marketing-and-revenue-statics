package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/db"
)

// ProfileUseCase handles user profile business logic
type ProfileUseCase struct {
	db db.PersistenceDB
}

// NewProfileUseCase creates a new profile usecase
func NewProfileUseCase(db db.PersistenceDB) *ProfileUseCase {
	return &ProfileUseCase{
		db: db,
	}
}

// UpdateProfileRequest represents profile update input
type UpdateProfileRequest struct {
	Name    string
	Bio     string
	Phone   string
	Picture string
}

// GetProfile retrieves user profile
func (p *ProfileUseCase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	key := fmt.Sprintf("user:%s", userID)
	userInterface, err := p.db.Read(ctx, key)
	if err != nil || userInterface == nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	user, ok := userInterface.(*domain.User)
	if !ok {
		return nil, errors.New("invalid user data format")
	}

	return user, nil
}

// UpdateProfile updates user profile
func (p *ProfileUseCase) UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*domain.User, error) {
	// Fetch existing user
	user, err := p.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Picture != "" {
		user.Picture = req.Picture
	}

	// Save updated profile
	key := fmt.Sprintf("user:%s", userID)
	if err := p.db.Update(ctx, key, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return user, nil
}
