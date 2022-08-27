package ttlstore

import (
	"context"
	"encoding/gob"

	"time"

	"bytes"

	redis "github.com/go-redis/redis/v9"
)

type RedisStore[K string, V any] struct {
	cl *redis.Client
}

func NewRedisStore[K string, V any](ctx context.Context, cfg TTLStoreConfig) TTLStore[K, V] {
	return &RedisStore[K, V]{cl: redis.NewClient(&redis.Options{
		Addr:     cfg.RedisStore.Addr,
		Password: cfg.RedisStore.Password,
		DB:       cfg.RedisStore.DB,
	})}
}

func (rs *RedisStore[K, V]) Close() error {
	return rs.cl.Close()
}

func (rs *RedisStore[K, V]) Load() error {
	return nil
}

func (rs *RedisStore[K, V]) Set(ctx context.Context, key K, val V, ttl time.Duration) error {
	// TODO: maby intoduce interface to allow user implement encode/decode
	switch any(val).(type) {
	case string:
		return rs.cl.Set(ctx, string(key), val, ttl).Err()
	case []byte:
		return rs.cl.Set(ctx, string(key), val, ttl).Err()
	default:
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(val); err != nil {
			return err
		}

		return rs.cl.Set(ctx, string(key), buf.Bytes(), ttl).Err()
	}
}

func (rs *RedisStore[K, V]) Get(ctx context.Context, key K) (V, bool) {
	var val V
	var buf bytes.Buffer
	res, err := rs.cl.Get(ctx, string(key)).Result()
	if err != nil {
		return val, false
	}
	buf.WriteString(res)
	dec := gob.NewDecoder(&buf)
	if err := dec.Decode(&val); err != nil {
		return val, false
	}

	return val, true
}
