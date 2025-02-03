# Purchase Processor Microservice
## Table of Contents
- [Overview]()
- [Features]()
- [Technology Stack]()
- [Setup Instructions]()
- [Environment Variables]()
- [Build and Run]()
- [API Usage]()
- [Tracing and Logging]()
- [Kafka Consumer]()
- [Development Notes]()

## Overview
The **Purchase Processor Microservice** is a backend service responsible for handling user purchase transactions. It communicates with an authentication service over gRPC, processes event-based data streams via Kafka, and exposes a tracing system for monitoring.
This service is built with **Go** and adheres to distributed systems design principles. It includes robust logging, tracing, and external service communication.
## Features
- Secure gRPC communication for interaction with the `auth-service`.
- Kafka consumer for processing purchase events.
- HTTP API for managing purchases and exposing service health.
- Integrated distributed tracing with **OpenTelemetry**.
- Structured logging via **zap** logger.
- Environment variable-based configuration for flexibility.

## Technology Stack
- **Programming Language**: Go
- **gRPC**: For communication with the `auth-service`
- **Kafka (Sarama)**: For consuming and processing event streams
- **OpenTelemetry (OTLP)**: For distributed tracing
- **Zap**: For structured logging
- **HTTP Server**: Exposed on port `8084` for API requests

## Setup Instructions
### Prerequisites
1. Install [Go 1.23]() or higher.
2. Have a running Kafka broker (e.g., `kafka:9092`) for consuming purchase events.
3. Ensure the `auth-service` gRPC server is running and accessible.
4. Install necessary dependencies by running:
``` bash
   go mod tidy
```
## Environment Variables
The service uses the following environment variables:

| Environment Variable | Description | Default Value |
| --- | --- | --- |
| `KAFKA_BROKERS` | Kafka broker addresses, comma-separated | `kafka:9092` |
| `KAFKA_TOPIC` | The Kafka topic to consume messages from | `purchases` |
| `KAFKA_GROUP` | Kafka consumer group for this service | `purchase-processor-group` |
To set these variables locally, you can create an `.env` file with the following content:
``` env
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC=purchases
KAFKA_GROUP=purchase-processor-group
```
## Build and Run
### Build the Project
To build the Go binary, run:
``` bash
go build -o purchase-processor
```
### Run the Service
Run the binary with:
``` bash
./purchase-processor
```
## API Usage
The service exposes an HTTP API on port `8084`. Below is an example API request to process a purchase:
### Example Request (via `wget`)
``` bash
wget http://localhost:8084/api/v1/purchases \
     --header="Content-Type: application/json" \
     --header="Authorization: Bearer <your_oauth2_token>" \
     --post-data='{
       "user_id": "4f8e3c79-9a0f-4a02-823b-9c6a21f46ad8",
       "username": "first_user",
       "amount": 100.50,
       "currency": "USD",
       "payment_method": "credit_card",
       "card_holder_name": "John Doe",
       "card_number": "4111111111111111",
       "card_expiry": "12/25",
       "card_cvc": "123",
       "billing_address": "123 Elm Street, Springfield",
       "transaction_date": "2025-01-19T15:04:05Z"
     }' -O /dev/null
```
## Tracing and Logging
### Tracing
The service is integrated with OpenTelemetry and exports traces:
- Traces are sent to an **OTLP Collector** at the default endpoint `otel-collector:4317`.
- You can adjust the tracing parameters in the `tracing.InitTracer` function.

### Logging
The service uses `zap` for structured logging. Log outputs include details about errors, warnings, and runtime activity.
## Kafka Consumer
The service consumes events from a Kafka topic using the `sarama` library.
### Kafka Consumer Overview
- **Brokers**: Configured in the `KAFKA_BROKERS` environment variable.
- **Topic**: Consumes events from the `KAFKA_TOPIC` topic.
- **Group**: Part of the `KAFKA_GROUP` consumer group for load balancing.

### Kafka Configuration
The config uses the `RoundRobin` balancing strategy, and offsets are managed via Kafka.
### Example Event Payload
Consumer processes purchase events with a payload similar to:
``` json
{
  "user_id": "4f8e3c79-9a0f-4a02-823b-9c6a21f46ad8",
  "username": "first_user",
  "amount": 100.50,
  "currency": "USD",
  "payment_method": "credit_card",
  "transaction_date": "2025-01-19T15:04:05Z"
}
```
Any errors during processing are logged using the `zap` logging system.
## Development Notes
### Project File Structure
``` 
.
├── main.go                    # Entry point of the application
├── logging/                   # Logging utilities
├── tracing/                   # OpenTelemetry tracing setup
├── middleware/                # gRPC interceptors (e.g., for tracing)
├── pb/                        # gRPC protobuf definitions
├── service/                   # Core business logic
├── purchase_handler/          # Kafka consumer group handler
├── go.mod                     # Dependency management
└── README.md                  # Project documentation
```

----

# Useful links:

Jaeger: http://localhost:16686/search

OTEL getting started
https://opentelemetry.io/blog/2024/getting-started-with-otelsql/

OTEL SQL
https://github.com/XSAM/otelsql/blob/main/example/otel-collector/otel-collector.yaml

OTEL Collector
https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor/memorylimiterprocessor

GO Grpc
https://github.com/grpc/grpc-go/blob/master/examples/features/reflection/server/main.go

Go Reflection
https://grpc.io/docs/guides/reflection/

kafka
https://dev.to/deeshath/apache-kafka-kraft-mode-setup-5nj

debezium
https://debezium.io/documentation/reference/stable/connectors/mongodb.html
https://github.com/debezium/debezium-examples/blob/main/unwrap-mongodb-smt/docker-compose.yaml