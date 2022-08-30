package dumpfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const DEFAULT_DUMP_NAME = ".temp.db"

type DumpFile struct {
	path   string
	reader *os.File
	writer *os.File
}

// NewDumpFile - Creates Dump File, if file name is not specified it's creating file with default name
func NewDumpFile(path string) (io.ReadWriteCloser, error) {
	df := &DumpFile{}
	dir, fname := filepath.Split(path)
	if len(fname) == 0 {
		df.path = dir + DEFAULT_DUMP_NAME
	} else {
		df.path = path
	}

	var err error

	df.writer, err = os.OpenFile(df.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	df.reader, err = os.OpenFile(df.path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	return df, nil
}

func (df *DumpFile) Read(b []byte) (n int, err error) {
	return df.reader.Read(b)
}

func (df *DumpFile) Write(b []byte) (n int, err error) {
	return df.writer.Write(b)
}

func (df *DumpFile) Close() error {
	if err := df.reader.Close(); err != nil {
		return err
	}

	if err := df.writer.Close(); err != nil {
		return err
	}

	return nil
}
