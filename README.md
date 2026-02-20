# Subscription Service

REST API service for aggregating user subscription data.

## Features

- CRUD operations for subscriptions
- Filter subscriptions by user, service, and date
- Calculate total cost for a period
- PostgreSQL database with migrations
- Docker support
- Swagger documentation
- Unit and integration tests

## Tech Stack

- Go 1.21
- Gin Web Framework
- PostgreSQL
- Docker & Docker Compose
- Swagger for API documentation
- Goose for migrations

## API Endpoints

### Subscriptions

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/subscriptions` | Create a new subscription |
| GET | `/api/v1/subscriptions` | List subscriptions with filters |
| GET | `/api/v1/subscriptions/total` | Get total price for period |
| GET | `/api/v1/subscriptions/:id` | Get subscription by ID |
| PUT | `/api/v1/subscriptions/:id` | Update subscription |
| DELETE | `/api/v1/subscriptions/:id` | Delete subscription |

### Swagger Documentation

After starting the service, visit:
http://localhost:8080/swagger/index.html

## Quick Start

### Using Docker (recommended)

```bash
# Clone repository
git clone https://github.com/t5129001t-jpg/subscription-service.git
cd subscription-service

# Start with Docker Compose
docker-compose up -d

# Service will be available at http://localhost:8080

# Copy environment file
cp .env.example .env

# Install dependencies
go mod download

# Run the service
go run cmd/main.go

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# View coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | 8080 |
| `SERVER_READ_TIMEOUT` | Read timeout (seconds) | 10 |
| `SERVER_WRITE_TIMEOUT` | Write timeout (seconds) | 10 |
| `DB_HOST` | PostgreSQL host | localhost |
| `DB_PORT` | PostgreSQL port | 5432 |
| `DB_USER` | Database user | postgres |
| `DB_PASSWORD` | Database password | postgres |
| `DB_NAME` | Database name | subscription_db |
| `DB_SSLMODE` | SSL mode | disable |

.
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/                  # Configuration
│   ├── handler/                  # HTTP handlers
│   ├── model/                    # Data models
│   ├── repository/               # Database operations
│   └── service/                  # Business logic
├── migrations/                   # Database migrations
├── tests/                        # Integration tests
│   ├── handler/                   # Handler tests
│   ├── repository/                # Repository tests
│   └── service/                   # Service tests
├── docs/                         # Swagger documentation
├── docker-compose.yaml           # Docker composition
├── Dockerfile                    # Docker build file
└── .env.example                  # Environment variables example

##API Examples

#Create a subscription

curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Netflix",
    "price": 1000,
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "start_date": "01-2024"
  }'

#Get total price for period

curl "http://localhost:8080/api/v1/subscriptions/total?user_id=123e4567-e89b-12d3-a456-426614174000&start_month=01-2024&end_month=12-2024"

##License

#MIT


## **Как вставить:**

1. **Откройте файл в nano:**
   ```bash
   nano README.md

