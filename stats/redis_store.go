package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"verve/config"
)

// RedisIDStore implements IDStore using Redis
type RedisIDStore struct {
	client *redis.Client
}

// NewRedisIDStore creates a new Redis-based ID store
func NewRedisIDStore(cfg *config.Config) (*RedisIDStore, error) {
	// Initialize Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &RedisIDStore{
		client: client,
	}, nil
}

// AddID adds an ID to the Redis set for a specific minute
func (s *RedisIDStore) AddID(ctx context.Context, minuteTimestamp time.Time, id int64) error {
	minuteKey := fmt.Sprintf("verve:minute:%d", minuteTimestamp.Unix())

	// Add id to the set of unique ids for this minute
	_, err := s.client.SAdd(ctx, minuteKey, id).Result()
	if err != nil {
		return fmt.Errorf("failed to add id to Redis set: %w", err)
	}

	// Set expiration on the key if it's new
	_, err = s.client.Expire(ctx, minuteKey, 2*time.Hour).Result() // Keep it for a reasonable time
	if err != nil {
		return fmt.Errorf("failed to set expiration on key: %w", err)
	}

	return nil
}

// CountIDs counts unique IDs for a specific minute
func (s *RedisIDStore) CountIDs(ctx context.Context, minuteTimestamp time.Time) (int64, error) {
	minuteKey := fmt.Sprintf("verve:minute:%d", minuteTimestamp.Unix())

	// Get count of unique IDs
	count, err := s.client.SCard(ctx, minuteKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get unique count from Redis: %w", err)
	}

	return count, nil
}

// StoreEndpoint stores an endpoint for a specific minute
func (s *RedisIDStore) StoreEndpoint(ctx context.Context, minuteTimestamp time.Time, endpoint string) error {
	endpointKey := fmt.Sprintf("verve:endpoint:%d", minuteTimestamp.Unix())

	// Store the endpoint
	err := s.client.Set(ctx, endpointKey, endpoint, 2*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store endpoint: %w", err)
	}

	return nil
}

// GetEndpoint retrieves an endpoint for a specific minute
func (s *RedisIDStore) GetEndpoint(ctx context.Context, minuteTimestamp time.Time) (string, error) {
	endpointKey := fmt.Sprintf("verve:endpoint:%d", minuteTimestamp.Unix())

	// Get the endpoint
	endpoint, err := s.client.Get(ctx, endpointKey).Result()
	if err == redis.Nil {
		return "", nil // No endpoint stored
	} else if err != nil {
		return "", fmt.Errorf("failed to get endpoint from Redis: %w", err)
	}

	return endpoint, nil
}

// Close cleans up resources
func (s *RedisIDStore) Close() error {
	return s.client.Close()
}
