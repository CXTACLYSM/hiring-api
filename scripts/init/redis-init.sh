#!/bin/sh
set -e

REDIS_CONF="/usr/local/etc/redis/redis.conf"

redis-server "${REDIS_CONF}" &
REDIS_PID=$!

echo "Waiting for Redis to start..."
for i in $(seq 1 30); do
    if redis-cli PING 2>/dev/null | grep -q PONG; then
        echo "Redis is ready."
        break
    fi
    if [ "$i" -eq 30 ]; then
        echo "ERROR: Redis failed to start within 30 seconds."
        exit 1
    fi
    sleep 1
done

echo "Creating application user: ${REDIS_USERNAME}"
redis-cli ACL SETUSER "${REDIS_USERNAME}" on ">${REDIS_PASSWORD}" "~*" "+@all"

redis-cli ACL SAVE
echo "ACL configuration saved."

redis-cli ACL SETUSER default on ">${REDIS_PASSWORD}" ~* +@all

if redis-cli --user "${REDIS_USERNAME}" -a "${REDIS_PASSWORD}" --no-auth-warning PING | grep -q PONG; then
    echo "Application user verified successfully."
else
    echo "ERROR: Cannot connect as application user!"
    exit 1
fi

wait $REDIS_PID