#!/usr/bin/env bash
set -euo pipefail

# Simple dev helper: starts a Postgres container (if not running), runs server locally.
# Usage: ./scripts/dev.sh

PROJECT_NAME=resume-app
DB_NAME=resumeapp
DB_CONTAINER=resume-db
DB_PORT=5432
DB_PASSWORD=postgres

function cleanup() {
  echo "Cleaning up..."
  if docker ps -a --format '{{.Names}}' | grep -q "^${DB_CONTAINER}$"; then
    echo "Stopping ${DB_CONTAINER}..."
    docker stop ${DB_CONTAINER} >/dev/null 2>&1 || true
    docker rm ${DB_CONTAINER} >/dev/null 2>&1 || true
  fi
}

# On exit, do not remove DB automatically to preserve data, but stop when asked.
trap 'echo Exiting dev.sh' EXIT

# Start Postgres if not present
if ! docker ps --format '{{.Names}}' | grep -q "^${DB_CONTAINER}$"; then
  echo "Starting Postgres container '${DB_CONTAINER}'..."
  docker run -d --name ${DB_CONTAINER} -e POSTGRES_PASSWORD=${DB_PASSWORD} -e POSTGRES_DB=${DB_NAME} -p ${DB_PORT}:5432 postgres:16-alpine
else
  echo "Postgres container '${DB_CONTAINER}' already running"
fi

# Export env for local server run
export DB_HOST=127.0.0.1
export DB_PORT=${DB_PORT}
export DB_USER=postgres
export DB_PASSWORD=${DB_PASSWORD}
export DB_NAME=${DB_NAME}

echo "Building and running server locally..."
# Build to catch compile errors first
if ! go build -o bin/server ./cmd/server; then
  echo "Build failed"
  exit 1
fi

# Run server in background and print PID
nohup ./bin/server > server.log 2>&1 &
PID=$!
sleep 1
echo "Server started with PID ${PID}. Logs: server.log"

# Show tail of logs
sleep 1
tail -n 120 server.log || true

echo "You can curl http://localhost:8080/health"

echo "To stop: kill ${PID} && docker stop ${DB_CONTAINER}"
