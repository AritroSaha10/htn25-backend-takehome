# Build stage
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# We need CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Production stage
FROM debian:bookworm-slim
WORKDIR /app

# Install SQLite dependencies and CA certificates
RUN apt-get update && apt-get install -y libsqlite3-0 ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy binary from builder
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
