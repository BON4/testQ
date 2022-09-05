package coder

import (
	"bytes"
	"encoding/gob"
	"io"
)

// Can be changed, by increceing scanner buffer size manualy
const SCAN_BUFFER_CAP = 64 * 1024

type Encoder[T any] struct {
	encoder    *gob.Encoder
	sizeWriter *bytes.Buffer
	w          io.Writer

	bytesCounter uint64
	objectSize   uint64
}

func NewEncoder[T any](w io.Writer) *Encoder[T] {
	b := bytes.NewBuffer([]byte{})
	return &Encoder[T]{
		sizeWriter: b,
		w:          w,
		encoder:    gob.NewEncoder(io.MultiWriter(w, b)),
	}
}

func (c *Encoder[T]) Encode(data *T) error {
	// Reset Encoder
	if c.bytesCounter+c.objectSize >= SCAN_BUFFER_CAP {
		// fmt.Println("RESET")
		c.encoder = gob.NewEncoder(io.MultiWriter(c.w, c.sizeWriter))
		c.bytesCounter = 0
	}

	if err := c.encoder.Encode(*data); err != nil {
		//TODO: Log err
		return err
	}

	c.bytesCounter += uint64(c.sizeWriter.Len())

	if c.objectSize == 0 {
		c.objectSize = c.bytesCounter
	}

	c.sizeWriter.Reset()
	return nil
}
