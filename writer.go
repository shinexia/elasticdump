package main

import (
	"io"
	"os"
)

type LazyWriter struct {
	filename string
	writer   io.WriteCloser
}

func newLazyWriter(filename string) *LazyWriter {
	w := &LazyWriter{
		filename: filename,
		writer:   nil,
	}
	return w
}

func (w *LazyWriter) Write(data []byte) (n int, err error) {
	if w.writer == nil {
		file, err := os.Create(w.filename)
		if err != nil {
			return 0, err
		}
		w.writer = file
	}
	return w.writer.Write(data)
}

func (w *LazyWriter) Close() error {
	if w.writer != nil {
		return w.writer.Close()
	}
	return nil
}
