package stats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"verve/config"
)

// KafkaWriter implements StatsWriter using Kafka
type KafkaWriter struct {
	writer *kafka.Writer
	logger Logger
	topic  string
}

// NewKafkaWriter creates a new Kafka-based stats writer
func NewKafkaWriter(cfg *config.Config, logger Logger) (*KafkaWriter, error) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Kafka.Brokers...),
		Topic:                  cfg.Kafka.Topic,
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: true,
	}

	return &KafkaWriter{
		writer: writer,
		logger: logger,
		topic:  cfg.Kafka.Topic,
	}, nil
}

// WriteStats writes stats to Kafka
func (w *KafkaWriter) WriteStats(ctx context.Context, timestamp time.Time, count int64) error {
	stats := map[string]interface{}{
		"timestamp": timestamp.Format(time.RFC3339),
		"count":     count,
	}

	value, err := json.Marshal(stats)
	if err != nil {
		return fmt.Errorf("failed to marshal stats for Kafka: %w", err)
	}

	err = w.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(fmt.Sprintf("%d", timestamp.Unix())),
			Value: value,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	w.logger.Printf("Sent stats to Kafka topic %s", w.topic)
	return nil
}

// Close cleans up resources
func (w *KafkaWriter) Close() error {
	return w.writer.Close()
}
