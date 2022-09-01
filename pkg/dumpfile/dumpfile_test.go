package dumpfile

import (
	"bufio"
	"io"
	"testing"
)

func TestDumpFile(t *testing.T) {
	path := "test.db"
	df, err := NewDumpFile(path)
	if err != nil {
		t.Error(err)
	}

	buffR := bufio.NewReader(df)
	df.Write([]byte("hello \n"))

	res, err := buffR.ReadString(byte('\n'))
	if err != nil {
		if err != io.EOF {
			t.Error(err)
		}
	}

	t.Logf("In file: %s\n", res)
	df.Close()
	// if err := os.Remove(path); err != nil {
	// 	t.Error(err)
	// }
}
