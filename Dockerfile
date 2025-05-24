FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o dailylog ./cmd/dailylog

EXPOSE 3000

CMD ["./nalo_workspace/cmd/api/"]
