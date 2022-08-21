package ttlstore

import (
	"github.com/go-redis/redis/v9"
	"context"
)

type RedisStore struct {
	cl *redis.Client
}

func NewRedisStore(ctx context.Context, cfg TTLStoreConfig) TTLStore {
	return nil
}
