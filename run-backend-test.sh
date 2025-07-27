#!/bin/bash

export GO_ENV=test

# Build and run the Go application
echo "Building Go backend..."
go build -o battery_recycle main.go

echo "Starting backend server on port 8036..."
./battery_recycle