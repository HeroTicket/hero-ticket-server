package redis

import (
	"context"
	"time"

	redis "github.com/go-redis/cache/v9"
	"github.com/heroticket/internal/cache"
)

type redisCache struct {
	c *redis.Cache
}

func NewCache(c *redis.Cache) cache.Cache {
	return &redisCache{c: c}
}

func (r *redisCache) Exists(ctx context.Context, key string) bool {
	return r.c.Exists(ctx, key)
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	err := r.c.Delete(ctx, key)
	if err == redis.ErrCacheMiss {
		return nil
	}
	return err
}

func (r *redisCache) Get(ctx context.Context, key string, value interface{}) error {
	err := r.c.Get(ctx, key, value)
	if err == redis.ErrCacheMiss {
		return cache.ErrCacheMiss
	}
	return err
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}, ttls ...time.Duration) error {
	item := &redis.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
	}

	if len(ttls) > 0 {
		item.TTL = ttls[0]
	}

	return r.c.Set(item)
}
