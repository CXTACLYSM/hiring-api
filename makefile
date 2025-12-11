# =============================================================================
# HIRING SERVICE - MAKEFILE
# =============================================================================

# Load environment files
ifneq (,$(wildcard ./.env.app))
    include .env.app
    export
endif
ifneq (,$(wildcard ./.env.infra))
    include .env.infra
    export
endif

# Colors for output
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RED    := \033[0;31m
NC     := \033[0m # No Color

# =============================================================================
# ENVIRONMENT MANAGEMENT
# =============================================================================

.PHONY: env
env: ## Generate .env from .env.app and .env.infra
	@echo "$(GREEN)Generating .env from .env.app and .env.infra...$(NC)"
	@cat /dev/null > .env
	@echo "# ==============================================================================" >> .env
	@echo "# AUTO-GENERATED - DO NOT EDIT" >> .env
	@echo "# Generated from .env.app and .env.infra" >> .env
	@echo "# Run: make env" >> .env
	@echo "# ==============================================================================" >> .env
	@echo "" >> .env
	@echo "# --- FROM .env.app ---" >> .env
	@grep -v '^\s*#' .env.app | grep -v '^\s*$$' >> .env || true
	@echo "" >> .env
	@echo "# --- FROM .env.infra ---" >> .env
	@grep -v '^\s*#' .env.infra | grep -v '^\s*$$' >> .env || true
	@echo "$(GREEN).env generated!$(NC)"

# =============================================================================
# DOCKER COMPOSE COMMANDS
# =============================================================================

.PHONY: up
up: check-env env ## Start all services
	@echo "$(GREEN)Starting all services...$(NC)"
	docker compose up -d
	@echo "$(GREEN)Services started!$(NC)"
	@$(MAKE) status

.PHONY: up-build
up-build: check-env env ## Start all services with rebuild
	@echo "$(GREEN)Building and starting all services...$(NC)"
	docker compose up -d --build
	@echo "$(GREEN)Services started!$(NC)"
	@$(MAKE) status

.PHONY: down
down: ## Stop all services
	@echo "$(YELLOW)Stopping all services...$(NC)"
	docker compose down
	@echo "$(GREEN)Services stopped!$(NC)"

.PHONY: down-v
down-v: ## Stop all services and remove volumes (DESTRUCTIVE!)
	@echo "$(RED)Stopping all services and removing volumes...$(NC)"
	@read -p "Are you sure? This will delete all data! [y/N] " confirm && [ "$$confirm" = "y" ]
	docker compose down -v
	@echo "$(GREEN)Services stopped and volumes removed!$(NC)"

.PHONY: restart
restart: down up ## Restart all services

.PHONY: logs
logs: ## Show logs for all services
	docker compose logs -f

.PHONY: logs-db
logs-db: ## Show PostgreSQL logs (primary and replica)
	docker compose logs -f postgres_primary postgres_replica

.PHONY: status
status: ## Show status of all services
	@echo ""
	@echo "$(GREEN)=== Services Status ===$(NC)"
	@docker compose ps
	@echo ""

# =============================================================================
# INFRASTRUCTURE ONLY (without app services)
# =============================================================================

.PHONY: infra-up
infra-up: check-env ## Start only infrastructure (postgres, redis, rabbitmq)
	@echo "$(GREEN)Starting infrastructure services...$(NC)"
	docker compose up -d postgres_primary postgres_replica redis rabbitmq
	@echo "$(GREEN)Infrastructure started!$(NC)"
	@$(MAKE) status

.PHONY: infra-down
infra-down: ## Stop infrastructure services
	@echo "$(YELLOW)Stopping infrastructure services...$(NC)"
	docker compose stop postgres_primary postgres_replica redis rabbitmq
	@echo "$(GREEN)Infrastructure stopped!$(NC)"

# =============================================================================
# POSTGRESQL REPLICATION
# =============================================================================

.PHONY: pg-replication-status
pg-replication-status: ## Check PostgreSQL replication status
	@echo "$(GREEN)=== Primary: Replication Status ===$(NC)"
	@docker exec -it $(POSTGRES_PRIMARY_CONTAINER_NAME) psql -U $(POSTGRES_SUPERUSER) -d $(POSTGRES_WRITE_DB) -c \
		"SELECT client_addr, state, sent_lsn, write_lsn, flush_lsn, replay_lsn, sync_state FROM pg_stat_replication;"
	@echo ""
	@echo "$(GREEN)=== Primary: Replication Slots ===$(NC)"
	@docker exec -it $(POSTGRES_PRIMARY_CONTAINER_NAME) psql -U $(POSTGRES_SUPERUSER) -d $(POSTGRES_WRITE_DB) -c \
		"SELECT slot_name, slot_type, active, restart_lsn FROM pg_replication_slots;"
	@echo ""
	@echo "$(GREEN)=== Replica: Recovery Status ===$(NC)"
	@docker exec -it $(POSTGRES_REPLICA_CONTAINER_NAME) psql -U $(POSTGRES_SUPERUSER) -d $(POSTGRES_READ_DB) -c \
		"SELECT pg_is_in_recovery() as is_replica, pg_last_wal_receive_lsn() as received, pg_last_wal_replay_lsn() as replayed;"

.PHONY: pg-replication-lag
pg-replication-lag: ## Check replication lag
	@echo "$(GREEN)=== Replication Lag ===$(NC)"
	@docker exec -it $(POSTGRES_REPLICA_CONTAINER_NAME) psql -U $(POSTGRES_SUPERUSER) -d $(POSTGRES_READ_DB) -c \
		"SELECT CASE WHEN pg_last_wal_receive_lsn() = pg_last_wal_replay_lsn() THEN 0 ELSE EXTRACT(EPOCH FROM now() - pg_last_xact_replay_timestamp()) END AS lag_seconds;"

.PHONY: pg-test-replication
pg-test-replication: ## Test replication by creating and checking a table
	@echo "$(GREEN)=== Testing Replication ===$(NC)"
	@echo "Creating test table on PRIMARY..."
	@docker exec -it $(POSTGRES_PRIMARY_CONTAINER_NAME) psql -U $(POSTGRES_SUPERUSER) -d $(POSTGRES_WRITE_DB) -c \
		"CREATE TABLE IF NOT EXISTS _replication_test (id serial, created_at timestamp default now()); INSERT INTO _replication_test DEFAULT VALUES;"
	@echo "Waiting for replication..."
	@sleep 2
	@echo "Checking table on REPLICA..."
	@docker exec -it $(POSTGRES_REPLICA_CONTAINER_NAME) psql -U $(POSTGRES_SUPERUSER) -d $(POSTGRES_READ_DB) -c \
		"SELECT * FROM _replication_test ORDER BY id DESC LIMIT 5;"
	@echo "$(GREEN)=== Replication test completed! ===$(NC)"

# =============================================================================
# DATABASE SHELLS
# =============================================================================

.PHONY: pg-primary
pg-primary: ## Connect to primary PostgreSQL as superuser
	docker exec -it $(POSTGRES_PRIMARY_CONTAINER_NAME) psql -U $(POSTGRES_SUPERUSER) -d $(POSTGRES_WRITE_DB)

.PHONY: pg-replica
pg-replica: ## Connect to replica PostgreSQL as superuser
	docker exec -it $(POSTGRES_REPLICA_CONTAINER_NAME) psql -U $(POSTGRES_SUPERUSER) -d $(POSTGRES_READ_DB)

.PHONY: pg-write
pg-write: ## Connect to primary as write user
	docker exec -it $(POSTGRES_PRIMARY_CONTAINER_NAME) psql -U $(POSTGRES_WRITE_USERNAME) -d $(POSTGRES_WRITE_DB)

.PHONY: pg-read
pg-read: ## Connect to replica as read user
	docker exec -it $(POSTGRES_REPLICA_CONTAINER_NAME) psql -U $(POSTGRES_READ_USERNAME) -d $(POSTGRES_READ_DB)

.PHONY: redis-cli
redis-cli: ## Connect to Redis CLI
	docker exec -it $(REDIS_CONTAINER_NAME) redis-cli -a $(REDIS_PASSWORD)

# =============================================================================
# MIGRATIONS
# =============================================================================

PG_MIGRATIONS_PATH := migrations/postgres
PG_NETWORK := hiring_backend
PG_DSN := postgres://$(POSTGRES_WRITE_USERNAME):$(POSTGRES_WRITE_PASSWORD)@$(POSTGRES_WRITE_HOST):5432/$(POSTGRES_WRITE_DB)?sslmode=disable

.PHONY: migrate-up
migrate-up: ## Run all pending migrations
	@echo "$(GREEN)Running migrations on PRIMARY...$(NC)"
	docker run --rm \
		-v $(PWD)/$(PG_MIGRATIONS_PATH):/migrations \
		--network $(PG_NETWORK) \
		migrate/migrate:latest \
		-path=/migrations \
		-database "$(PG_DSN)" \
		up
	@echo "$(GREEN)Migrations completed!$(NC)"
	@echo "$(YELLOW)Note: Changes will automatically replicate to REPLICA$(NC)"

.PHONY: migrate-down
migrate-down: ## Rollback last migration
	@echo "$(YELLOW)Rolling back last migration on PRIMARY...$(NC)"
	docker run --rm \
		-v $(PWD)/$(PG_MIGRATIONS_PATH):/migrations \
		--network $(PG_NETWORK) \
		migrate/migrate:latest \
		-path=/migrations \
		-database "$(PG_DSN)" \
		down 1

.PHONY: migrate-status
migrate-status: ## Show migration status
	docker run --rm \
		-v $(PWD)/$(PG_MIGRATIONS_PATH):/migrations \
		--network $(PG_NETWORK) \
		migrate/migrate:latest \
		-path=/migrations \
		-database "$(PG_DSN)" \
		version

.PHONY: migrate-create
migrate-create: ## Create new migration (usage: make migrate-create name=create_users)
	@if [ -z "$(name)" ]; then \
		echo "$(RED)Error: Please provide migration name$(NC)"; \
		echo "Usage: make migrate-create name=create_users"; \
		exit 1; \
	fi
	docker run --rm \
		-v $(PWD)/$(PG_MIGRATIONS_PATH):/migrations \
		migrate/migrate:latest \
		create -ext sql -dir /migrations -seq $(name)
	@echo "$(GREEN)Migration created!$(NC)"

# =============================================================================
# BUILD
# =============================================================================

.PHONY: build
build: ## Build Go applications
	@echo "$(GREEN)Building applications...$(NC)"
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/api ./cmd/api/main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/worker ./cmd/worker/main.go
	@echo "$(GREEN)Build completed!$(NC)"

# =============================================================================
# UTILITIES
# =============================================================================

.PHONY: check-env
check-env: ## Check if .env.app and .env.infra files exist
	@if [ ! -f .env.app ]; then \
		echo "$(RED)Error: .env.app not found!$(NC)"; \
		echo "Run: cp .env.app.example .env.app"; \
		exit 1; \
	fi
	@if [ ! -f .env.infra ]; then \
		echo "$(RED)Error: .env.infra not found!$(NC)"; \
		echo "Run: cp .env.infra.example .env.infra"; \
		exit 1; \
	fi

.PHONY: env-init
env-init: ## Create .env.app and .env.infra from examples
	@if [ -f .env.app ]; then \
		echo "$(YELLOW).env.app already exists. Overwrite? [y/N]$(NC)"; \
		read confirm && [ "$$confirm" = "y" ] && cp .env.app.example .env.app && echo "$(GREEN).env.app created$(NC)" || echo "Skipped"; \
	else \
		cp .env.app.example .env.app; \
		echo "$(GREEN).env.app created$(NC)"; \
	fi
	@if [ -f .env.infra ]; then \
		echo "$(YELLOW).env.infra already exists. Overwrite? [y/N]$(NC)"; \
		read confirm && [ "$$confirm" = "y" ] && cp .env.infra.example .env.infra && echo "$(GREEN).env.infra created$(NC)" || echo "Skipped"; \
	else \
		cp .env.infra.example .env.infra; \
		echo "$(GREEN).env.infra created$(NC)"; \
	fi

.PHONY: clean
clean: ## Clean build artifacts and docker resources
	@echo "$(YELLOW)Cleaning...$(NC)"
	rm -rf bin/
	rm -f .env
	docker compose down --rmi local --remove-orphans
	@echo "$(GREEN)Cleaned!$(NC)"

.PHONY: help
help: ## Show this help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help