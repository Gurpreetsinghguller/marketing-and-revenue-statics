package db

import "context"

type PersistenceDB interface {
	Create(ctx context.Context, key string, value interface{}) error
	Read(ctx context.Context, key string) (interface{}, error)
	Update(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]interface{}, error)
	Close() error
}
