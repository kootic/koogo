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
RUN CGO_ENABLED=0 GOOS=linux go build -o koogo ./cmd/koogo

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install CA certificates for HTTPS connections
RUN apk add --no-cache ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/koogo .

# Set the binary as executable
RUN chmod +x /app/koogo

# Expose the application port (adjust as needed)
EXPOSE 8080

# Run the application
CMD ["/app/koogo"]
