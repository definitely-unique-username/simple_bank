#!/bin/sh

set -e 

echo "Running DB migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Starting"
exec "$@"