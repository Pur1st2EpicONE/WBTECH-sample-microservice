.PHONY: all up create-topic run-service test test-unit test-integration coverage coverage-html migrate-test

all: services db-load migrate-up create-topic run-service

services:
	@cat .env.example > .env
	@docker-compose up -d postgres kafka
	@echo "Waiting for Kafka to start..."
	@sleep 7
	
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

test: test-pg

test-unit:
	go test -v -cover ./...

test-pg:
	@docker-compose up -d postgres_test
	@until docker exec postgres_test pg_isready -U Neo > /dev/null 2>&1; do sleep 0.5; done
	@sleep 10
	@migrate -path ./schema -database 'postgres://Neo:0451@localhost:5434/wb-service-db-test?sslmode=disable' up
	@go test ./... -coverpkg=./... -coverprofile=coverage.out -v
	@go tool cover -html=coverage.out -o cover.html
	@grep -v "/mocks/" coverage.out > coverage_filtered.out
	@go tool cover -func=coverage_filtered.out
	@rm -f coverage.out
	@rm -f coverage_filtered.out
	-@docker-compose down

coverage:
	go test ./... -coverprofile=coverage.out -v
	go tool cover -html=coverage.out -o cover.html
	