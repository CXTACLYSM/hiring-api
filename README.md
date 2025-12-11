# Hiring API

Go backend service with CQRS architecture.

## Architecture

```
┌─────────┐     ┌─────────┐
│   API   │     │ Worker  │
└────┬────┘     └────┬────┘
     │ write         │ consume
     ▼               ▼
┌─────────┐     ┌──────────┐
│ Primary │     │ RabbitMQ │
│ (PG)    │     └──────────┘
└────┬────┘
     │ WAL stream
     ▼
┌─────────┐     ┌─────────┐
│ Replica │     │  Redis  │
│ (PG)    │     │ (cache) │
└─────────┘     └─────────┘
```

- **API** — HTTP server, writes to Primary, reads from Replica
- **Worker** — background jobs consumer (RabbitMQ)
- **PostgreSQL Primary** — write operations (Commands)
- **PostgreSQL Replica** — read operations (Queries), streaming replication
- **Redis** — caching, sessions
- **RabbitMQ** — async message queue

## Quick Start

```bash
# 1. Clone and setup env files
cp .env.app.example .env.app
cp .env.infra.example .env.infra

# 2. Start all services
make up

# 3. Check replication status
make pg-replication-status

# 4. Run migrations
make migrate-up
```

## Useful Commands

| Command | Description |
|---------|-------------|
| `make up` | Start all services |
| `make down` | Stop services |
| `make down-v` | Stop + remove volumes |
| `make logs` | Show logs |
| `make pg-primary` | Connect to primary DB |
| `make pg-replica` | Connect to replica DB |
| `make pg-replication-status` | Check replication |
| `make migrate-up` | Run migrations |
| `make migrate-create name=xxx` | Create migration |

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/jackc/pgx/v5` | PostgreSQL driver |
| `github.com/rabbitmq/amqp091-go` | RabbitMQ client |
| `github.com/spf13/viper` | Configuration |
| `github.com/golang-jwt/jwt/v5` | JWT authentication |
| `github.com/google/uuid` | UUID generation |
| `github.com/ClickHouse/clickhouse-go/v2` | ClickHouse driver (analytics) |

## Project Structure

```
cmd/
  api/          # HTTP server entrypoint
  worker/       # Background worker entrypoint
internal/
  command/      # CQRS commands (write)
  query/        # CQRS queries (read)
  domain/       # Business entities
  handler/      # HTTP handlers
  service/      # Business logic
pkg/            # Shared packages (connectors)
migrations/     # SQL migrations
builds/         # Dockerfiles, configs
scripts/        # Init scripts
```
