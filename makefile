include .env

PG_MIGRATIONS_PATH := migrations/postgres
PG_NETWORK := weather-by-geo_persistence-db
PG_DSN := postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):5432/$(POSTGRES_DB)?sslmode=disable

CH_MIGRATIONS_PATH := migrations/clickhouse
CH_NETWORK := weather-by-geo_analytic-db
CH_DSN := clickhouse://$(CLICKHOUSE_HOST):9000?username=$(CLICKHOUSE_USERNAME)&password=$(CLICKHOUSE_PASSWORD)&database=$(CLICKHOUSE_DB)

.PHONY: migrate-pg-up
migrate-pg-up:
	docker run --rm \
		-v $(PWD)/$(PG_MIGRATIONS_PATH):/migrations \
		--network $(PG_NETWORK) \
		migrate/migrate:latest \
		-path=/migrations \
		-database "$(PG_DSN)" \
		up

.PHONY: migrate-pg-down
migrate-pg-down:
	docker run --rm \
		-v $(PWD)/$(PG_MIGRATIONS_PATH):/migrations \
		--network $(PG_NETWORK) \
		migrate/migrate:latest \
		-path=/migrations \
		-database "$(PG_DSN)" \
		down 1

.PHONY: migrate-pg-status
migrate-pg-status:
	docker run --rm \
		-v $(PWD)/$(PG_MIGRATIONS_PATH):/migrations \
		--network $(PG_NETWORK) \
		migrate/migrate:latest \
		-path=/migrations \
		-database "$(PG_DSN)" \
		version

.PHONY: migrate-pg-create
migrate-pg-create:
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "Usage: make migrate-pg-create <migration_name>"; \
		exit 1; \
	fi
	docker run --rm \
		-v $(PWD)/$(PG_MIGRATIONS_PATH):/migrations \
		migrate/migrate:latest \
		create -ext sql -dir /migrations -seq $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-ch-up
migrate-ch-up:
	docker run --rm \
		-v $(PWD)/$(CH_MIGRATIONS_PATH):/migrations \
		--network $(CH_NETWORK) \
		migrate/migrate:latest \
		-path=/migrations \
		-database "$(CH_DSN)" \
		up

.PHONY: migrate-ch-down
migrate-ch-down:
	docker run --rm \
		-v $(PWD)/$(CH_MIGRATIONS_PATH):/migrations \
		--network $(CH_NETWORK) \
		migrate/migrate:latest \
		-path=/migrations \
		-database "$(CH_DSN)" \
		down 1

.PHONY: migrate-ch-status
migrate-ch-status:
	docker run --rm \
		-v $(PWD)/$(CH_MIGRATIONS_PATH):/migrations \
		--network $(CH_NETWORK) \
		migrate/migrate:latest \
		-path=/migrations \
		-database "$(CH_DSN)" \
		version

.PHONY: migrate-ch-create
migrate-ch-create:
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "Usage: make migrate-ch-create <migration_name>"; \
		exit 1; \
	fi
	docker run --rm \
		-v $(PWD)/$(CH_MIGRATIONS_PATH):/migrations \
		migrate/migrate:latest \
		create -ext sql -dir /migrations -seq $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up: migrate-pg-up migrate-ch-up

.PHONY: migrate-down
migrate-down: migrate-ch-down migrate-pg-down

%:
	@: