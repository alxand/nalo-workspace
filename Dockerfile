FROM golang:1.24-alpine

WORKDIR /app

# Install git (needed for private repos sometimes)
RUN apk add --no-cache git

# Copy go.mod and go.sum first (for caching)
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o nalo_workspace ./cmd/api/main.go

EXPOSE 3000

