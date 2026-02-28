package db

import "context"

// this is Repository pattern
// abstracts the underlying database operations for the application
// The interface defines the "strategy"
type PersistenceDB interface {
	Create(ctx context.Context, key string, value interface{}) error
	Read(ctx context.Context, key string) (interface{}, error)
	Update(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]interface{}, error)
	Close() error
}
