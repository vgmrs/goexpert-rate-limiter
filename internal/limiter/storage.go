package limiter

import (
	"context"
)

type Storage interface {
	Increment(ctx context.Context, key string, expiration int64) (int64, error)
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration int64) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Close() error
}
