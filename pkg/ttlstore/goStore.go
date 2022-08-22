package ttlstore

import (
	"context"
	"sync"
	"time"
)

type MapStore[K string, V any] struct {
	store *sync.Map
	ctx   context.Context
}

func (ms *MapStore[T, V]) runDaemon(store *sync.Map, ctx context.Context, dRt time.Duration) {
	tiker := time.NewTicker(dRt)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tiker.C:
			store.Range(func(k, v any) bool {
				if val, ok := v.(TTLStoreEntity[V]); ok {
					eTime := val.GetTTL()
					if !(eTime <= 0) && eTime < time.Now().Unix() {
						store.Delete(k)
					}
				}
				return true
			})
		}
	}
}

func NewMapStore[K string, V any](ctx context.Context, cfg TTLStoreConfig) TTLStore[K, V] {
	ms := &MapStore[K, V]{
		store: &sync.Map{},
		ctx:   ctx,
	}

	go ms.runDaemon(ms.store, ms.ctx, cfg.MapStore.GCRefresh)
	return ms
}

func (ms *MapStore[K, V]) Set(_ context.Context, key K, val V, ttl time.Duration) error {
	var t int64 = -1
	if ttl == 0 {
		return nil
	} else if ttl > 0 {
		t = time.Now().Add(ttl).Unix()
	}

	se := TTLStoreEntity[V]{
		Entity: val,
	}

	se.SetTTL(t)
	ms.store.Store(key, se)

	return nil
}

func (ms *MapStore[K, V]) Get(_ context.Context, key K) (V, bool) {
	var ent TTLStoreEntity[V]
	if val, ok := ms.store.Load(key); ok {
		if ent, ok := val.(TTLStoreEntity[V]); ok {
			eTime := ent.GetTTL()
			// 0 | 0 -> 0
			// 1 | 0 -> 1
			// 0 | 1 -> 1
			// 1 | 1 -> 1
			if eTime > time.Now().Unix() || (eTime <= 0) {
				return ent.Entity, true
			}
		}
	}
	return ent.Entity, false
}
