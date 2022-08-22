package ttlstore

import (
	"context"
	"testing"
	"time"

	models "github.com/BON4/timedQ/internal/models"
)

func TestRedisGetSet(t *testing.T) {
	cfg := newRedisStoreConfig("localhost:6379", "", 0)
	ety := &models.Entity{
		Payload: "Hello",
	}

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
