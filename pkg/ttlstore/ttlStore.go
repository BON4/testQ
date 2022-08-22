package ttlstore

import (
	"context"
	"time"
)

type TTLStoreEntity[T any] struct {
	Entity T
	ttl    int64
}

func (te TTLStoreEntity[T]) GetTTL() int64 {
	return te.ttl
}

func (te *TTLStoreEntity[T]) SetTTL(ttl int64) {
	te.ttl = ttl
}

type TTLStore[K string, V any] interface {
	Get(ctx context.Context, key K) (V, bool)
	Set(ctx context.Context, key K, val V, ttl time.Duration) error
}
