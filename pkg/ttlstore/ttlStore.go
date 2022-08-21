package ttlstore

import (
	"time"
)

type TTLStoreEntity interface {
	GetTTL() time.Time
	SetTTL(time.Time)
}


