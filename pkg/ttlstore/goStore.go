package ttlstore

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BON4/timedQ/pkg/coder"
)

const DEFAULT_DUMP_NAME = ".temp.db"

// Separator in gob encoded file, where 6d6170456e74697479 - _.....mapEntity in HEX
var SEP, _ = hex.DecodeString("6d6170456e74697479")

// var SEP, _ = hex.DecodeString("5fff81030101406d6170456e74697479")
var PRE_SEP_LEN = 16
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
	wg       *sync.WaitGroup
	cancel   context.CancelFunc
	store    *sync.Map
	ctx      context.Context
	save     chan mapEntity[K, TTLStoreEntity[V]]
	cfg      TTLStoreConfig
	dump     *os.File
	dumpPath string
}

func runSaveDaemon[K string, V any](ctx context.Context, kv chan mapEntity[K, TTLStoreEntity[V]], wg *sync.WaitGroup, file io.Writer) {
	wg.Add(1)

	defer wg.Done()

	encoder := coder.NewEncoder[*mapEntity[K, TTLStoreEntity[V]]](file)

	for {
		select {
		case data := <-kv:
			if err := encoder.Encode(&data); err != nil {
				panic(err)
			}
		case <-ctx.Done():
			//SET FLAG TO DONE. WHEN ALL SAVED EXIT.
			return
		}
	}
}

func runGcDaemon[K string, V any](ctx context.Context, store *sync.Map, wg *sync.WaitGroup, dRt time.Duration) {
	wg.Add(1)

	defer wg.Done()

	tiker := time.NewTicker(dRt)
	for {
		select {
		case <-tiker.C:
			store.Range(func(k, v any) bool {
				if val, ok := v.(TTLStoreEntity[V]); ok {
					eTime := val.GetTTL()
					if !(eTime <= 0) && eTime < time.Now().Unix() {
						store.Delete(k)
					}
				} else {
					// TODO: Make proper logger
					fmt.Println("Invalid type value found while saving")
				}
				return true
			})
		case <-ctx.Done():
			return
		}
	}
}

// TODO: handle error
func NewMapStore[K string, V any](ctx context.Context, cfg TTLStoreConfig) *MapStore[K, V] {
	msctx, cancel := context.WithCancel(ctx)

	ms := &MapStore[K, V]{
		store:  &sync.Map{},
		ctx:    msctx,
		cancel: cancel,
		//TODO: CHANEL SIZE?
		save: make(chan mapEntity[K, TTLStoreEntity[V]], 100),
		cfg:  cfg,
		dump: nil,
		wg:   &sync.WaitGroup{},
	}

	dir, fname := filepath.Split(cfg.MapStore.SavePath)
	if len(fname) == 0 {
		ms.dumpPath = dir + DEFAULT_DUMP_NAME
	} else {
		ms.dumpPath = cfg.MapStore.SavePath
	}

	go runGcDaemon[K, V](ms.ctx, ms.store, ms.wg, ms.cfg.MapStore.GCRefresh)
	return ms
}

func (ms *MapStore[K, V]) Path() string {
	return ms.dumpPath
}

// Run - runs a damon that saves map content to file
// Call Run, only in case of where cfg.MapStore.Save == true
// WRRNING: Load shoud be called before Run
func (ms *MapStore[K, V]) Run() error {
	if ms.cfg.MapStore.Save {
		var err error
		ms.dump, err = os.OpenFile(ms.dumpPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println(err)
			return err
		}

		go runSaveDaemon[K, V](ms.ctx, ms.save, ms.wg, ms.dump)
	}

	return nil
}

func (ms *MapStore[K, V]) Close() error {
	ms.cancel()
	ms.wg.Wait()
	if ms.cfg.MapStore.Save {
		//TODO: Error when close without run
		return ms.dump.Close()
	}
	return nil
}

