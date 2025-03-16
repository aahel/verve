package stats

import (
	"context"
	"fmt"
	"os"
	"time"
)

// FileWriter implements StatsWriter using file output
type FileWriter struct {
	file *os.File
}

// NewFileWriter creates a new file-based stats writer
func NewFileWriter(filePath string) (*FileWriter, error) {
	// Open file with append mode, create if not exists
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open stats file: %w", err)
	}
	return &FileWriter{
		file: file,
	}, nil
}

// WriteStats writes stats to the file
func (w *FileWriter) WriteStats(ctx context.Context, timestamp time.Time, count int64) error {
	_, err := fmt.Fprintf(w.file, "Unique requests in minute %s: %d\n",
		timestamp.Format("15:04:05"), count)
	if err != nil {
		return fmt.Errorf("failed to write stats to file: %w", err)
	}
	return nil
}

// Close implements StatsWriter interface by closing the file
func (w *FileWriter) Close() error {
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}
