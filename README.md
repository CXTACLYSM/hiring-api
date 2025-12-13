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
# 1. Setup environment
make check-env
make env
make env-infra

# 2. Start infrastructure
make infra-up

# 3. Check replication
make pg-replication-status
make pg-test-replication

# 4. Build and run services
make docker-build
make run

# 5. Run migrations
make migrate-up
```

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
  api/              # HTTP server entrypoint
  worker/           # Background worker entrypoint
internal/
  command/          # CQRS commands (write)
  query/            # CQRS queries (read)
  domain/           # Business entities
  handler/          # HTTP handlers
  service/          # Business logic
pkg/                # Shared packages (connectors)
migrations/         # SQL migrations
builds/             # Dockerfiles, configs
scripts/            # Init scripts
.envs/
  app/
    services/       # Generated app service env files
    templates/      # App env templates
  infra/
    services/       # Generated infra service env files
    templates/      # Infra env templates
```

## Development Workflow
```bash
# Daily workflow
make check-env          # Verify environment
make infra-up           # Start infrastructure
make audit              # Check system status
make run                # Start API and Worker

# Debugging
make pg-write-log       # Watch queries on primary
make pg-read-log        # Watch queries on replica
make pg-write-monitor   # Monitor active connections
make logs               # View all service logs

# Database
make pg-primary         # Shell into primary DB
make migrate-up         # Apply migrations
make pg-replication-status  # Verify replication

# Cleanup
make down               # Stop services
make down-v             # Stop and remove all data
make clean              # Full cleanup
```