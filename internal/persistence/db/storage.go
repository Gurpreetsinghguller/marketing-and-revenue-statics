package db

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
)

// // Ensure StorageMgr implements PersistenceDB
// var _ PersistenceDB = (*StorageMgr)(nil)

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
	mu sync.RWMutex // Ensure thread-safe file operations
}

// NewStorageMgr creates a new storage manager
func NewStorageMgr() *StorageMgr {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(DataDir, 0755); err != nil {
		panic("Failed to create data directory: " + err.Error())
	}
	return &StorageMgr{}
}

// ReadJSON reads a JSON file and unmarshals into target
func (s *StorageMgr) ReadJSON(filePath string, target interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, return empty slice
			return nil
		}
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, target)
}

// WriteJSON marshals data and writes to JSON file
func (s *StorageMgr) WriteJSON(filePath string, data interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, jsonData, 0644)
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

// FindInSlice finds an item in a slice using a predicate function
func FindInSlice(slice interface{}, fn func(interface{}) bool) (interface{}, error) {
	// This is a placeholder - actual implementation depends on slice type
	return nil, ErrNotFound
}
