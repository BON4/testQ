package ttlstore

import (
	"sync"
	"context"
	"time"
)

type MapStore[K any, V TTLStoreEntity] struct {
	store *sync.Map
	ctx context.Context
}

func(ms* MapStore[T,V]) runDaemon(store *sync.Map, ctx context.Context, dRt time.Duration) {
	tiker := time.NewTicker(dRt)
	for {
		select {
		case <-ctx.Done(): return 
		case <-tiker.C:
			store.Range(func(k, v any) bool {
				if val, ok := v.(V); ok {
					eTime := val.GetTTL()
					if !eTime.IsZero() && eTime.After(time.Now()) {
						store.Delete(k)
					}
				}
				return true
			})
		}
	}
}

func NewMapStore[K any, V TTLStoreEntity](ctx context.Context, daemonRefreshTime time.Duration) TTLStore {
	ms := &MapStore[K, V]{
		store: &sync.Map{},
		ctx: ctx,
	}

	go ms.runDaemon(ms.store, ms.ctx, daemonRefreshTime)
	return ms
}


func (ms* MapStore[K, V]) Set(key K, val V, ttl time.Duration) {
	t := time.Time{}
	if ttl == 0 {
		return
	} else if ttl > 0 {
		t = time.Now().Add(ttl)		
	}

	val.SetTTL(&t)
	ms.store.Store(key, val)
}

func (ms* MapStore[K, V]) Get(key K) (V, bool) {
	var ent V
	if val, ok := ms.store.Load(key); ok {
		if ent, ok := val.(V); ok {
			eTime := ent.GetTTL()
			if eTime.IsZero() && eTime.Before(time.Now()) {
				return ent, true
			}			
		}
	}
	return ent, false
}
