# Implementation Approach and Design Considerations

## Overview

The implementation is a high-throughput REST service in Go designed to handle at least 10K requests per second. The service provides a `/api/verve/accept` endpoint that accepts an integer ID as a mandatory parameter and an optional HTTP endpoint for callbacks.

## Architecture

The system is structured in several packages to maintain separation of concerns:

1. **main** - Application entry point, initialization, and signal handling
2. **config** - Configuration management, loading from environment variables
3. **server** - HTTP server implementation and request handling
4. **stats** - Statistics collection, processing, and reporting

## Design Decisions

### Performance Considerations

- **Efficient HTTP handling**: Using Go's standard library HTTP server with tuned timeouts
- **Connection pooling**: Redis connections are pooled for optimal performance
- **Asynchronous processing**: Statistics processing happens in a separate goroutine
- **Minimal memory allocations**: Careful string and struct usage to minimize GC pressure
- **Context usage**: All external calls use contexts with appropriate timeouts

### Scalability

- **Stateless application**: The service is designed to be horizontally scalable
- **Redis for deduplication**: Using Redis SET data structure for efficient ID deduplication
- **Distributed state**: All state is maintained in Redis, allowing multiple instances to run behind a load balancer

### Extensibility

 Adding new stats writers is easy since they are abstacted using StatsWriter interface. Currently there are two implementations kafka writer and file writer implmenting this interface. Note :- You can set `KAFKA_ENABLED` to false in docker-compose.yaml to enable file writer.

### Extension 1: Handling Duplicate IDs Behind Load Balancer

To ensure ID deduplication works when the service is behind a load balancer:

- Using Redis's SADD command which guarantees atomic operations
- Each unique ID is added to a Redis SET with a key based on the current minute
- If multiple instances receive the same ID simultaneously, only one will succeed in adding it to the set
- All instances read from the same Redis instance, ensuring consistent counts

### Extension 2: Distributed Streaming

For sending statistics to a distributed streaming service:

- Implemented Kafka integration for publishing statistics
- Each minute's unique ID count is published to a Kafka topic
- Messages are keyed by timestamp for proper partitioning
- Configurable through environment variables

### Error Handling and Resilience

- Comprehensive error handling throughout the codebase
- Graceful shutdown handling for clean termination
- Connection retry mechanisms for Redis and Kafka
- Timeout handling for all external operations

### Configuration and Deployment

- Environment variable-based configuration
- Docker support for containerized deployment
- Reasonable defaults for quick startup
- Configurable timeouts and intervals

## Testing Considerations

While not implemented in this submission, the following testing approaches would be recommended:

- Unit tests for each package
- Integration tests for Redis and Kafka interactions
- Load testing to verify 10K requests/second throughput
- Chaos testing to verify behavior during network partitions

## Production Readiness

Additional considerations for a production environment:

- Metrics collection (Prometheus)
- Distributed tracing (OpenTelemetry)
- Health check endpoints
- Rate limiting
- Circuit breakers for external dependencies
- More comprehensive logging and monitoring