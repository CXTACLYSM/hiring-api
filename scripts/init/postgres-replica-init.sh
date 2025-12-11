#!/bin/bash
# =============================================================================
# PostgreSQL REPLICA - INITIALIZATION SCRIPT
# =============================================================================
# This script runs as the entrypoint for the replica container
# It initializes from primary via pg_basebackup if needed
# =============================================================================

set -e

# Check if data directory is already initialized
if [ -s /var/lib/postgresql/data/PG_VERSION ]; then
echo "Replica: Data directory already initialized, starting PostgreSQL..."
else
echo "Replica: Initializing from primary via pg_basebackup..."

# Wait for primary to be ready
until PGPASSWORD="${REPLICATOR_PASSWORD}" pg_isready -h "${PRIMARY_HOST}" -U "${REPLICATOR_USER}"; do
echo "Waiting for primary at ${PRIMARY_HOST}..."
sleep 2
done

echo "Primary is ready, starting base backup..."

# Perform base backup
# -h: primary host
# -U: replication user
# -D: data directory
# -Fp: plain format
# -Xs: stream WAL during backup
# -P: show progress
# -R: create standby.signal and postgresql.auto.conf
# -S: use replication slot
PGPASSWORD="${REPLICATOR_PASSWORD}" pg_basebackup \
-h "${PRIMARY_HOST}" \
-U "${REPLICATOR_USER}" \
-D /var/lib/postgresql/data \
-Fp -Xs -P -R \
-S replica_slot

# Copy our custom configs
cp /etc/postgresql/postgresql.conf /var/lib/postgresql/data/
cp /etc/postgresql/pg_hba.conf /var/lib/postgresql/data/

# Fix permissions
chown -R postgres:postgres /var/lib/postgresql/data
chmod 700 /var/lib/postgresql/data

echo "Replica: Base backup completed!"
fi

# Start PostgreSQL as postgres user
exec su-exec postgres postgres \
-c config_file=/var/lib/postgresql/data/postgresql.conf \
-c hba_file=/var/lib/postgresql/data/pg_hba.conf