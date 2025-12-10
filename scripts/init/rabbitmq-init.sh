#!/bin/bash
set -e

until rabbitmqctl status > /dev/null 2>&1; do
    echo "Waiting for RabbitMQ..."
    sleep 2
done

rabbitmqadmin -u rabbitmq -p rabbitmq -V weather declare exchange name=weather.events type=topic durable=true
rabbitmqadmin -u rabbitmq -p rabbitmq -V weather declare exchange name=weather.dlx type=direct durable=true

rabbitmqadmin -u rabbitmq -p rabbitmq -V weather declare queue name=weather.notifications durable=true arguments='{"x-dead-letter-exchange":"weather.dlx","x-dead-letter-routing-key":"weather.notifications.dlq"}'
rabbitmqadmin -u rabbitmq -p rabbitmq -V weather declare queue name=weather.analytics durable=true arguments='{"x-dead-letter-exchange":"weather.dlx","x-dead-letter-routing-key":"weather.analytics.dlq"}'

rabbitmqadmin -u rabbitmq -p rabbitmq -V weather declare queue name=weather.notifications.dlq durable=true
rabbitmqadmin -u rabbitmq -p rabbitmq -V weather declare queue name=weather.analytics.dlq durable=true

rabbitmqadmin -u rabbitmq -p rabbitmq -V weather declare binding source=weather.events destination=weather.notifications routing_key="notification.#"
rabbitmqadmin -u rabbitmq -p rabbitmq -V weather declare binding source=weather.events destination=weather.analytics routing_key="analytics.#"
rabbitmqadmin -u rabbitmq -p rabbitmq -V weather declare binding source=weather.dlx destination=weather.notifications.dlq routing_key="weather.notifications.dlq"
rabbitmqadmin -u rabbitmq -p rabbitmq -V weather declare binding source=weather.dlx destination=weather.analytics.dlq routing_key="weather.analytics.dlq"

echo "RabbitMQ initialized successfully"