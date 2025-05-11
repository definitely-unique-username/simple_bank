#!/bin/sh

set -e 

echo "Running DB migration"
source /app/.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Starting"
exec "$@"