package ttlstore

type TTLStoreEntity[T any] struct {
	Entity T
	ttl    int64
}

func (te TTLStoreEntity[T]) GetTTL() int64 {
	return te.ttl
}

func (te *TTLStoreEntity[T]) SetTTL(ttl int64) {
	te.ttl = ttl
}
