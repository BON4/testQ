package ttlstore

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"encoding/hex"
	"io"
	"os"
	"sync"
	"time"
)

var SEP, _ = hex.DecodeString("5fff81030101406d6170456e74697479")
var SEP_LEN = len(SEP)

type mapEntity[K string, V any] struct {
	Key K
	Val V
}

func mapEtitySplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	dataLen := len(data)

	if atEOF && dataLen == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, SEP); i >= 0 {
		return i + SEP_LEN, data[0:i], nil
	}

	if atEOF {
		return dataLen, data, bufio.ErrFinalToken
	}

	return 0, nil, nil
}

type MapStore[K string, V any] struct {
	store *sync.Map
	ctx   context.Context
	save  chan mapEntity[K, V]
	cfg   TTLStoreConfig
}

func runSaveDaemon[K string, V any](ctx context.Context, kv chan mapEntity[K, V], path string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	for {
		select {
		case <-ctx.Done():
			//SET FLAG TO DONE. WHEN ALL SAVED EXIT.
			return
		case data := <-kv:

			//TODO: provide using custom encode algorithm
			if err := enc.Encode(data); err != nil {
				//TODO: Log err
			}
		}
	}
}

func runGcDaemon[K string, V any](store *sync.Map, ctx context.Context, dRt time.Duration) {
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

func NewMapStore[K string, V any](ctx context.Context, cfg TTLStoreConfig) *MapStore[K, V] {
	ms := &MapStore[K, V]{
		store: &sync.Map{},
		ctx:   ctx,
		//TODO: CHANEL SIZE?
		save: make(chan mapEntity[K, V], 100),
		cfg:  cfg,
	}

	go runGcDaemon[K, V](ms.store, ms.ctx, cfg.MapStore.GCRefresh)
	go runSaveDaemon[K, V](ctx, ms.save, cfg.MapStore.SavePath)

	return ms
}

func (ms *MapStore[K, V]) Load() error {
	file, err := os.OpenFile(ms.cfg.MapStore.SavePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	fileBuf := bufio.NewScanner(file)

	fileBuf.Split(mapEtitySplitFunc)

	var ent mapEntity[K, V]
	for fileBuf.Scan() {
		b := fileBuf.Bytes()
		if len(b) > 0 {
			dec := gob.NewDecoder(bytes.NewReader(append(SEP, b...)))
			for {
				if err := dec.Decode(&ent); err != nil {
					if err == io.EOF {
						break
					} else {
						return err
					}
				}
				ms.store.Store(ent.Key, ent.Val)
			}
		}
	}
	return nil
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
	ms.save <- mapEntity[K, V]{Key: key, Val: val}
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
