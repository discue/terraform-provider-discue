#!/bin/sh

set -eux

# Start test server in background and ensure it's killed on exit (success or failure)
(cd test-server && go run .) &
SERVER_PID=$!
trap 'kill "${SERVER_PID}" 2>/dev/null || true; wait "${SERVER_PID}" 2>/dev/null || true' EXIT

# Wait for server to become available
ready=0
for i in $(seq 1 50); do
    if command -v curl >/dev/null 2>&1; then
        if curl -sSf http://localhost:3000/ >/dev/null 2>&1; then ready=1; break; fi
        elif command -v nc >/dev/null 2>&1; then
        if nc -z localhost 3000 >/dev/null 2>&1; then ready=1; break; fi
    else
        sleep 0.1
    fi
    sleep 0.1
done

if [ "${ready}" -ne 1 ]; then
    echo "test server did not start" >&2
    kill "${SERVER_PID}" 2>/dev/null || true
    exit 1
fi

go test -v -cover ./internal/validators/
# Ensure provider uses the test server
export DISCUE_API_ENDPOINT="http://localhost:3000"
export DISCUE_API_KEY="test"

TF_ACC=1 go test -v -cover ./internal/provider/