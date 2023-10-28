package infrastructure

import (
	"fmt"
	"os"
	"sync"
)

func NewConcurrentFileHandler(filePath string) ConcurrentFileHandler {
	var mu sync.Mutex
	return ConcurrentFileHandler{filePath, &mu}
}

type ConcurrentFileHandler struct {
	filePath string
	mu       *sync.Mutex
}

func (h *ConcurrentFileHandler) Read() ([]byte, os.FileInfo, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	empty := []byte{}

	fileContent, err := os.ReadFile(h.filePath)
	if err != nil {
		return empty, nil, err
	}

	fileStat, err := os.Stat(h.filePath)
	if err != nil {
		return empty, nil, err
	}

	return fileContent, fileStat, nil
}

func (h *ConcurrentFileHandler) Write(data []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	file, err := os.Create(h.filePath) // Create or open the file for writing (truncating if it already exists)
	if err != nil {
		return fmt.Errorf("error creating/opening file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to file '%s': %w", h.filePath, err)
	}

	return nil
}
