package entity


type Entity struct {
	Payload string `json:"payload"`
	TTL     int64  `json:"ttl"`
}

//TODO: compare perfomanse with passing time by ref (*int64)
//or by using UNIX timestamp
func (e Entity) GetTTL() int64  {
	return e.TTL
}

func (e *Entity) SetTTL(t int64) {
	e.TTL = t
}



