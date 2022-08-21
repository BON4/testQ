package ttlstore

import (
	"testing"
	"context"
	"time"
	"github.com/BON4/timedQ/internal/models"
)

func TestGet(t *testing.T) {
  ms := NewMapStore[string, models.Entity](context.Background(), time.Second)
	
}
