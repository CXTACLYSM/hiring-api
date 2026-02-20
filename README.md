# Hiring API

Go microservices backend with CQRS architecture, event-driven communication, and streaming replication.

## Architecture

```
┌──────────────┐  REST   ┌──────────────┐  REST
│    Client    │────────▶│    Client    │────────▶
└──────┬───────┘         └──────┬───────┘
       │                        │
       ▼                        ▼
┌──────────────┐  gRPC   ┌──────────────┐
│     Auth     │◀────────│     Blog     │
│  :8080/:50051│         │    :8081     │
└──────┬───────┘         └──────┬───────┘
       │ publish                │ publish
       ▼                        ▼
┌──────────────────────────────────────────┐
│              Apache Kafka                │
│         (user.created, post.*)           │
└──────────────────┬───────────────────────┘
                   │ consume
                   ▼
            ┌──────────────┐
            │ Notifications│
            │    :8082     │
            └──────────────┘

┌──────────────┐  WAL   ┌──────────────┐
│  PostgreSQL  │───────▶│  PostgreSQL  │
│   Primary    │ stream │   Replica    │
│   (write)    │        │   (read)     │
└──────────────┘        └──────────────┘

┌──────────────┐
│    Redis     │  token cache, sessions
└──────────────┘
```

### Services

| Service | Port | Role |
|---------|------|------|
| **Auth** | 8080 (HTTP), 50051 (gRPC) | JWT authentication, user management |
| **Blog** | 8081 | CRUD posts, authenticates via gRPC to Auth |
| **Notifications** | 8082 | Kafka consumer, processes events |

### Infrastructure

| Component | Purpose |
|-----------|---------|
| **PostgreSQL Primary** | Write operations (Commands) |
| **PostgreSQL Replica** | Read operations (Queries), streaming replication |
| **Redis** | Token cache, session store |
| **Apache Kafka** | Async event-driven communication (KRaft, no Zookeeper) |

## Quick Start

```bash
# Start everything (build + infrastructure + services)
make up

# Check status
make status

# Stop
make down
```

## Development

```bash
# Rebuild and restart a single service
make run-auth
make run-blog
make run-notifications

# Run migrations
make migrate-up path=auth
make migrate-up path=blog

# View logs
make logs-auth
make logs-blog
make logs-notifications
make logs-kafka
```

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/go-chi/chi/v5` | HTTP router |
| `github.com/jackc/pgx/v5` | PostgreSQL driver |
| `github.com/gomodule/redigo` | Redis client |
| `github.com/IBM/sarama` | Kafka client |
| `github.com/spf13/viper` | Configuration |
| `github.com/golang-jwt/jwt/v5` | JWT authentication |
| `github.com/google/uuid` | UUID generation |
| `github.com/go-playground/validator/v10` | Request validation |
| `github.com/golang-migrate/migrate` | Database migrations |
| `google.golang.org/grpc` | gRPC framework |
| `google.golang.org/protobuf` | Protocol Buffers |

## Project Structure

```
cmd/
  auth/                 # Auth service entrypoint
  blog/                 # Blog service entrypoint
  notifications/        # Notifications service entrypoint
configs/
  auth/                 # Auth service config (app, postgres, redis, kafka)
  blog/                 # Blog service config
  notifications/        # Notifications service config
internal/
  auth/
    user/
      application/      # Commands, Queries, Services, DTOs
      domain/           # Entities
      infrastructure/   # Handlers, DB implementations
    di/                 # Dependency injection container
    grpc/               # gRPC server (ValidateToken)
  blog/
    post/               # Same DDD structure as auth/user
    di/
  notifications/
    di/
    consumer/           # Kafka event handlers
pkg/
  postgres/             # PostgreSQL connector (read/write pools)
  redis/                # Redis connector (auth/resource pools)
  kafka/                # Kafka producer/consumer
  events/               # Shared event definitions
  grpc/auth/v1/         # Generated protobuf/gRPC code
  shared/               # Shared utilities (handlers, middlewares, validation)
migrations/
  postgres/
    auth/               # Auth migrations
    blog/               # Blog migrations
builds/
  Dockerfile            # Multi-stage build (auth, blog, notifications)
  local/                # Local config files (postgres, redis)
scripts/
  init/                 # DB init scripts, Kafka topic creation
.envs/
  app/
    services/           # Generated app service env files
    templates/          # App env templates
  infra/
    services/           # Generated infra service env files
    templates/          # Infra env templates
```

## Useful Commands

```bash
# Infrastructure only (without app services)
make infra-up
make infra-down

# PostgreSQL
make pg-primary             # Shell into primary
make pg-replica             # Shell into replica
make pg-replication-status  # Check replication
make pg-replication-lag     # Check lag
make pg-write-monitor       # Monitor active queries (live)
make pg-write-log           # Tail query log

# Redis
make redis-cli

# Kafka
make kafka-topics
make kafka-create-topic name=user.created partitions=3
make kafka-consumer-groups

# System
make audit                  # Full system audit
make clean                  # Remove all artifacts
make down-v                 # Stop + remove volumes (DESTRUCTIVE)
```