package did

import (
	"context"
	"time"
)

var DefaultCacheExpiry = 10 * time.Minute

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttls ...time.Duration) error
	Get(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
}
