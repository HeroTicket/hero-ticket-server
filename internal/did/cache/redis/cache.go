package redis

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/heroticket/internal/did"
)

type redisCache struct {
	c *cache.Cache
}

func New(c *cache.Cache) did.Cache {
	return &redisCache{c: c}
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.c.Delete(ctx, key)
}

func (r *redisCache) Get(ctx context.Context, key string, value interface{}) error {
	return r.c.Get(ctx, key, value)
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}, ttls ...time.Duration) error {
	item := &cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
	}

	if len(ttls) > 0 {
		item.TTL = ttls[0]
	}

	return r.c.Set(item)
}
