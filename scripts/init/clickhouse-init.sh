#!/bin/bash
set -e

clickhouse-client --user "$CLICKHOUSE_USER" --password "$CLICKHOUSE_PASSWORD" <<-EOSQL
    -- Create database if not exists
    CREATE DATABASE IF NOT EXISTS ${CLICKHOUSE_DB:-weather_by_geo};

    SELECT 'ClickHouse initialized successfully';
EOSQL