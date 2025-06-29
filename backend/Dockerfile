# ---- Build Stage ----
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Set Go proxy for faster downloads
ENV GOPROXY=https://goproxy.cn,direct

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/main ./cmd/server

# ---- Release Stage ----
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main /app/main

# Expose port
EXPOSE 8036

# Run the application
CMD ["/app/main"]