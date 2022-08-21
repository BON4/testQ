package ttlstore

import (
	"time"
)

type TTLStoreEntity interface {
	GetTTL() int64
	SetTTL(int64)
}

type TTLStore[K any, V TTLStoreEntity] interface {
	Get(key K) (V, bool)
	Set(key K, val V, ttl time.Duration)
}


