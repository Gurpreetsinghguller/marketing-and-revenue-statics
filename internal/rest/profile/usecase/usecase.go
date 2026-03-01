package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
)

// ProfileUseCase handles user profile business logic
type ProfileUseCase struct {
	userRepo domain.UserRepo
}
type ProfileUseCaseInterface interface {
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID string, updates *domain.User) (*domain.User, error)
}

// NewProfileUseCase creates a new profile usecase
func NewProfileUseCase(userRepo domain.UserRepo) *ProfileUseCase {
	return &ProfileUseCase{
		userRepo: userRepo,
	}
}

// GetProfile retrieves user profile
func (p *ProfileUseCase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	_ = ctx

	user, err := p.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

// UpdateProfile updates user profile
func (p *ProfileUseCase) UpdateProfile(ctx context.Context, userID string, updates *domain.User) (*domain.User, error) {
	if updates == nil {
		return nil, errors.New("user profile updates are required")
	}

	// Fetch existing user
	user, err := p.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if updates.Name != "" {
		user.Name = updates.Name
	}
	if updates.Bio != "" {
		user.Bio = updates.Bio
	}
	if updates.Phone != "" {
		user.Phone = updates.Phone
	}
	if updates.Picture != "" {
		user.Picture = updates.Picture
	}

	// Save updated profile
	if err := p.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return user, nil
}
