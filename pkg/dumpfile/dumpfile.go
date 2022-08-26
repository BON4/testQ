package dumpfile

import (
	"io"
	"os"
	"path/filepath"
)

const DEFAULT_DUMP_NAME = ".temp.db"

type DumpFile struct {
	temp   bool
	path   string
	reader *os.File
	writer *os.File
}

// NewDumpFile - Creates Dump File, if file name is not specified it's creating temp file
func NewDumpFile(path string) (io.ReadWriteCloser, error) {
	df := &DumpFile{}
	dir, fname := filepath.Split(path)
	if len(fname) == 0 {
		df.temp = true
		df.path = dir + DEFAULT_DUMP_NAME

		var err error
		df.reader, err = os.CreateTemp("", df.path)
		if err != nil {
			return nil, err
		}

		df.writer = df.reader

	} else {
		df.path = path

		var err error

		df.writer, err = os.OpenFile(df.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			return nil, err
		}

		df.reader, err = os.OpenFile(df.path, os.O_RDONLY, 0644)
		if err != nil {
			return nil, err
		}

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
	if df.temp {
		if err := df.reader.Close(); err != nil {
			return err
		}

		if err := df.writer.Close(); err != nil {
			return err
		}

		if err := os.Remove(df.reader.Name()); err != nil {
			return err
		}
	} else {
		if err := df.reader.Close(); err != nil {
			return err
		}

		if err := df.writer.Close(); err != nil {
			return err
		}
	}
	return nil
}
