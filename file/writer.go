package file

import (
	"os"
	"path/filepath"
)

type Writer struct {
	path string
	file *os.File
}

func NewWriter(path string) (*Writer, error) {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	writer := &Writer{path: path}
	return writer, nil
}

func (w *Writer) Write(data []byte) (int, error) {
	return w.file.Write(data)
}

func (w *Writer) Create(name string) error {
	path := filepath.Join(w.path, name[1:]+".md")
	var err error
	w.file, err = os.Create(path)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) Close() error {
	return nil
}
