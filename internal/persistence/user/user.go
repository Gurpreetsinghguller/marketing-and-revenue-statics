package persistence

import (
	"math/rand"
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
)

// UserRepository implements domain.UserRepo using JSON file storage
type UserRepository struct {
	storage *StorageMgr
	users   []domain.User
}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	repo := &UserRepository{
		storage: NewStorageMgr(),
		users:   []domain.User{},
	}
	// Load existing users from file
	repo.storage.ReadJSON(UsersFile, &repo.users)
	return repo
}

// Create saves a new user
func (r *UserRepository) Create(user *domain.User) error {
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()

	// Check if email already exists
	for _, u := range r.users {
		if u.Email == user.Email {
			return ErrExists
		}
	}

	// Generate ID if not provided
	if user.ID == "" {
		user.ID = "user_" + generateRandomString(12)
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	r.users = append(r.users, *user)

	return r.storage.WriteJSON(UsersFile, r.users)
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	for i := range r.users {
		if r.users[i].ID == id {
			return &r.users[i], nil
		}
	}
	return nil, ErrNotFound
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	for i := range r.users {
		if r.users[i].Email == email {
			return &r.users[i], nil
		}
	}
	return nil, ErrNotFound
}

// Update updates an existing user
func (r *UserRepository) Update(user *domain.User) error {
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()

	for i := range r.users {
		if r.users[i].ID == user.ID {
			user.UpdatedAt = time.Now()
			r.users[i] = *user
			return r.storage.WriteJSON(UsersFile, r.users)
		}
	}
	return ErrNotFound
}

// Delete removes a user
func (r *UserRepository) Delete(id string) error {
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()

	for i := range r.users {
		if r.users[i].ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return r.storage.WriteJSON(UsersFile, r.users)
		}
	}
	return ErrNotFound
}

// GetAll retrieves all users
func (r *UserRepository) GetAll() ([]domain.User, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	return r.users, nil
}

// GetByRole retrieves users by role
func (r *UserRepository) GetByRole(role string) ([]domain.User, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	var result []domain.User
	for i := range r.users {
		if string(r.users[i].Role) == role {
			result = append(result, r.users[i])
		}
	}
	return result, nil
}

// EmailExists checks if email already exists
func (r *UserRepository) EmailExists(email string) bool {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	for _, u := range r.users {
		if u.Email == email {
			return true
		}
	}
	return false
}

// Helper function to generate random string
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
