package workers

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/BON4/timedQ/pkg/ttlstore"
)

func TestManager(t *testing.T) {
	ctx := context.Background()

	storeCount := 5

	stores := make([]*ttlstore.MapStore[string, string], storeCount)

	for i := 0; i < len(stores); i++ {
		path := fmt.Sprintf("/home/home/go/src/timedQ/internal/workers/#temp%d.db", i)
		stores[i] = ttlstore.NewMapStore[string, string](ctx, ttlstore.NewMapStoreConfig(time.Second/3, 1, path, true))
		if err := stores[i].Load(); err != nil {
			t.Logf("Error at %s:", path)
			t.Error(err)
			return
		}

		if err := stores[i].Load(); err != nil {
			t.Error(err)
		}

		if err := stores[i].Run(); err != nil {
			t.Error(err)
		}

		defer stores[i].Close()
	}

	wmcfg := newManagerConfig(uint(storeCount), time.Minute)
	wm := NewWorkerManager(context.Background(), stores, wmcfg)

	wm.Run()

	//res := wm.Get("test|")

	//t.Logf("Got: %s", res)

	time.Sleep(time.Second * 2)

	wm.Stop()
}

func TestManagerBigChunks(t *testing.T) {
	ctx := context.Background()

	storeCount := 5

	stores := make([]*ttlstore.MapStore[string, string], storeCount)

	for i := 0; i < len(stores); i++ {
		path := fmt.Sprintf("/home/home/go/src/timedQ/internal/workers/#temp%d.db", i)
		stores[i] = ttlstore.NewMapStore[string, string](ctx, ttlstore.NewMapStoreConfig(time.Second/3, 1, path, true))
		if err := stores[i].Load(); err != nil {
			t.Logf("Error at %s:", path)
			t.Error(err)
			return
		}

		if err := stores[i].Run(); err != nil {
			t.Error(err)
		}

		defer stores[i].Close()
	}

	wmcfg := newManagerConfig(uint(storeCount), time.Minute)
	wm := NewWorkerManager(context.Background(), stores, wmcfg)

	wm.Run()

	nWorks := 10
	nReq := 50
	wg := &sync.WaitGroup{}

	//TODO: check for errors like:     manager_test.go:66: gob: duplicate type received

	for j := 0; j < nWorks; j++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for i := 0; i < nReq; i++ {
				wm.Get(fmt.Sprintf("test{%d}", i))
			}
		}(wg)
	}

	wg.Wait()

	time.Sleep(time.Second)

	wm.Stop()
}
