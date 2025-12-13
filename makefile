# =============================================================================
# HIRING SERVICE - MAKEFILE
# =============================================================================

# Environment directories
ENV_DIR := .envs
ENV_APP_DIR := $(ENV_DIR)/app
ENV_INFRA_DIR := $(ENV_DIR)/infra
ENV_APP_SERVICES_DIR := $(ENV_APP_DIR)/services
ENV_INFRA_SERVICES_DIR := $(ENV_INFRA_DIR)/services
ENV_APP_TEMPLATES_DIR := $(ENV_APP_DIR)/templates
ENV_INFRA_TEMPLATES_DIR := $(ENV_INFRA_DIR)/templates

# Load environment file
ifneq (,$(wildcard ./$(ENV_INFRA_DIR)/.env))
    include $(ENV_INFRA_DIR)/.env
    export
endif

# Colors for output
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RED    := \033[0;31m
NC     := \033[0m # No Color

# Docker compose with env file
DC := docker compose --env-file $(ENV_DIR)/.env

# =============================================================================
# ENVIRONMENT MANAGEMENT
# =============================================================================

.PHONY: env
env: ## Generate root .env from infra/.env
	@echo "$(GREEN)Generating .env...$(NC)"
	@cat /dev/null > $(ENV_DIR)/.env
	@echo "# ==============================================================================" >> $(ENV_DIR)/.env
	@echo "# AUTO-GENERATED - DO NOT EDIT" >> $(ENV_DIR)/.env
	@echo "# Generated from infra/.env" >> $(ENV_DIR)/.env
	@echo "# Run: make env" >> $(ENV_DIR)/.env
	@echo "# ==============================================================================" >> $(ENV_DIR)/.env
	@echo "" >> $(ENV_DIR)/.env
	@grep -v '^\s*#' $(ENV_INFRA_DIR)/.env | grep -v '^\s*$$' >> $(ENV_DIR)/.env || true
	@echo "$(GREEN).env generated!$(NC)"

.PHONY: env-infra
env-infra: ## Generate service env files from templates
	@echo "$(GREEN)Generating service env files from templates...$(NC)"
	@# Ensure directories exist
	@mkdir -p $(ENV_APP_SERVICES_DIR)/api
	@mkdir -p $(ENV_APP_SERVICES_DIR)/worker
	@mkdir -p $(ENV_INFRA_SERVICES_DIR)/postgres-primary
	@mkdir -p $(ENV_INFRA_SERVICES_DIR)/postgres-replica
	@mkdir -p $(ENV_INFRA_SERVICES_DIR)/redis
	@mkdir -p $(ENV_INFRA_SERVICES_DIR)/rabbitmq
	@# Load .env and generate service env files
	@set -a && . $(ENV_INFRA_DIR)/.env && set +a && \
		if [ -f $(ENV_APP_TEMPLATES_DIR)/.env.api.template ]; then \
			envsubst < $(ENV_APP_TEMPLATES_DIR)/.env.api.template > $(ENV_APP_SERVICES_DIR)/api/.env; \
			echo "  $(GREEN)✓$(NC) app/services/api/.env"; \
		fi && \
		if [ -f $(ENV_APP_TEMPLATES_DIR)/.env.worker.template ]; then \
			envsubst < $(ENV_APP_TEMPLATES_DIR)/.env.worker.template > $(ENV_APP_SERVICES_DIR)/worker/.env; \
			echo "  $(GREEN)✓$(NC) app/services/worker/.env"; \
		fi && \
		envsubst < $(ENV_INFRA_TEMPLATES_DIR)/.env.postgres-primary.template > $(ENV_INFRA_SERVICES_DIR)/postgres-primary/.env && \
		echo "  $(GREEN)✓$(NC) infra/services/postgres-primary/.env" && \
		envsubst < $(ENV_INFRA_TEMPLATES_DIR)/.env.postgres-replica.template > $(ENV_INFRA_SERVICES_DIR)/postgres-replica/.env && \
		echo "  $(GREEN)✓$(NC) infra/services/postgres-replica/.env" && \
		envsubst < $(ENV_INFRA_TEMPLATES_DIR)/.env.redis.template > $(ENV_INFRA_SERVICES_DIR)/redis/.env && \
		echo "  $(GREEN)✓$(NC) infra/services/redis/.env" && \
		envsubst < $(ENV_INFRA_TEMPLATES_DIR)/.env.rabbitmq.template > $(ENV_INFRA_SERVICES_DIR)/rabbitmq/.env && \
		echo "  $(GREEN)✓$(NC) infra/services/rabbitmq/.env"
	@echo "$(GREEN)All service env files generated!$(NC)"

.PHONY: check-env
check-env: ## Check if infra/.env exists
	@if [ ! -f $(ENV_INFRA_DIR)/.env ]; then \
		echo "$(RED)Error: $(ENV_INFRA_DIR)/.env not found!$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Infra env exists!$(NC)"

# =============================================================================
# DOCKER BUILD
# =============================================================================

# Image names
API_IMAGE := ${COMPANY_DOCKER_REGISTRY_HOST}/${COMPANY_HIRING_PROJECT_NAME}/${COMPANY_HIRING_API_SERVICE_NAME}:${COMPANY_HIRING_VERSION}
WORKER_IMAGE := ${COMPANY_DOCKER_REGISTRY_HOST}/${COMPANY_HIRING_PROJECT_NAME}/${COMPANY_HIRING_WORKER_SERVICE_NAME}:${COMPANY_HIRING_VERSION}

.PHONY: docker-build-api
docker-build-api: ## Build API Docker image
	@echo "$(GREEN)Building API image...$(NC)"
	docker build -f builds/Dockerfile -t $(API_IMAGE) --target api .
	@echo "$(GREEN)API image built: $(API_IMAGE)$(NC)"

.PHONY: docker-build-worker
docker-build-worker: ## Build Worker Docker image
	@echo "$(GREEN)Building Worker image...$(NC)"
	docker build -f builds/Dockerfile -t $(WORKER_IMAGE) --target worker .
	@echo "$(GREEN)Worker image built: $(WORKER_IMAGE)$(NC)"

.PHONY: docker-build
docker-build: docker-build-api docker-build-worker ## Build all Docker images
	@echo "$(GREEN)All images built!$(NC)"

# =============================================================================
# RUN SERVICES
# =============================================================================

.PHONY: check-api-image
check-api-image:
	@if ! docker image inspect $(API_IMAGE) >/dev/null 2>&1; then \
		echo "$(YELLOW)API image not found, building...$(NC)"; \
		$(MAKE) docker-build-api; \
	fi

.PHONY: check-worker-image
check-worker-image:
	@if ! docker image inspect $(WORKER_IMAGE) >/dev/null 2>&1; then \
		echo "$(YELLOW)Worker image not found, building...$(NC)"; \
		$(MAKE) docker-build-worker; \
	fi

.PHONY: run-api
run-api: check-env env env-infra docker-build-api ## Run API service
	@echo "$(GREEN)Starting API service...$(NC)"
	$(DC) up -d api --force-recreate
	@echo "$(GREEN)API started!$(NC)"
	@$(MAKE) status

.PHONY: run-worker
run-worker: check-env env env-infra docker-build-worker ## Run Worker service
	@echo "$(GREEN)Starting Worker service...$(NC)"
	$(DC) up -d worker --force-recreate
	@echo "$(GREEN)Worker started!$(NC)"
	@$(MAKE) status

.PHONY: run
run: check-env env env-infra check-api-image check-worker-image ## Run API and Worker services
	@echo "$(GREEN)Starting API and Worker services...$(NC)"
	$(DC) up -d api worker
	@echo "$(GREEN)Services started!$(NC)"
	@$(MAKE) status

# =============================================================================
# DOCKER COMPOSE COMMANDS
# =============================================================================

.PHONY: config
config: ## Show resolved docker compose configuration
	$(DC) config

.PHONY: up
up: check-env env env-infra ## Start all services
	@echo "$(GREEN)Starting all services...$(NC)"
	$(DC) up -d
	@echo "$(GREEN)Services started!$(NC)"
	@$(MAKE) status

.PHONY: up-build
up-build: check-env env env-infra ## Start all services with rebuild
	@echo "$(GREEN)Building and starting all services...$(NC)"
	$(DC) up -d --build
	@echo "$(GREEN)Services started!$(NC)"
	@$(MAKE) status

.PHONY: down
down: ## Stop all services
	@echo "$(YELLOW)Stopping all services...$(NC)"
	$(DC) down
	@echo "$(GREEN)Services stopped!$(NC)"

.PHONY: down-v
down-v: ## Stop all services and remove volumes (DESTRUCTIVE!)
	@echo "$(RED)Stopping all services and removing volumes...$(NC)"
	@read -p "Are you sure? This will delete all data! [y/N] " confirm && [ "$$confirm" = "y" ]
	$(DC) down -v
	@echo "$(GREEN)Services stopped and volumes removed!$(NC)"

.PHONY: restart
restart: down up ## Restart all services

.PHONY: logs
logs: ## Show logs for all services
	$(DC) logs -f

.PHONY: logs-db
logs-db: ## Show PostgreSQL logs (primary and replica)
	$(DC) logs -f postgres_primary postgres_replica

.PHONY: logs-app
logs-app: ## Show api and worker logs
	$(DC) logs -f api worker

.PHONY: logs-api
logs-api: ## Show api logs
	$(DC) logs -f api

.PHONY: logs-worker
logs-worker: ## Show worker logs
	$(DC) logs -f worker

.PHONY: status
status: ## Show status of all services
	@echo ""
	@echo "$(GREEN)=== Services Status ===$(NC)"
	@$(DC) ps
	@echo ""

# =============================================================================
# INFRASTRUCTURE ONLY (without app services)
# =============================================================================

.PHONY: infra-up
infra-up: check-env env env-infra ## Start only infrastructure (postgres, redis, rabbitmq)
	@echo "$(GREEN)Starting infrastructure services...$(NC)"
	$(DC) up -d postgres_primary postgres_replica redis rabbitmq
	@echo "$(GREEN)Infrastructure started!$(NC)"
	@$(MAKE) status

.PHONY: infra-down
infra-down: ## Stop infrastructure services
	@echo "$(YELLOW)Stopping infrastructure services...$(NC)"
	$(DC) stop postgres_primary postgres_replica redis rabbitmq
	@echo "$(GREEN)Infrastructure stopped!$(NC)"

# =============================================================================
# POSTGRESQL REPLICATION
# =============================================================================

.PHONY: pg-replication-status
pg-replication-status: ## Check PostgreSQL replication status
	@echo "$(GREEN)=== Primary: Replication Status ===$(NC)"
	@docker exec -i -e PGPASSWORD=$(COMPANY_HIRING_POSTGRES_SUPERUSER_PASSWORD) \
	$(COMPANY_HIRING_POSTGRES_PRIMARY_HOST) \
	psql -U $(COMPANY_HIRING_POSTGRES_SUPERUSER_USERNAME) -d $(COMPANY_HIRING_POSTGRES_DATABASE) -c \
		"SELECT client_addr, state, sent_lsn, write_lsn, flush_lsn, replay_lsn, sync_state FROM pg_stat_replication;"
	@echo ""
	@echo "$(GREEN)=== Primary: Replication Slots ===$(NC)"
	@docker exec -i -e PGPASSWORD=$(COMPANY_HIRING_POSTGRES_SUPERUSER_PASSWORD) \
	$(COMPANY_HIRING_POSTGRES_PRIMARY_HOST) \
	psql -U $(COMPANY_HIRING_POSTGRES_SUPERUSER_USERNAME) -d $(COMPANY_HIRING_POSTGRES_DATABASE) -c \
		"SELECT slot_name, slot_type, active, restart_lsn FROM pg_replication_slots;"
	@echo ""
	@echo "$(GREEN)=== Replica: Recovery Status ===$(NC)"
	@docker exec -i -e PGPASSWORD=$(COMPANY_HIRING_POSTGRES_SUPERUSER_PASSWORD) \
	$(COMPANY_HIRING_POSTGRES_REPLICA_HOST) \
	psql -U $(COMPANY_HIRING_POSTGRES_SUPERUSER_USERNAME) -d $(COMPANY_HIRING_POSTGRES_DATABASE) -c \
		"SELECT pg_is_in_recovery() as is_replica, pg_last_wal_receive_lsn() as received, pg_last_wal_replay_lsn() as replayed;"

.PHONY: pg-replication-lag
pg-replication-lag: ## Check replication lag
	@echo "$(GREEN)=== Replication Lag ===$(NC)"
	@docker exec -it $(COMPANY_HIRING_POSTGRES_REPLICA_HOST) psql -U $(COMPANY_HIRING_POSTGRES_SUPERUSER_USERNAME) -d $(COMPANY_HIRING_POSTGRES_DATABASE) -c \
		"SELECT CASE WHEN pg_last_wal_receive_lsn() = pg_last_wal_replay_lsn() THEN 0 ELSE EXTRACT(EPOCH FROM now() - pg_last_xact_replay_timestamp()) END AS lag_seconds;"

.PHONY: pg-test-replication
pg-test-replication: ## Test replication by creating and checking a table
	@echo "$(GREEN)=== Testing Replication ===$(NC)"
	@echo "Creating test table on PRIMARY..."
	@docker exec -i $(COMPANY_HIRING_POSTGRES_PRIMARY_HOST) psql -U $(COMPANY_HIRING_POSTGRES_SUPERUSER_USERNAME) -d $(COMPANY_HIRING_POSTGRES_DATABASE) -c \
		"CREATE TABLE IF NOT EXISTS _replication_test (id serial, created_at timestamp default now()); INSERT INTO _replication_test DEFAULT VALUES;"
	@echo "Waiting for replication..."
	@sleep 2
	@echo "Checking table on REPLICA..."
	@docker exec -i $(COMPANY_HIRING_POSTGRES_REPLICA_HOST) psql -U $(COMPANY_HIRING_POSTGRES_SUPERUSER_USERNAME) -d $(COMPANY_HIRING_POSTGRES_DATABASE) -c \
		"SELECT * FROM _replication_test ORDER BY id DESC LIMIT 5;"
	@echo "$(GREEN)=== Replication test completed! ===$(NC)"

# =============================================================================
# DATABASE SHELLS
# =============================================================================

.PHONY: pg-primary
pg-primary: ## Connect to primary PostgreSQL as superuser
	docker exec -it -e PGPASSWORD=$(COMPANY_HIRING_POSTGRES_SUPERUSER_PASSWORD) \
		$(COMPANY_HIRING_POSTGRES_PRIMARY_HOST) \
		psql -U $(COMPANY_HIRING_POSTGRES_SUPERUSER_USERNAME) -d $(COMPANY_HIRING_POSTGRES_DATABASE)

.PHONY: pg-replica
pg-replica: ## Connect to replica PostgreSQL as superuser
	docker exec -it -e PGPASSWORD=$(COMPANY_HIRING_POSTGRES_SUPERUSER_PASSWORD) \
	$(COMPANY_HIRING_POSTGRES_REPLICA_HOST) \
	psql -U $(COMPANY_HIRING_POSTGRES_SUPERUSER_USERNAME) -d $(COMPANY_HIRING_POSTGRES_DATABASE)

.PHONY: pg-write
pg-write: ## Connect to primary as write user
	docker exec -it -e PGPASSWORD=$(COMPANY_HIRING_POSTGRES_PRIMARY_PASSWORD) \
	$(COMPANY_HIRING_POSTGRES_PRIMARY_HOST) \
	psql -U $(COMPANY_HIRING_POSTGRES_PRIMARY_USERNAME) -d $(COMPANY_HIRING_POSTGRES_DATABASE)

.PHONY: pg-read
pg-read: ## Connect to replica as read user
	docker exec -it -e PGPASSWORD=$(COMPANY_HIRING_POSTGRES_REPLICA_PASSWORD) \
	$(COMPANY_HIRING_POSTGRES_REPLICA_HOST) \
	psql -U $(COMPANY_HIRING_POSTGRES_REPLICA_USERNAME) -d $(COMPANY_HIRING_POSTGRES_DATABASE)

.PHONY: redis-cli
redis-cli: ## Connect to Redis CLI
	docker exec -it $(COMPANY_HIRING_REDIS_HOST) redis-cli -a $(COMPANY_HIRING_REDIS_PASSWORD)

# =============================================================================
# DATABASE MONITORING
# =============================================================================

.PHONY: pg-write-monitor
pg-write-monitor: ## Monitor active queries on primary (write) - live refresh
	@echo "$(GREEN)Monitoring queries on PRIMARY (Ctrl+C to exit)...$(NC)"
	watch -n 2 'docker exec -e PGPASSWORD=$(COMPANY_HIRING_POSTGRES_SUPERUSER_PASSWORD) \
		$(COMPANY_HIRING_POSTGRES_PRIMARY_HOST) \
		psql -U $(COMPANY_HIRING_POSTGRES_SUPERUSER_USERNAME) -d $(COMPANY_HIRING_POSTGRES_DATABASE) -c \
		"SELECT pid, usename, client_addr, state, wait_event_type, left(query, 80) as query, now() - query_start as duration FROM pg_stat_activity WHERE datname = '\''$(COMPANY_HIRING_POSTGRES_DATABASE)'\'' AND pid <> pg_backend_pid() ORDER BY query_start DESC NULLS LAST LIMIT 20;"'

.PHONY: pg-read-monitor
pg-read-monitor: ## Monitor active queries on replica (read) - live refresh
	@echo "$(GREEN)Monitoring queries on REPLICA (Ctrl+C to exit)...$(NC)"
	watch -n 2 'docker exec -e PGPASSWORD=$(COMPANY_HIRING_POSTGRES_SUPERUSER_PASSWORD) \
		$(COMPANY_HIRING_POSTGRES_REPLICA_HOST) \
		psql -U $(COMPANY_HIRING_POSTGRES_SUPERUSER_USERNAME) -d $(COMPANY_HIRING_POSTGRES_DATABASE) -c \
		"SELECT pid, usename, client_addr, state, wait_event_type, left(query, 80) as query, now() - query_start as duration FROM pg_stat_activity WHERE datname = '\''$(COMPANY_HIRING_POSTGRES_DATABASE)'\'' AND pid <> pg_backend_pid() ORDER BY query_start DESC NULLS LAST LIMIT 20;"'

# =============================================================================
# DATABASE QUERY LOGS
# =============================================================================

.PHONY: pg-write-log
pg-write-log: ## Tail live query log on primary (write)
	@echo "$(GREEN)Tailing query log on PRIMARY (Ctrl+C to exit)...$(NC)"
	@echo "$(YELLOW)Note: Requires log_statement='all' in postgresql.conf$(NC)"
	@echo ""
	docker logs -f $(COMPANY_HIRING_POSTGRES_PRIMARY_HOST) 2>&1 | grep -E --line-buffered "(statement:|duration:|ERROR|FATAL)"

.PHONY: pg-read-log
pg-read-log: ## Tail live query log on replica (read)
	@echo "$(GREEN)Tailing query log on REPLICA (Ctrl+C to exit)...$(NC)"
	@echo "$(YELLOW)Note: Requires log_statement='all' in postgresql.conf$(NC)"
	@echo ""
	docker logs -f $(COMPANY_HIRING_POSTGRES_REPLICA_HOST) 2>&1 | grep -E --line-buffered "(statement:|duration:|ERROR|FATAL)"

.PHONY: pg-write-log-full
pg-write-log-full: ## Tail full log on primary (unfiltered)
	docker logs -f $(COMPANY_HIRING_POSTGRES_PRIMARY_HOST)

.PHONY: pg-read-log-full
pg-read-log-full: ## Tail full log on replica (unfiltered)
	docker logs -f $(COMPANY_HIRING_POSTGRES_REPLICA_HOST)

# =============================================================================
# MIGRATIONS
# =============================================================================

PG_MIGRATIONS_PATH := migrations/postgres
PG_NETWORK := hiring_backend
PG_DSN := postgres://$(COMPANY_HIRING_POSTGRES_PRIMARY_USERNAME):$(COMPANY_HIRING_POSTGRES_PRIMARY_PASSWORD)@$(COMPANY_HIRING_POSTGRES_PRIMARY_HOST):5432/$(COMPANY_HIRING_POSTGRES_DATABASE)?sslmode=disable

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

.PHONY: clean
clean: ## Clean build artifacts and docker resources
	@echo "$(YELLOW)Cleaning...$(NC)"
	rm -rf bin/
	rm -f $(ENV_DIR)/.env
	rm -f $(ENV_APP_SERVICES_DIR)/api/.env
	rm -f $(ENV_APP_SERVICES_DIR)/worker/.env
	rm -f $(ENV_INFRA_SERVICES_DIR)/postgres-primary/.env
	rm -f $(ENV_INFRA_SERVICES_DIR)/postgres-replica/.env
	rm -f $(ENV_INFRA_SERVICES_DIR)/redis/.env
	rm -f $(ENV_INFRA_SERVICES_DIR)/rabbitmq/.env
	$(DC) down --rmi local --remove-orphans 2>/dev/null || true
	@echo "$(GREEN)Cleaned!$(NC)"

# =============================================================================
# AUDIT
# =============================================================================

.PHONY: audit
audit: ## Run full system audit (env, images, status, replication)
	@echo "$(GREEN)========================================$(NC)"
	@echo "$(GREEN)  SYSTEM AUDIT$(NC)"
	@echo "$(GREEN)========================================$(NC)"
	@echo ""
	@echo "$(GREEN)[1/5] Checking environment...$(NC)"
	@$(MAKE) check-env
	@sleep 1
	@echo ""
	@echo "$(GREEN)[2/5] Checking Docker images...$(NC)"
	@if docker image inspect $(API_IMAGE) >/dev/null 2>&1; then \
		echo "  $(GREEN)✓$(NC) API image: $(API_IMAGE)"; \
	else \
		echo "  $(RED)✗$(NC) API image: $(API_IMAGE) - NOT FOUND"; \
	fi
	@if docker image inspect $(WORKER_IMAGE) >/dev/null 2>&1; then \
		echo "  $(GREEN)✓$(NC) Worker image: $(WORKER_IMAGE)"; \
	else \
		echo "  $(RED)✗$(NC) Worker image: $(WORKER_IMAGE) - NOT FOUND"; \
	fi
	@sleep 1
	@echo ""
	@echo "$(GREEN)[3/5] Services status...$(NC)"
	@$(DC) ps
	@sleep 1
	@echo ""
	@echo "$(GREEN)[4/5] PostgreSQL replication status...$(NC)"
	@if docker ps --format '{{.Names}}' | grep -q $(COMPANY_HIRING_POSTGRES_PRIMARY_HOST); then \
		$(MAKE) pg-replication-status; \
	else \
		echo "  $(YELLOW)⚠$(NC) PostgreSQL Primary not running, skipping replication check"; \
	fi
	@sleep 1
	@echo ""
	@echo "$(GREEN)[5/5] Replication lag...$(NC)"
	@if docker ps --format '{{.Names}}' | grep -q $(COMPANY_HIRING_POSTGRES_REPLICA_HOST); then \
		$(MAKE) pg-replication-lag; \
	else \
		echo "  $(YELLOW)⚠$(NC) PostgreSQL Replica not running, skipping lag check"; \
	fi
	@echo ""
	@echo "$(GREEN)========================================$(NC)"
	@echo "$(GREEN)  AUDIT COMPLETE$(NC)"
	@echo "$(GREEN)========================================$(NC)"

# =============================================================================
# HELP
# =============================================================================

.PHONY: help
help: ## Show this help
	@echo "Usage: make [target]"
	@echo ""
	@echo "$(GREEN)Environment:$(NC)"
	@awk -F ':.*?## ' '/^env:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "env", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^env-infra:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "env-infra", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^check-env:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "check-env", $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Docker Compose:$(NC)"
	@awk -F ':.*?## ' '/^up:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "up", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^up-build:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "up-build", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^down:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "down", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^down-v:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "down-v", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^restart:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "restart", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^logs:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "logs", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^logs-db:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "logs-db", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^status:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "status", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^config:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "config", $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Infrastructure:$(NC)"
	@awk -F ':.*?## ' '/^infra-up:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "infra-up", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^infra-down:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "infra-down", $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Docker Build:$(NC)"
	@awk -F ':.*?## ' '/^docker-build:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "docker-build", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^docker-build-api:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "docker-build-api", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^docker-build-worker:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "docker-build-worker", $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Run Services:$(NC)"
	@awk -F ':.*?## ' '/^run:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "run", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^run-api:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "run-api", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^run-worker:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "run-worker", $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)PostgreSQL Replication:$(NC)"
	@awk -F ':.*?## ' '/^pg-replication-status:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "pg-replication-status", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^pg-replication-lag:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "pg-replication-lag", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^pg-test-replication:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "pg-test-replication", $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Database Shells:$(NC)"
	@awk -F ':.*?## ' '/^pg-primary:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "pg-primary", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^pg-replica:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "pg-replica", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^pg-write:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "pg-write", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^pg-read:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "pg-read", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^redis-cli:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "redis-cli", $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Migrations:$(NC)"
	@awk -F ':.*?## ' '/^migrate-up:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "migrate-up", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^migrate-down:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "migrate-down", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^migrate-status:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "migrate-status", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^migrate-create:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "migrate-create", $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Build & Utilities:$(NC)"
	@awk -F ':.*?## ' '/^build:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "build", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^clean:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "clean", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^audit:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "audit", $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Database Monitoring:$(NC)"
	@awk -F ':.*?## ' '/^pg-write-monitor:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "pg-write-monitor", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^pg-read-monitor:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "pg-read-monitor", $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Database Query Logs:$(NC)"
	@awk -F ':.*?## ' '/^pg-write-log:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "pg-write-log", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^pg-read-log:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "pg-read-log", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^pg-write-log-full:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "pg-write-log-full", $$2}' $(MAKEFILE_LIST)
	@awk -F ':.*?## ' '/^pg-read-log-full:.*## /{printf "  \033[36m%-24s\033[0m %s\n", "pg-read-log-full", $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help