#!/bin/bash

# Set environment variables for test database
export DB_HOST=182.92.150.161
export DB_PORT=3006
export DB_USER=root
export DB_PASSWORD="yoge@coder%%%123321!"
export DB_NAME=doc_ai
export SEED_DATABASE=true
export JWT_SECRET=your_jwt_secret_key

# Navigate to backend directory
cd backend

# Build and run the Go application
echo "Building Go backend..."
go build -o server cmd/server/main.go

echo "Starting backend server on port 8036..."
./server