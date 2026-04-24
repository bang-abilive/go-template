package database

import "context"

type Database interface {
	Health(ctx context.Context) (map[string]string, error)
	Close(ctx context.Context) error
}
