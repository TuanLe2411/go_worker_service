FROM golang:1.24.0-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app

# Copy go mod and sum files first for better cache utilization
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -ldflags="-s -w" -o worker_service cmd/main/main.go
# Run stage
FROM alpine:3.21

# Add necessary runtime packages and security configurations
RUN apk add --no-cache ca-certificates tzdata curl && \
    addgroup -S appgroup && adduser -S appuser -G appgroup && \
    mkdir -p /app/log && chown appuser:appgroup /app/log && chmod 777 /app/log

# Set the timezone to Asia/Ho_Chi_Minh
ENV TZ=Asia/Ho_Chi_Minh
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/worker_service .

# Copy environment files
COPY .env.* ./

USER appuser

# Expose the port
EXPOSE 8888

# Set environment variable
ENV ENV=production

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8888/health || exit 1

# Run the binary
ENTRYPOINT ["./worker_service"]