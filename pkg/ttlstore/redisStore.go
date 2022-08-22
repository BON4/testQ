package ttlstore

import (
	"context"

	"time"

	"github.com/BON4/timedQ/pkg/buffpool"
	redis "github.com/go-redis/redis/v9"
)

type RedisStore[K string, V any] struct {
	cl *redis.Client
	bp buffpool.BufferPool
}

func NewRedisStore[K string, V any](ctx context.Context, cfg TTLStoreConfig) TTLStore[K, V] {
	return &RedisStore[K, V]{cl: redis.NewClient(&redis.Options{
		Addr:     cfg.RedisStore.Addr,
		Password: cfg.RedisStore.Password,
		DB:       cfg.RedisStore.DB,
	})}
}

func (rs *RedisStore[K, V]) Set(ctx context.Context, key K, val V, ttl time.Duration) error {

	return rs.cl.SetEx(ctx, string(key), val, ttl).Err()
}

func (rs *RedisStore[K, V]) Get(ctx context.Context, key K) (V, bool) {
	rs.cl.Get(ctx, string(key)).Result()
}
