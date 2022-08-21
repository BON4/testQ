package ttlstore

import (
	"testing"
	"context"
	"time"
	models "github.com/BON4/timedQ/internal/models"
)

func TestGet(t *testing.T) {
	ety := &models.Entity{
		Payload: "Hello",
	}
  ms := NewMapStore[string, *models.Entity](context.Background(), time.Second)
	ms.Set("test", ety, time.Minute)
 	if providedEty, ok := ms.Get("test"); ok {
		if ety.Payload == providedEty.Payload {
			return
		}
		t.Log("Payloads dont match")
	} else {
		t.Error("Cant get entity")
	}

	time.Sleep(time.Second*2)

	if _, ok := ms.Get("test"); ok {
		t.Error("Expected error, entity not deleted after ttl")
	}

}