// Load - loads all contents from file to internal map, then clears a file and dump all contents to fresh file
// WRRNING: Load shoud be called before Run
func (ms *MapStore[K, V]) Load() error {
	// Using custom split function we will get output in bytes where:
	// FIRST SCAN:
	// *******@mapEntity
	// ^_____^ - this is pre_separator (always len of 7)
	//
	// SECOND SCAN:
	// separator ... pre_sepator
	// ..............^__________ - now pre_separator will be at the end, we need to trim it from end and append it to start.
	if ms.cfg.MapStore.Save {
		var err error
		reader, err := os.OpenFile(ms.dumpPath, os.O_CREATE|os.O_RDONLY, 0666)
		if err != nil {
			return err
		}

		decoder := coder.NewDecoder[*mapEntity[K, TTLStoreEntity[V]]](reader, SEP)
		decoder.Decode(func(ent *mapEntity[K, TTLStoreEntity[V]]) {
			ms.store.Store(ent.Key, ent.Val)
		})
		// fileBuf := bufio.NewScanner(reader)

		// fileBuf.Split(mapEtitySplitFunc)

		// var ent mapEntity[K, TTLStoreEntity[V]]

		// var pre_separator []byte = SEP
		// for fileBuf.Scan() {
		// 	b := fileBuf.Bytes()
		// 	if len(b) > 8 {
		// 		encoded := append(pre_separator, bytes.TrimRight(b, string(pre_separator))...)

		// 		dec := gob.NewDecoder(bytes.NewReader(encoded))
		// 		for {
		// 			if err := dec.Decode(&ent); err != nil {
		// 				if err == io.EOF {
		// 					break
		// 				} else {
		// 					return err
		// 				}
		// 			}

		// 			ms.store.Store(ent.Key, ent.Val)
		// 		}
		// 	} else {
		// 		// Store pre_separator value
		// 		pre_separator = make([]byte, len(b))
		// 		copy(pre_separator, b)
		// 		pre_separator = append(pre_separator, SEP...)
		// 	}
		// }

		reader.Close()

		writer, err := os.OpenFile(ms.dumpPath, os.O_TRUNC|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println(err)
			return err
		}

		defer writer.Close()
		encoder := coder.NewEncoder[*mapEntity[K, TTLStoreEntity[V]]](writer)

		ms.store.Range(func(key any, val any) bool {
			if okKey, ok := key.(K); ok {
				if okVal, ok := val.(TTLStoreEntity[V]); ok {
					if err = encoder.Encode(&mapEntity[K, TTLStoreEntity[V]]{
						Key: okKey,
						Val: okVal,
					}); err != nil {
						return false
					}
				}
			}
			return true
		})

		if err != nil {
			return err
		}
		//encoder.Encode(
		// 	double := bytes.NewBuffer([]byte{})

		// 	enc := gob.NewEncoder(io.MultiWriter(writer, double))

		// 	var bytesCounter uint64 = 0
		// 	var objectSize uint64 = 0
		// 	ms.store.Range(func(key any, val any) bool {
		// 		if okKey, ok := key.(K); ok {
		// 			if okVal, ok := val.(TTLStoreEntity[V]); ok {
		// 				if bytesCounter+objectSize >= SCAN_BUFFER_CAP {
		// 					enc = gob.NewEncoder(io.MultiWriter(writer, double))
		// 					bytesCounter = 0
		// 				}

		// 				if err := enc.Encode(mapEntity[K, TTLStoreEntity[V]]{Key: okKey, Val: okVal}); err != nil {
		// 					panic(err)
		// 				}

		// 				bytesCounter += uint64(double.Len())

		// 				if objectSize == 0 {
		// 					objectSize = bytesCounter
		// 				}

		// 				double.Reset()
		// 			}
		// 		}

		// 		return true
		// 	})
		// }

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
	if ms.cfg.MapStore.Save {
		ms.save <- mapEntity[K, TTLStoreEntity[V]]{Key: key, Val: se}
	}

	return nil
}

func (ms *MapStore[K, V]) Get(_ context.Context, key K) (V, bool) {
	var ent TTLStoreEntity[V]
	if val, ok := ms.store.Load(key); ok {
		//fmt.Println("Loaded")
		if ent, ok := val.(TTLStoreEntity[V]); ok {
			eTime := ent.GetTTL()
			// 0 | 0 -> 0
			// 1 | 0 -> 1
			// 0 | 1 -> 1
			// 1 | 1 -> 1
			//fmt.Printf("Found with Get: %v\n", ent.Entity)
			if eTime > time.Now().Unix() || (eTime <= 0) {
				return ent.Entity, true
			} else {
				fmt.Printf("Cant assert, got: %+v", val)
			}
		}
	}
	return ent.Entity, false
}

func (ms *MapStore[K, V]) Range(f func(key K, val V) bool) {
	ms.store.Range(func(key any, value any) bool {
		if okKey, ok := key.(K); ok {
			if okVal, ok := value.(TTLStoreEntity[V]); ok {
				return f(okKey, okVal.Entity)
			}
		}
		return true
	})
}
