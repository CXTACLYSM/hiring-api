#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Extensions
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    CREATE EXTENSION IF NOT EXISTS "pgcrypto";

    -- Schemas (optional)
    -- CREATE SCHEMA IF NOT EXISTS app;

    RAISE NOTICE 'PostgreSQL initialized successfully';
EOSQL