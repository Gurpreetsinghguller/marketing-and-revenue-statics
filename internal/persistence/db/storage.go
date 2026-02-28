package db

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Ensure StorageMgr implements PersistenceDB
var _ PersistenceDB = (*StorageMgr)(nil)

// Storage paths
const (
	DataDir       = "data"
	UsersFile     = "data/users.json"
	CampaignsFile = "data/campaigns.json"
	EventsFile    = "data/events.json"
	ReportsFile   = "data/reports.json"
)

// Common errors
var (
	ErrNotFound = errors.New("resource not found")
	ErrExists   = errors.New("resource already exists")
)

// StorageMgr manages JSON file storage
type StorageMgr struct {
	mu      sync.RWMutex // Ensure thread-safe file operations
	baseDir string
}

// NewStorageMgr creates a new storage manager
func NewStorageMgr() *StorageMgr {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(DataDir, 0755); err != nil {
		panic("Failed to create data directory: " + err.Error())
	}
	return &StorageMgr{baseDir: DataDir}
}

func (s *StorageMgr) resolvePath(key string) (string, error) {
	cleanKey := filepath.Clean(strings.TrimPrefix(strings.TrimSpace(key), "/"))
	if cleanKey == "." || cleanKey == "" {
		return "", errors.New("invalid key")
	}
	if strings.HasPrefix(cleanKey, "..") || filepath.IsAbs(cleanKey) {
		return "", errors.New("invalid key")
	}

	return filepath.Join(s.baseDir, cleanKey+".json"), nil
}

// Create stores a value at key and fails if it already exists.
func (s *StorageMgr) Create(ctx context.Context, key string, value interface{}) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	path, err := s.resolvePath(key)
	if err != nil {
		return err
	}

	if _, err = os.Stat(path); err == nil {
		return ErrExists
	}
	if !os.IsNotExist(err) {
		return err
	}

	jsonData, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0644)
}

// Read retrieves a value by key.
func (s *StorageMgr) Read(ctx context.Context, key string) (interface{}, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()

	path, err := s.resolvePath(key)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if len(data) == 0 {
		return map[string]interface{}{}, nil
	}

	var value interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}

	return value, nil
}

// Update updates or replaces the value at key.
func (s *StorageMgr) Update(ctx context.Context, key string, value interface{}) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	path, err := s.resolvePath(key)
	if err != nil {
		return err
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0644)
}

// Delete removes a value by key.
func (s *StorageMgr) Delete(ctx context.Context, key string) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	path, err := s.resolvePath(key)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

// List returns all values with keys matching prefix.
func (s *StorageMgr) List(ctx context.Context, prefix string) ([]interface{}, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()

	cleanPrefix := filepath.Clean(strings.TrimPrefix(strings.TrimSpace(prefix), "/"))
	if cleanPrefix == "." {
		cleanPrefix = ""
	}
	if strings.HasPrefix(cleanPrefix, "..") || filepath.IsAbs(cleanPrefix) {
		return nil, errors.New("invalid prefix")
	}

	basePath := s.baseDir
	if cleanPrefix != "" {
		basePath = filepath.Join(s.baseDir, cleanPrefix)
	}

	entries := make([]interface{}, 0)
	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			return nil
		}

		var value interface{}
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		entries = append(entries, value)
		return nil
	})

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) || os.IsNotExist(err) {
			return []interface{}{}, nil
		}
		return nil, err
	}

	return entries, nil
}

// Close releases storage resources.
func (s *StorageMgr) Close() error {
	return nil
}

// Helper function to generate UUID
func GenerateID(prefix string) string {
	return prefix + "_" + randomString(12)
}

// Simple random string generator
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
