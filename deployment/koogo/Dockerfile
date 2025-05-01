# Build stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o koogo ./cmd/koogo

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install CA certificates for HTTPS connections and Atlas for database migrations
RUN apk add --no-cache ca-certificates curl
ENV ATLAS_VERSION=v0.32.0
RUN curl -sSf https://atlasgo.sh | sh -s -- --community

# Copy the binary from the builder stage
COPY --from=builder /app/koogo .

# Copy the migrations from the builder stage
COPY --from=builder /app/internal/repo/postgres/migrations ./migrations

# Set the binary as executable
RUN chmod +x /app/koogo

# Expose the application port (adjust as needed)
EXPOSE 8080

# Run the application
ENTRYPOINT ["/app/koogo"]

CMD ["start"]
