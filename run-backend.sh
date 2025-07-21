#!/bin/bash

# Set environment variables for test database
export DB_HOST=182.92.150.161
export DB_PORT=3006
export DB_USER=root
export DB_PASSWORD="yoge@coder%%%123321!"   
export DB_NAME=battery_recycle_erp
export SEED_DATABASE=true
export JWT_SECRET=your_jwt_secret_key


# Build and run the Go application
echo "Building Go backend..."
go build -o battery_recycle main.go

echo "Starting backend server on port 8036..."
./battery_recycle