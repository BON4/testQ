package coder

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io"
)

func getSlpitFunc(sep []byte, sep_len int) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		dataLen := len(data)

		if atEOF && dataLen == 0 {
			return 0, nil, nil
		}

		if i := bytes.Index(data, sep); i >= 0 {

			return i + sep_len, data[0:i], nil
		}

		if atEOF {
			return dataLen, data, bufio.ErrFinalToken
		}

		return 0, nil, nil
	}
}

type Decoder[T any] struct {
	decoder *gob.Decoder
	buffer  *bufio.Scanner
	entity  T
	sep     []byte
}

// NewDecoder - creats new gob decoder wrapper.
// separator - is hex encoded name-string of T
func NewDecoder[T any](r io.Reader, separator []byte) *Decoder[T] {
	buf := bufio.NewScanner(r)
	buf.Split(getSlpitFunc(separator, len(separator)))
	return &Decoder[T]{
		buffer: buf,
		sep:    separator,
	}
}

func (d *Decoder[T]) Decode(callback func(T)) error {
	var entity T
	var pre_separator []byte
	for d.buffer.Scan() {
		b := d.buffer.Bytes()
		if len(b) > 8 {

			encoded := append(pre_separator, bytes.TrimRight(b, string(pre_separator))...)
			dec := gob.NewDecoder(bytes.NewReader(encoded))
			for {
				if err := dec.Decode(entity); err != nil {
					if err == io.EOF {
						break
					} else {
						return err
					}

				}
				callback(entity)
			}
		} else {
			// Store pre_separator value
			pre_separator = make([]byte, len(b))
			copy(pre_separator, b)
			pre_separator = append(pre_separator, d.sep...)
		}
	}

	return nil
}
