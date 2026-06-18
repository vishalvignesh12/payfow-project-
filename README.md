# Payflow

Payflow is a payment processing service built with Go. It handles transaction processing, fraud detection, and retry mechanisms for failed payments.

## Features

- Transaction creation with idempotency
- Automatic retry for failed transactions (with exponential backoff)
- Fraud scoring and flagging
- Audit logging for transactions
- Health check endpoint
- PostgreSQL for data storage
- Redis for caching / temporary storage
- Docker Compose for easy local development
- Database migrations using golang-migrate

## Tech Stack

- **Language**: Go 1.22+
- **Web Framework**: [chi](https://github.com/go-chi/chi)
- **Database**: PostgreSQL with [sqlx](https://github.com/jmoiron/sqlx)
- **Caching**: Redis with [go-redis/v9](https://github.com/redis/go-redis/v9)
- **Migration**: [golang-migrate/migrate/v4](https://github.com/golang-migrate/migrate/v4)
- **Environment Variables**: [godotenv](https://github.com/joho/godotenv)
- **Logging**: [zerolog](https://github.com/rs/zerolog)
- **Dev**: Docker Compose

## Project Structure

```
payflow/
├── cmd/
│   ├── server/          # HTTP server entrypoint
│   └── migrate/         # Database migration tool
├── internal/
│   ├── config/          # Configuration loading
│   ├── db/              # Database initialization (Postgres, Redis)
│   ├── models/          # Data models (transaction, audit, fraud)
│   ├── repository/      # Data access layer for models
│   └── statemachine/    # Transaction state machine logic
├── migrations/          # SQL migration files
├── docker/              # Dockerfiles
├── docker-compose.yml   # Docker Compose for local dev
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.22 or later
- Docker and Docker Compose (for local development)
- PostgreSQL and Redis (if not using Docker Compose)

### Environment Variables

The service configures itself via environment variables. Create a `.env` file in the project root or set them in your environment.

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | **required** |
| `REDIS_URL` | Redis connection string | **required** |
| `PORT` | Port for the HTTP server | `8080` |
| `WORKER_COUNT` | Number of workers for processing | `5` |
| `MAX_RETRIES` | Maximum retry attempts for failed transactions | `5` |
| `RETRY_BASE_DELAY_S` | Base delay (seconds) for retry backoff | `30` |
| `RATE_LIMIT_RPM` | Rate limit per minute for API endpoints | `100` |
| `BANK_FAILURE_RATE` | Simulated bank failure rate (0.0-1.0) | `0.2` |
| `HIGH_VALUE_THRESHOLD` | Amount (in smallest currency unit) considered high value | `10000000` |
| `ML_SERVICE_URL` | URL for external ML fraud service (optional) | *empty* |
| `ML_TIMEOUT_MS` | Timeout for ML service calls (milliseconds) | `100` |
| `VELOCITY_WINDOW_S` | Time window (seconds) for velocity checks | `60` |
| `VELOCITY_MAX_COUNT` | Max transactions allowed in velocity window | `5` |

### Running with Docker Compose

The easiest way to run Payflow locally is with Docker Compose:

```bash
docker-compose up
```

This will start:
- API server on `http://localhost:8080`
- PostgreSQL on `localhost:5433`
- Redis on `localhost:6380`

### Running Locally (without Docker)

1. Install dependencies:

```bash
go mod download
```

2. Set up PostgreSQL and Redis locally, then set the environment variables:
   - `DATABASE_URL` (e.g., `postgres://user:pass@localhost:5432/payflow?sslmode=disable`)
   - `REDIS_URL` (e.g., `redis://localhost:6379`)

3. Run database migrations:

```bash
go run ./cmd/migrate
```

4. Start the server:

```bash
go run ./cmd/server
```

## API Endpoints

Currently, Payflow exposes a single health check endpoint:

#### GET `/health`

Returns the health status of the service and its dependencies.

**Response:**

```json
{
  "status": "healthy",
  "db": "ok",
  "redis": "ok"
}
```

If the database is unreachable, returns `503 Service Unavailable` with:

```json
{
  "status": "unhealthy",
  "db": "error"
}
```

## Database Schema

Payflow uses the following tables (managed via migrations):

#### `transactions`
- `id`: UUID (primary key)
- `idempotency_key`: VARCHAR(255) (unique) - ensures duplicate requests are ignored
- `amount`: BIGINT - transaction amount in smallest currency unit (e.g., cents)
- `currency`: VARCHAR(3) - ISO currency code (default: `INR`)
- `status`: VARCHAR(20) - transaction status (`CREATED`, `PROCESSING`, `SUCCESS`, `FAILED`, etc.)
- `payer_id`: VARCHAR(255) - identifier of the payer
- `payee_id`: VARCHAR(255) - identifier of the payee
- `metadata`: JSONB - optional additional data
- `failure_reason`: TEXT - reason for failure if applicable
- `retry_count`: INT - number of retry attempts
- `next_retry_at`: TIMESTAMPTZ - when the next retry is scheduled
- `fraud_score`: FLOAT - score from fraud detection (0.0-1.0)
- `fraud_flagged`: BOOLEAN - whether transaction was flagged as fraud
- `permanent_failure`: BOOLEAN - whether transaction has failed permanently
- `created_at`: TIMESTAMPTZ - creation timestamp
- `updated_at`: TIMESTAMPTZ - last update timestamp

#### `audit_log`
Logs changes to transactions for compliance and debugging.

#### Fraud-related tables
Tables for storing fraud rules, models, and logs (see migration files).

### Models

The `internal/models` package contains Go structs that map to database tables:

- `Transactions`: Represents a payment transaction with fields matching the `transactions` table.
- `AuditLog`: Represents an audit entry for state changes.
- `Fraud`: Represents a fraud detection result linked to a transaction.

Each model includes struct tags for database mapping (`db`) and JSON serialization (`json`).

### Repository Layer

The `internal/repository` package provides data access functions for each model:

- **TransactionRepo**: CRUD operations for transactions, including creating, retrieving, updating status, updating retry counts, and fetching pending retries.
- **AuditRepo**: Writing audit logs and retrieving audit entries by transaction ID.
- **FraudRepo**: Writing fraud flags and retrieving recent fraud records.

Repositories use `sqlx` for database interactions and accept optional transaction (`*sqlx.Tx`) parameters for atomic operations.

### State Machine

The `internal/statemachine` package defines valid state transitions for a transaction:

- Created → Pending
- Pending → Processing
- Processing → Success or Failed
- Failed → Processing (retry)

The `IsValid` and `ValidTransition` functions enforce these rules, preventing invalid state changes.

## Migrations

Database migrations are managed with [golang-migrate](https://github.com/golang-migrate/migrate). The migration files are in the `migrations/` directory.

To apply migrations:

```bash
go run ./cmd/migrate
```

To rollback migrations (not implemented in the current migrate command, but can be done manually with the migrate CLI).

## Development

### Running Tests

There are currently no tests in the repository. Consider adding unit and integration tests.


--- 

*Note: This README was generated based on the source code. For any discrepancies, the source code is the ultimate source of truth.*