all: services db-load migrate-up create-topic run-service

services:
	@docker-compose up -d
	@echo "Waiting for Kafka to start..."
	@sleep 10

db-load:
	@until docker exec postgres pg_isready -U postgres > /dev/null 2>&1; do sleep 0.5; done

create-topic:
	docker exec kafka /opt/kafka/bin/kafka-topics.sh --create --topic orders --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
	@echo "Topic created"

run-service:
	@go run ./cmd/service/main.go

migrate-up:
	@migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5433/postgres?sslmode=disable' up

migrate-down:
	@migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5433/postgres?sslmode=disable' down
