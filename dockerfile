FROM golang:1.24.0-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64

WORKDIR /app

# Copy go mod and sum files first for better cache utilization
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o worker_service cmd/main/main.go
# Run stage
FROM alpine:3.21

# Add ca certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/worker_service .

# Copy environment files
COPY .env.* ./

# Create a non-root user and set permissions
RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app
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