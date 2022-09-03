package ttlstore

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"math/rand"

	models "github.com/BON4/timedQ/internal/models"
)

type TestStruct struct {
	A int64
	B int64
	C int64
}

func TestMapLarge(t *testing.T) {
	n := 5
	m := 200

	filename := "#temp.db"
	os.Remove(filename)

	cfg := NewMapStoreConfig(time.Second/3, 1, filename, true)

	for i := 0; i < n; i++ {
		ms := NewMapStore[string, *TestStruct](context.Background(), cfg)
		if err := ms.Load(); err != nil {
			t.Error(err)
			return
		}

		if err := ms.Run(); err != nil {
			t.Error(err)
			return
		}

		for i := 0; i < m; i++ {
			ety := &TestStruct{
				A: rand.Int63(),
				B: rand.Int63(),
				C: rand.Int63(),
			}

			ms.Set(context.Background(), fmt.Sprintf("%d", i), ety, -1)
		}

		time.Sleep(time.Second / 2)

		ms.Close()

		if stat, err := os.Stat(filename); err == nil {
			t.Logf("File Size: %d mb", stat.Size()/1000)
		} else {
			t.Error(err)
		}

	}

	if err := os.Remove(filename); err != nil {
		t.Error(err)
	}
}

func TestMapGetSet(t *testing.T) {
	ety := &models.Entity{
		Payload: "Hello\n",
	}

	cfg := NewMapStoreConfig(time.Second/3, 1, "#temp.db", false)

	ms := NewMapStore[string, *models.Entity](context.Background(), cfg)

	if err := ms.Run(); err != nil {
		t.Error(err)
	}

	defer ms.Close()

	ms.Set(context.Background(), "test", ety, time.Second*2)
	time.Sleep(time.Second)
	if providedEty, ok := ms.Get(context.Background(), "test"); ok {
		if !(ety.Payload == providedEty.Payload) {
			t.Log("Payloads dont match")
		}
	} else {
		t.Error("Cant get entity")
	}

	time.Sleep(time.Second * 2)

	if _, ok := ms.Get(context.Background(), "test"); ok {
		t.Error("Expected error, entity not deleted after ttl")
	}
}

func TestMultipleInstans(t *testing.T) {
	n := 1000
	ents := make(map[string]models.Entity, n)

	ctx := context.Background()

	cfg1 := NewMapStoreConfig(time.Second/3, 1, "#temp1.db", false)
	cfg2 := NewMapStoreConfig(time.Second/3, 1, "#temp2.db", false)

	store1 := NewMapStore[string, *models.Entity](context.Background(), cfg1)

	if err := store1.Run(); err != nil {
		t.Error(err)
	}

	defer store1.Close()

	store2 := NewMapStore[string, *models.Entity](context.Background(), cfg2)

	if err := store2.Run(); err != nil {
		t.Error(err)
	}

	defer store2.Close()

	for i := 0; i < n; i++ {
		key := fmt.Sprintf("%d", i)
		val := models.Entity{
			Payload: fmt.Sprintf("test:%d", i),
		}
		ents[key] = val

		if i < n/2 {
			if err := store1.Set(ctx, key, &val, -1); err != nil {
				t.Error(err)
			}
		} else {
			if err := store2.Set(ctx, key, &val, -1); err != nil {
				t.Error(err)
			}
		}
	}

	for k, v := range ents {
		if providedEty, ok := store1.Get(ctx, k); ok {
			if !(v.Payload == providedEty.Payload) {
				t.Log("Payloads dont match")
			}
		} else if providedEty, ok := store2.Get(ctx, k); ok {
			if !(v.Payload == providedEty.Payload) {
				t.Log("Payloads dont match")
			}
		} else {
			t.Errorf("Cant find entry with k:%s v:%s", k, v.Payload)
		}
	}

}

// TODO: Create proper test for Load
func TestMapLoad(t *testing.T) {
	filename := "#temp.db"
	os.Remove(filename)

	cfg := NewMapStoreConfig(time.Second/3, 1, filename, true)

	for i := 0; i < 3; i++ {
		ms := NewMapStore[string, *models.Entity](context.Background(), cfg)

		if err := ms.Run(); err != nil {
			t.Error(err)
		}

		for i := 0; i < 5; i++ {
			ety := &models.Entity{
				Payload: fmt.Sprintf("test:%d", i),
			}

			ms.Set(context.Background(), fmt.Sprintf("%d", i), ety, -1)
		}

		time.Sleep(time.Second / 2)

		ms.Close()
	}

	if stat, err := os.Stat(filename); err == nil {
		t.Logf("File Size: %d mb", stat.Size()/1000)
	} else {
		t.Error(err)
	}

	newMs := NewMapStore[string, *models.Entity](context.Background(), cfg)
	if err := newMs.Load(); err != nil {
		t.Error(err)
	}

	if err := os.Remove(filename); err != nil {
		t.Error(err)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func BenchmarkMapGetSet(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	cfg := NewMapStoreConfig(time.Second/3, 1, "#temp.db", false)

	ms := NewMapStore[string, *models.Entity](context.Background(), cfg)
	defer ms.Close()
	for i := 0; i < b.N; i++ {
		key := randSeq(10)
		ety := &models.Entity{
			Payload: randSeq(40),
		}
		ms.Set(context.Background(), key, ety, -1)
		if providedEty, ok := ms.Get(context.Background(), key); ok {
			if !(ety.Payload == providedEty.Payload) {
				b.Log("Payloads dont match")
			}
		} else {
			b.Error("Cant get entity")
		}
	}
}
