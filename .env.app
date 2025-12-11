# =============================================================================
# APPLICATION CONFIG
# =============================================================================
# These variables are used by Go application
# =============================================================================

# -----------------------------------------------------------------------------
# APP
# -----------------------------------------------------------------------------
APP_HOST=0.0.0.0
APP_PORT=8080
JWT_SECRET=

# -----------------------------------------------------------------------------
# POSTGRESQL - WRITE (Commands)
# -----------------------------------------------------------------------------
POSTGRES_WRITE_HOST=hiring_primary
POSTGRES_WRITE_PORT=5432
POSTGRES_WRITE_USERNAME=company_pg_write_user
POSTGRES_WRITE_PASSWORD=company_pg_write_password
POSTGRES_WRITE_DB=hiring

# -----------------------------------------------------------------------------
# POSTGRESQL - READ (Queries)
# -----------------------------------------------------------------------------
POSTGRES_READ_HOST=hiring_replica
POSTGRES_READ_PORT=5432
POSTGRES_READ_USERNAME=company_pg_read_user
POSTGRES_READ_PASSWORD=company_pg_read_password
POSTGRES_READ_DB=hiring

# -----------------------------------------------------------------------------
# REDIS
# -----------------------------------------------------------------------------
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=redis_password
REDIS_AUTH_DB=0
REDIS_RESOURCE_DB=1

# -----------------------------------------------------------------------------
# RABBITMQ
# -----------------------------------------------------------------------------
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_USERNAME=rabbitmq
RABBITMQ_PASSWORD=rabbitmq_password
RABBITMQ_VHOST=hiring