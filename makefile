# =============================================================================
# HIRING SERVICE - MAKEFILE
# =============================================================================

# =============================================================================
# SERVICE DEFINITIONS
# =============================================================================

# Application services (used for build, run, env generation)
APP_SERVICES := auth blog notifications

# Infrastructure services (used for env generation, infra-up/down)
INFRA_SERVICES := postgres-primary postgres-replica redis kafka

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

# Kafka CLI base path (apache/kafka image)
KAFKA_BIN := /opt/kafka/bin

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
	@set -a && . $(ENV_INFRA_DIR)/.env && set +a && \
		for svc in $(APP_SERVICES); do \
			mkdir -p $(ENV_APP_SERVICES_DIR)/$$svc; \
			if [ -f $(ENV_APP_TEMPLATES_DIR)/.env.$$svc.template ]; then \
				envsubst < $(ENV_APP_TEMPLATES_DIR)/.env.$$svc.template > $(ENV_APP_SERVICES_DIR)/$$svc/.env; \
				echo "  $(GREEN)✓$(NC) app/services/$$svc/.env"; \
			else \
				echo "  $(YELLOW)⚠$(NC) app/services/$$svc/.env — template not found, skipped"; \
			fi; \
		done && \
		for svc in $(INFRA_SERVICES); do \
			mkdir -p $(ENV_INFRA_SERVICES_DIR)/$$svc; \
			if [ -f $(ENV_INFRA_TEMPLATES_DIR)/.env.$$svc.template ]; then \
				envsubst < $(ENV_INFRA_TEMPLATES_DIR)/.env.$$svc.template > $(ENV_INFRA_SERVICES_DIR)/$$svc/.env; \
				echo "  $(GREEN)✓$(NC) infra/services/$$svc/.env"; \
			else \
				echo "  $(YELLOW)⚠$(NC) infra/services/$$svc/.env — template not found, skipped"; \
			fi; \
		done
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

# Image name helper: $(call image_name,service)
define image_name
${COMPANY_DOCKER_REGISTRY_HOST}/${COMPANY_HIRING_PROJECT_NAME}/$(1):${COMPANY_HIRING_VERSION}
endef

AUTH_IMAGE := $(call image_name,auth)
BLOG_IMAGE := $(call image_name,blog)
NOTIFICATIONS_IMAGE := $(call image_name,notifications)

.PHONY: docker-build-auth
docker-build-auth: ## Build Auth Docker image
	@echo "$(GREEN)Building Auth image...$(NC)"
	docker build -f builds/Dockerfile -t $(AUTH_IMAGE) --target auth .
	@echo "$(GREEN)Auth image built: $(AUTH_IMAGE)$(NC)"

.PHONY: docker-build-blog
docker-build-blog: ## Build Blog Docker image
	@echo "$(GREEN)Building Blog image...$(NC)"
	docker build -f builds/Dockerfile -t $(BLOG_IMAGE) --target blog .
	@echo "$(GREEN)Blog image built: $(BLOG_IMAGE)$(NC)"

.PHONY: docker-build-notifications
docker-build-notifications: ## Build Notifications Docker image
	@echo "$(GREEN)Building Notifications image...$(NC)"
	docker build -f builds/Dockerfile -t $(NOTIFICATIONS_IMAGE) --target notifications .
	@echo "$(GREEN)Notifications image built: $(NOTIFICATIONS_IMAGE)$(NC)"

.PHONY: docker-build
docker-build: docker-build-auth docker-build-blog docker-build-notifications ## Build all Docker images
	@echo "$(GREEN)All images built!$(NC)"

# =============================================================================
# RUN SERVICES
# =============================================================================

.PHONY: check-images
check-images:
	@for svc in $(APP_SERVICES); do \
		img=$(call image_name,$$svc); \
		if ! docker image inspect $$img >/dev/null 2>&1; then \
			echo "$(YELLOW)$$svc image not found, building...$(NC)"; \
			$(MAKE) docker-build-$$svc; \
		fi; \
	done

.PHONY: run-auth
run-auth: check-env env env-infra docker-build-auth ## Run Auth service
	@echo "$(GREEN)Starting Auth service...$(NC)"
	$(DC) up -d auth --force-recreate
	@echo "$(GREEN)Auth started!$(NC)"
	@$(MAKE) status

.PHONY: run-blog
run-blog: check-env env env-infra docker-build-blog ## Run Blog service
	@echo "$(GREEN)Starting Blog service...$(NC)"
	$(DC) up -d blog --force-recreate
	@echo "$(GREEN)Blog started!$(NC)"
	@$(MAKE) status

.PHONY: run-notifications
run-notifications: check-env env env-infra docker-build-notifications ## Run Notifications service
	@echo "$(GREEN)Starting Notifications service...$(NC)"
	$(DC) up -d notifications --force-recreate
	@echo "$(GREEN)Notifications started!$(NC)"
	@$(MAKE) status

.PHONY: run
run: check-env env env-infra check-images ## Run all application services
	@echo "$(GREEN)Starting all application services...$(NC)"
	$(DC) up -d auth blog notifications
	@echo "$(GREEN)Services started!$(NC)"
	@$(MAKE) status

# =============================================================================
# DOCKER COMPOSE COMMANDS
# =============================================================================

.PHONY: config
config: ## Show resolved docker compose configuration
	$(DC) config

.PHONY: up
up: check-env env env-infra docker-build ## Start all services
	@echo "$(GREEN)Starting all services...$(NC)"
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
logs-app: ## Show all application service logs
	$(DC) logs -f auth blog notifications

.PHONY: logs-auth
logs-auth: ## Show auth logs
	$(DC) logs -f auth

.PHONY: logs-blog
logs-blog: ## Show blog logs
	$(DC) logs -f blog

.PHONY: logs-notifications
logs-notifications: ## Show notifications logs
	$(DC) logs -f notifications

.PHONY: logs-kafka
logs-kafka: ## Show Kafka logs
	$(DC) logs -f kafka

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
infra-up: check-env env env-infra ## Start only infrastructure (postgres, redis, kafka)
	@echo "$(GREEN)Starting infrastructure services...$(NC)"
	$(DC) up -d postgres_primary postgres_replica redis kafka
	@echo "$(GREEN)Infrastructure started!$(NC)"
	@$(MAKE) status

.PHONY: infra-down
infra-down: ## Stop infrastructure services
	@echo "$(YELLOW)Stopping infrastructure services...$(NC)"
	$(DC) stop postgres_primary postgres_replica redis kafka
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
# KAFKA
# =============================================================================

