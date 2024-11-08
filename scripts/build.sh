#!/bin/bash

# Exit on any error
set -e

# Get the directory of the script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${DIR}/.."

# Ensure dependencies are up to date
go mod tidy

# Run tests
go test ./...

# Build the binary
go build -o bin/metrics-agent ./cmd/agent/main.go

echo "Build complete: bin/metrics-agent"