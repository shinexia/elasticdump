/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package dumpdata

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

var (
	newLine = []byte("\n")
)

type WriteDataFunc func(hits []json.RawMessage) (int, error)

type LazyDataWriter struct {
	outputFile string
	file       *os.File
}

func NewLazyDataWriter(outputFile string) *LazyDataWriter {
	return &LazyDataWriter{
		outputFile: outputFile,
	}
}

func (w *LazyDataWriter) Write(hits []json.RawMessage) (int, error) {
	if w.file == nil {
		f, err := os.Create(w.outputFile)
		if err != nil {
			return 0, errors.Cause(err)
		}
		w.file = f
	}
	count := 0
	for _, hit := range hits {
		_, err := w.file.Write(hit)
		if err != nil {
			return count, errors.Cause(err)
		}
		_, err = w.file.Write(newLine)
		if err != nil {
			return count, errors.Cause(err)
		}
		count += 1
	}
	return count, nil
}

func (w *LazyDataWriter) Close() error {
	if w.file != nil {
		f := w.file
		w.file = nil
		return f.Close()
	}
	return nil
}
