package entity

import (
	"time"
)

type Entity struct {
	Payload string     `json:"payload"`
	TTL *time.Duration `json:"ttl"`
}

func (e Entity) GetTTL() *time.Duration {
	return e.TTL
}




