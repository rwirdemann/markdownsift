package os

import (
	"io"
	"os"
)

type Writer struct {
	writer io.Writer
}

func NewWriter() *Writer {
	return &Writer{writer: os.Stdout}
}

func (w *Writer) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

func (w *Writer) Create(name string) error {
	// nothing to do
	return nil
}

func (w *Writer) Close() error {
	// nothing to do
	return nil
}
