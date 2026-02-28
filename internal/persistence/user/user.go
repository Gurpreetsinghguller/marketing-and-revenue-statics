package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/db"
)

const userPrefix = "users"

type userDBModel struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Bio       string    `json:"bio"`
	Phone     string    `json:"phone"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toUserDBModel(u domain.User) userDBModel {
	return userDBModel{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		Name:      u.Name,
		Role:      string(u.Role),
		Bio:       u.Bio,
		Phone:     u.Phone,
		Picture:   u.Picture,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (m userDBModel) toDomain() domain.User {
	return domain.User{
		ID:        m.ID,
		Email:     m.Email,
		Password:  m.Password,
		Name:      m.Name,
		Role:      domain.Role(m.Role),
		Bio:       m.Bio,
		Phone:     m.Phone,
		Picture:   m.Picture,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func decodeUserDBModel(value interface{}) (userDBModel, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return userDBModel{}, err
	}
	var model userDBModel
	if err := json.Unmarshal(b, &model); err != nil {
		return userDBModel{}, err
	}
	return model, nil
}

func userKey(id string) string {
	return userPrefix + "/" + id
}

// UserRepository implements domain.UserRepo.
type UserRepository struct {
	storage db.PersistenceDB
}

// NewUserRepository creates a new user repository.
func NewUserRepository(storage ...db.PersistenceDB) *UserRepository {
	selected := db.PersistenceDB(db.NewStorageMgr())
	if len(storage) > 0 && storage[0] != nil {
		selected = storage[0]
	}
	return &UserRepository{storage: selected}
}

func (r *UserRepository) getAll() ([]domain.User, error) {
	stored, err := r.storage.List(context.Background(), userPrefix)
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(stored))
	for _, item := range stored {
		model, err := decodeUserDBModel(item)
		if err != nil {
			return nil, err
		}
		users = append(users, model.toDomain())
	}

	return users, nil
}

// Create saves a new user.
func (r *UserRepository) Create(user *domain.User) error {
	if r.EmailExists(user.Email) {
		return db.ErrExists
	}

	if user.ID == "" {
		user.ID = db.GenerateID("user")
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return r.storage.Create(context.Background(), userKey(user.ID), toUserDBModel(*user))
}

// GetByID retrieves a user by ID.
func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	stored, err := r.storage.Read(context.Background(), userKey(id))
	if err != nil {
		return nil, err
	}

	model, err := decodeUserDBModel(stored)
	if err != nil {
		return nil, err
	}

	entity := model.toDomain()
	return &entity, nil
}

// GetByEmail retrieves a user by email.
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	users, err := r.getAll()
	if err != nil {
		return nil, err
	}
	for i := range users {
		if users[i].Email == email {
			return &users[i], nil
		}
	}
	return nil, db.ErrNotFound
}

// Update updates an existing user.
func (r *UserRepository) Update(user *domain.User) error {
	existing, err := r.GetByID(user.ID)
	if err != nil {
		return err
	}

	if user.CreatedAt.IsZero() {
		user.CreatedAt = existing.CreatedAt
	}
	user.UpdatedAt = time.Now()

	return r.storage.Update(context.Background(), userKey(user.ID), toUserDBModel(*user))
}

// Delete removes a user.
func (r *UserRepository) Delete(id string) error {
	return r.storage.Delete(context.Background(), userKey(id))
}

// GetAll retrieves all users.
func (r *UserRepository) GetAll() ([]domain.User, error) {
	return r.getAll()
}

// GetByRole retrieves users by role.
func (r *UserRepository) GetByRole(role string) ([]domain.User, error) {
	users, err := r.getAll()
	if err != nil {
		return nil, err
	}

	result := make([]domain.User, 0)
	for i := range users {
		if string(users[i].Role) == role {
			result = append(result, users[i])
		}
	}

	return result, nil
}

// EmailExists checks if email already exists.
func (r *UserRepository) EmailExists(email string) bool {
	_, err := r.GetByEmail(email)
	return err == nil
}
