FROM golang:1.24.0-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -ldflags="-s -w" -o worker_service cmd/main/main.go

# Run stage
FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/worker_service .
EXPOSE 8888
ENTRYPOINT ["./worker_service"]