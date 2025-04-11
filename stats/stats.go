package stats

import (
	"context"
	"time"
)

// StatsWriter defines the interface for writing stats
type StatsWriter interface {
	// WriteStats writes stats for a given timestamp with unique count
	WriteStats(ctx context.Context, timestamp time.Time, count int64) error

	// Close cleans up resources
	Close() error
}

// Collector handles the collection and processing of request statistics
type Collector struct {
	idStore       IDStore
	writer        StatsWriter
	logger        Logger
	flushInterval time.Duration
	httpClient    HTTPClient
}

// IDStore defines the interface for storing and retrieving unique IDs
type IDStore interface {
	// AddID adds an ID to the store for a specific minute
	AddID(ctx context.Context, minuteTimestamp time.Time, id int64) error

	// CountIDs counts unique IDs for a specific minute
	CountIDs(ctx context.Context, minuteTimestamp time.Time) (int64, error)

	// Close cleans up resources
	Close() error
}

// Logger defines the interface for logging
type Logger interface {
	Printf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}
