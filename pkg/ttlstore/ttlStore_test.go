package ttlstore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"math/rand"

	models "github.com/BON4/timedQ/internal/models"
)

func TestMapGetSet(t *testing.T) {
	ety := &models.Entity{
		Payload: "Hello\n",
	}

	cfg := newMapStoreConfig(time.Second/3, 1, "#temp.db")

	ms := NewMapStore[string, *models.Entity](context.Background(), cfg)
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

// TODO: Create proper test for Load
func TestMapLoad(t *testing.T) {
	cfg := newMapStoreConfig(time.Second/3, 1, "#temp.db")

	for i := 0; i < 4; i++ {
		ms := NewMapStore[string, *models.Entity](context.Background(), cfg)

		for i := 0; i < 5; i++ {
			ety := &models.Entity{
				Payload: fmt.Sprintf("test:%d", i),
			}

			ms.Set(context.Background(), fmt.Sprintf("%d", i), ety, -1)
		}

		time.Sleep(time.Second / 2)

		ms.Close()
	}

	newMs := NewMapStore[string, *models.Entity](context.Background(), cfg)
	if err := newMs.Load(); err != nil {
		t.Error(err)
	}
	newMs.Close()
}

func TestRedisGetSet(t *testing.T) {
	ety := &models.Entity{
		Payload: "Hello",
	}

	cfg := newRedisStoreConfig("localhost:6379", "", 0)

	reds := NewRedisStore[string, *models.Entity](context.Background(), cfg)
	reds.Set(context.Background(), "test", ety, time.Second*2)
	time.Sleep(time.Second)
	if providedEty, ok := reds.Get(context.Background(), "test"); ok {
		if !(ety.Payload == providedEty.Payload) {
			t.Log("Payloads dont match")
		}
	} else {
		t.Error("Cant get entity")
	}

	time.Sleep(time.Second * 2)

	if _, ok := reds.Get(context.Background(), "test"); ok {
		t.Error("Expected error, entity not deleted after ttl")
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
	cfg := newMapStoreConfig(time.Second/3, 1, "#temp.db")

	ms := NewMapStore[string, *models.Entity](context.Background(), cfg)

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

func BenchmarkRedisGetSet(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	cfg := newRedisStoreConfig("localhost:6379", "", 0)

	ms := NewRedisStore[string, *models.Entity](context.Background(), cfg)

	for i := 0; i < b.N; i++ {
		key := randSeq(10)
		ety := &models.Entity{
			Payload: randSeq(40),
		}
		if err := ms.Set(context.Background(), key, ety, 0); err != nil {
			b.Error(err)
		}
		if providedEty, ok := ms.Get(context.Background(), key); ok {
			if !(ety.Payload == providedEty.Payload) {
				b.Log("Payloads dont match")
			}
		} else {
			b.Error("Cant get entity")
		}
	}
}
