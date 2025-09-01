FROM golang:1.24-bullseye AS builder

RUN apt-get update && apt-get install -y \
    git \
    build-essential \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN cat .env
COPY .env .
COPY config.yaml .

RUN go build -o wb-service ./cmd/wb-service
RUN go build -o producer ./cmd/producer/main.go 

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/wb-service .
COPY --from=builder /app/.env .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/producer . 

EXPOSE 8080

CMD ["./wb-service"]
