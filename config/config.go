package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	Server struct {
		Addr         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
	Kafka struct {
		Enabled bool
		Brokers []string
		Topic   string
	}
	LogFilePath string
	Stats       struct {
		FlushInterval time.Duration
	}
}

// Load loads configuration from environment variables with defaults
func Load() (*Config, error) {
	cfg := &Config{}

	// Server configuration
	cfg.Server.Addr = getEnv("SERVER_ADDR", ":8080")
	cfg.Server.ReadTimeout = getDurationEnv("SERVER_READ_TIMEOUT", 5*time.Second)
	cfg.Server.WriteTimeout = getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second)

	// Redis configuration
	cfg.Redis.Addr = getEnv("REDIS_ADDR", "localhost:6379")
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	cfg.Redis.DB = getIntEnv("REDIS_DB", 0)

	// Kafka configuration
	cfg.Kafka.Enabled = getBoolEnv("KAFKA_ENABLED", false)
	cfg.Kafka.Brokers = []string{getEnv("KAFKA_BROKER", "localhost:9092")}
	cfg.Kafka.Topic = getEnv("KAFKA_TOPIC", "verve-stats")

	// LOG file path
	cfg.LogFilePath = getEnv("LOG_FILE_PATH", "/app/host/stats.log")
	// Stats configuration
	cfg.Stats.FlushInterval = getDurationEnv("STATS_FLUSH_INTERVAL", 60*time.Second)

	return cfg, nil
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if durationValue, err := time.ParseDuration(value); err == nil {
			return durationValue
		}
	}
	return defaultValue
}