.PHONY: kafka-topics
kafka-topics: ## List all Kafka topics
	@echo "$(GREEN)=== Kafka Topics ===$(NC)"
	docker exec -it $(COMPANY_HIRING_KAFKA_HOST) $(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server localhost:9092 --list

.PHONY: kafka-create-topic
kafka-create-topic: ## Create Kafka topic (usage: make kafka-create-topic name=post.created partitions=3)
	@if [ -z "$(name)" ]; then \
		echo "$(RED)Error: Please provide topic name$(NC)"; \
		echo "Usage: make kafka-create-topic name=post.created partitions=3"; \
		exit 1; \
	fi
	docker exec -it $(COMPANY_HIRING_KAFKA_HOST) $(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server localhost:9092 \
		--create --topic $(name) \
		--partitions $(or $(partitions),3) \
		--replication-factor 1
	@echo "$(GREEN)Topic '$(name)' created!$(NC)"

.PHONY: kafka-describe-topic
kafka-describe-topic: ## Describe Kafka topic (usage: make kafka-describe-topic name=post.created)
	@if [ -z "$(name)" ]; then \
		echo "$(RED)Error: Please provide topic name$(NC)"; \
		echo "Usage: make kafka-describe-topic name=post.created"; \
		exit 1; \
	fi
	docker exec -it $(COMPANY_HIRING_KAFKA_HOST) $(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server localhost:9092 \
		--describe --topic $(name)

.PHONY: kafka-consumer-groups
kafka-consumer-groups: ## List Kafka consumer groups
	@echo "$(GREEN)=== Kafka Consumer Groups ===$(NC)"
	docker exec -it $(COMPANY_HIRING_KAFKA_HOST) $(KAFKA_BIN)/kafka-consumer-groups.sh \
		--bootstrap-server localhost:9092 --list

# =============================================================================
# MIGRATIONS
# =============================================================================

PG_MIGRATIONS_PATH := migrations/postgres
PG_NETWORK := hiring_backend
PG_DSN := postgres://$(COMPANY_HIRING_POSTGRES_PRIMARY_USERNAME):$(COMPANY_HIRING_POSTGRES_PRIMARY_PASSWORD)@$(COMPANY_HIRING_POSTGRES_PRIMARY_HOST):5432/$(COMPANY_HIRING_POSTGRES_DATABASE)?sslmode=disable

path ?= .
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
	@for svc in $(APP_SERVICES); do \
		echo "  Building $$svc..."; \
		CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/$$svc ./cmd/$$svc/main.go; \
		echo "  $(GREEN)✓$(NC) bin/$$svc"; \
	done
	@echo "$(GREEN)Build completed!$(NC)"

# =============================================================================
# UTILITIES
# =============================================================================

.PHONY: clean
clean: ## Clean build artifacts and docker resources
	@echo "$(YELLOW)Cleaning...$(NC)"
	rm -rf bin/
	rm -f $(ENV_DIR)/.env
	@for svc in $(APP_SERVICES); do \
		rm -f $(ENV_APP_SERVICES_DIR)/$$svc/.env; \
	done
	@for svc in $(INFRA_SERVICES); do \
		rm -f $(ENV_INFRA_SERVICES_DIR)/$$svc/.env; \
	done
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
	@for svc in $(APP_SERVICES); do \
		img="${COMPANY_DOCKER_REGISTRY_HOST}/${COMPANY_HIRING_PROJECT_NAME}/$$svc:${COMPANY_HIRING_VERSION}"; \
		if docker image inspect $$img >/dev/null 2>&1; then \
			echo "  $(GREEN)✓$(NC) $$svc image: $$img"; \
		else \
			echo "  $(RED)✗$(NC) $$svc image: $$img - NOT FOUND"; \
		fi; \
	done
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
	@grep -E '^(env|env-infra|check-env):.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)Docker Compose:$(NC)"
	@grep -E '^(up|up-build|down|down-v|restart|logs|logs-db|logs-app|logs-auth|logs-blog|logs-notifications|logs-kafka|status|config):.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)Infrastructure:$(NC)"
	@grep -E '^(infra-up|infra-down):.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)Docker Build:$(NC)"
	@grep -E '^(docker-build|docker-build-auth|docker-build-blog|docker-build-notifications):.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)Run Services:$(NC)"
	@grep -E '^(run|run-auth|run-blog|run-notifications):.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)PostgreSQL Replication:$(NC)"
	@grep -E '^(pg-replication-status|pg-replication-lag|pg-test-replication):.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)Database Shells:$(NC)"
	@grep -E '^(pg-primary|pg-replica|pg-write|pg-read|redis-cli):.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)Kafka:$(NC)"
	@grep -E '^(kafka-topics|kafka-create-topic|kafka-describe-topic|kafka-consumer-groups):.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)Migrations:$(NC)"
	@grep -E '^(migrate-up|migrate-down|migrate-status|migrate-create):.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)Database Monitoring:$(NC)"
	@grep -E '^(pg-write-monitor|pg-read-monitor):.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)Database Query Logs:$(NC)"
	@grep -E '^(pg-write-log|pg-read-log|pg-write-log-full|pg-read-log-full):.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)Build & Utilities:$(NC)"
	@grep -E '^(build|clean|audit):.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help