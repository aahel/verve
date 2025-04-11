# Running the Application

## Prerequisites

- Go
- Docker and Docker Compose (for containerized deployment)
- Redis and Kafka (if running without Docker)

## Option 1: Using Docker Compose (Recommended)
This is the simplest way to run the entire stack including Redis and Kafka:

### Start all services using Docker Compose:
```bash
docker compose up --build
```

This will start:
- verve-service
- Redis
- Zookeeper
- Kafka

### Test the endpoint:
```bash
curl "http://localhost:8080/api/verve/accept?id=123&endpoint=https://webhook.site/b07a22a4-148c-4426-8f42-6c52b213a9b3"
```

## Option 2: Running Locally

1. Make sure Redis is running locally:
   ```bash
   docker run -d -p 6379:6379 redis:alpine
   ```

2. If you want to use Kafka, set it up (you can use Confluent's quickstart):
   ```bash
   docker run -d -p 2181:2181 -p 9092:9092 -e ADVERTISED_HOST=localhost confluentinc/cp-kafka:latest
   ```
   Also make sure to set `KAFKA_ENABLED` env var to `true`.
   ```bash
   export KAFKA_ENABLED=true
   ```
   
3. Install the dependencies:
   ```bash
   go mod tidy
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

5. Test the endpoint:
   ```bash
   curl "http://localhost:8080/api/verve/accept?id=123&endpoint=https://webhook.site/acca7e6e-d531-4f0a-81c3-dcd8a52ca04b"
   ```

## Environment Variables

You can customize the application behavior with these environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_ADDR` | HTTP server address | `:8080` |
| `SERVER_READ_TIMEOUT` | HTTP read timeout | `5s` |
| `SERVER_WRITE_TIMEOUT` | HTTP write timeout | `10s` |
| `REDIS_ADDR` | Redis address | `localhost:6379` |
| `REDIS_PASSWORD` | Redis password | _(empty)_ |
| `REDIS_DB` | Redis database number | `0` |
| `KAFKA_ENABLED` | Enable Kafka producer | `true` |
| `KAFKA_BROKER` | Kafka broker address | `localhost:9092` |
| `KAFKA_TOPIC` | Kafka topic name | `verve-stats` |
| `STATS_FLUSH_INTERVAL` | Interval for stats processing | `60s` |
| `LOG_FILE_PATH` | Path to log file.Will be used only if `KAFKA_ENABLED` is false. | `./stats.log` |

