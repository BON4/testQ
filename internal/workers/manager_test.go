package workers

import (
	"context"
	"fmt"
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
		stores[i] = ttlstore.NewMapStore[string, string](ctx, ttlstore.NewMapStoreConfig(time.Second, 1, path, true))
		if err := stores[i].Load(); err != nil {
			t.Logf("Error at %s:", path)
			t.Error(err)
			return
		}
		defer stores[i].Close()
	}

	wmcfg := newManagerConfig(uint(storeCount), time.Minute)
	wm := NewWorkerManager(context.Background(), stores, wmcfg)

	wm.Run()

	res := wm.Get("test|")

	t.Logf("Got: %s", res)

	wm.Stop()
}
