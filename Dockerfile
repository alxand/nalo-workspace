# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and curl for health checks
RUN apk add --no-cache git curl

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./

RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nalo-workspace ./cmd/api/main.go

# Final stage
FROM alpine:latest

# Install curl for health checks
RUN apk --no-cache add curl ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/nalo-workspace .

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:3000/health || exit 1

# Run the application
CMD ["./nalo-workspace"]

