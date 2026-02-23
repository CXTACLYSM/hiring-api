#!/bin/bash

set -e

BOOTSTRAP_SERVER="${KAFKA_HOST}:${KAFKA_PORT}"
KAFKA_BIN="/opt/kafka/bin"

echo "Waiting for Kafka at ${BOOTSTRAP_SERVER}..."
for i in $(seq 1 60); do
    if ${KAFKA_BIN}/kafka-topics.sh --bootstrap-server "${BOOTSTRAP_SERVER}" --list >/dev/null 2>&1; then
        echo "Kafka is ready."
        break
    fi
    if [ "$i" -eq 60 ]; then
        echo "ERROR: Kafka failed to start within 60 seconds."
        exit 1
    fi
    sleep 2
done

create_topic() {
    local topic=$1
    local partitions=$2
    local replication=$3

    if ${KAFKA_BIN}/kafka-topics.sh --bootstrap-server "${BOOTSTRAP_SERVER}" --describe --topic "${topic}" >/dev/null 2>&1; then
        echo "Topic '${topic}' already exists, skipping."
    else
        echo "Creating topic '${topic}' (partitions=${partitions}, replication=${replication})..."
        ${KAFKA_BIN}/kafka-topics.sh --bootstrap-server "${BOOTSTRAP_SERVER}" \
            --create \
            --topic "${topic}" \
            --partitions "${partitions}" \
            --replication-factor "${replication}"
        echo "Topic '${topic}' created."
    fi
}

create_topic "user.created" 3 1

echo ""
echo "=== Current topics ==="
${KAFKA_BIN}/kafka-topics.sh --bootstrap-server "${BOOTSTRAP_SERVER}" --list
echo ""
echo "Kafka topic initialization complete."