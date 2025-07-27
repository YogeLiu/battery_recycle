#!/bin/bash

# Set environment variables for test database
export GO_ENV=prod

# Build and run the Go application
echo "Building Go backend..."
go build -o battery_recycle main.go

echo "Starting backend server on port 8036..."
./battery_recycle