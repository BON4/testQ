package entity

import (
	"time"
)

type Entity struct {
	Payload string `json:"payload"`
	TTL *time.Time `json:"ttl"`
}

func (e Entity) GetTTL() *time.Time {
	return e.TTL
}

func (e *Entity) SetTTL(t *time.Time) {
	e.TTL = t
}



