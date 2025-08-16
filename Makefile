.PHONY: all up create-topic run-service

all: services db-load migrate-up create-topic run-service

services:
	@cat .env.example > .env
	@docker-compose up -d
	@echo "Waiting for Kafka to start..."
	@sleep 10

down:
	@docker-compose down
	@rm -rf ./logs
	@rm -rf .env

db-load:
	@until docker exec postgres pg_isready -U Neo > /dev/null 2>&1; do sleep 0.5; done

create-topic:
	docker exec kafka /opt/kafka/bin/kafka-topics.sh --create --topic orders --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
	@echo "Topic created"

run-service:
	go run ./cmd/wb-service/main.go

migrate-up:
	@migrate -path ./schema -database 'postgres://Neo:0451@localhost:5433/wb-service-db?sslmode=disable' up

migrate-down:
	@migrate -path ./schema -database 'postgres://Neo:0451@localhost:5433/wb-service-db?sslmode=disable' down
