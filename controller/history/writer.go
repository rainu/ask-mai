package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type Writer struct {
	historyPath string
}

func NewWriter(historyPath string) *Writer {
	return &Writer{
		historyPath: historyPath,
	}
}

func (w *Writer) openFile() (*os.File, error) {
	parentDir := path.Dir(w.historyPath)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		// Directory does not exist, create it
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create history file directory %s: %s", parentDir, err)
		}
	}

	f, err := os.OpenFile(w.historyPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open history file %s: %s", w.historyPath, err)
	}
	return f, nil
}

func (w *Writer) Write(entry Entry) error {
	if w.historyPath == "" {
		return nil
	}

	hf, err := w.openFile()
	if err != nil {
		return err
	}
	defer hf.Close()

	err = json.NewEncoder(hf).Encode(entry)
	if err != nil {
		return fmt.Errorf("failed to write history entry: %s", err)
	}

	return nil
}
