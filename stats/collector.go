package stats

import (
	"context"
	"fmt"
	"time"

	"verve/config"
)

// NewCollector creates a new stats collector
func NewCollector(cfg *config.Config, logger Logger) (*Collector, error) {
	// Initialize ID store (Redis)
	idStore, err := NewRedisIDStore(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ID store: %w", err)
	}

	// Initialize stats writer based on configuration
	var writer StatsWriter
	if cfg.Kafka.Enabled {
		kafkaWriter, err := NewKafkaWriter(cfg, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Kafka writer: %w", err)
		}
		writer = kafkaWriter
	} else {
		writer, err = NewFileWriter(cfg.LogFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize File writer: %w", err)
		}
	}

	return &Collector{
		idStore:       idStore,
		writer:        writer,
		logger:        logger,
		flushInterval: cfg.Stats.FlushInterval,
		httpClient:    NewStandardHTTPClient(),
	}, nil
}

// ProcessRequest processes a request with the given id
func (c *Collector) ProcessRequest(id int64, endpoint string) error {
	// Get current minute timestamp (truncated to minute)
	now := time.Now().Truncate(time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Store the endpoint for later use (if provided)
	if endpoint != "" {
		err := c.idStore.StoreEndpoint(ctx, now, endpoint)
		if err != nil {
			return fmt.Errorf("failed to store endpoint: %w", err)
		}
	}

	// Add id to the set of unique ids for this minute
	err := c.idStore.AddID(ctx, now, id)
	if err != nil {
		return fmt.Errorf("failed to add id: %w", err)
	}

	return nil
}

// Run starts the stats collection and processing loop
func (c *Collector) Run() {
	ticker := time.NewTicker(c.flushInterval)
	defer ticker.Stop()

	for range ticker.C {
		c.processMinuteStats()
	}
}

// processMinuteStats processes the stats for the previous minute
func (c *Collector) processMinuteStats() {
	// Get the previous minute timestamp
	prevMinute := time.Now().Add(-c.flushInterval).Truncate(time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get count of unique IDs
	uniqueCount, err := c.idStore.CountIDs(ctx, prevMinute)
	if err != nil {
		c.logger.Printf("Failed to get unique count: %v", err)
		return
	}

	// Write stats to the configured writer
	err = c.writer.WriteStats(ctx, prevMinute, uniqueCount)
	if err != nil {
		c.logger.Printf("Failed to write stats: %v", err)
	}

	// Send HTTP POST if endpoint was provided
	c.sendHTTPNotification(prevMinute, uniqueCount)
}

// sendHTTPNotification sends stats via HTTP POST if an endpoint was provided
func (c *Collector) sendHTTPNotification(timestamp time.Time, count int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to get the endpoint for this minute
	endpoint, err := c.idStore.GetEndpoint(ctx, timestamp)
	if err != nil {
		c.logger.Printf("Failed to get endpoint: %v", err)
		return
	}

	if endpoint == "" {
		// No endpoint stored for this minute
		return
	}

	// Add count as a query parameter
	url := fmt.Sprintf("%s?count=%d", endpoint, count)

	// Send HTTP POST request
	resp, err := c.httpClient.Post(url, "text/plain", nil)
	if err != nil {
		c.logger.Printf("Failed to send HTTP POST to endpoint %s: %v", endpoint, err)
		return
	}
	defer resp.Close()

	// Log the HTTP status code
	c.logger.Printf("HTTP POST to endpoint %s returned status: %d", endpoint, resp.StatusCode())
}

// Close cleans up resources
func (c *Collector) Close() {
	if c.idStore != nil {
		_ = c.idStore.Close()
	}
	if c.writer != nil {
		_ = c.writer.Close()
	}
}
