FROM golang:1.20-alpine

WORKDIR /app

# Install git (needed for private repos sometimes)
RUN apk add --no-cache git

# Copy go.mod and go.sum first (for caching)
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o nalo_workspace ./cmd/api

EXPOSE 3000

CMD ["./nalo_workspace/cmd/api/"]
