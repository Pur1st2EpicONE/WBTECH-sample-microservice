.PHONY: all create-topic run-service test test-unit test-integration coverage coverage-html migrate-test

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

test-unit:
	@cat .env.example > .env	
	@cat .env.example > ./internal/configs/.env
	@cat config.yaml > ./internal/configs/config.yaml	
	go test -v -cover -short ./...
	@rm -f coverage.out
	@rm -f coverage_filtered.out
	@rm -f coverage_res.out
	@rm -f ./internal/configs/.env
	@rm -f ./internal/configs/config.yaml	

test:
	@cat .env.example > .env	
	@cat .env.example > ./internal/configs/.env
	@cat config.yaml > ./internal/configs/config.yaml		
	@docker-compose up -d postgres_test kafka_test
	@until docker exec postgres_test pg_isready -U Neo > /dev/null 2>&1; do sleep 0.5; done
	@sleep 10
	@migrate -path ./schema -database 'postgres://Neo:0451@localhost:5434/wb-service-db-test?sslmode=disable' up
	@docker exec kafka_test /opt/kafka/bin/kafka-topics.sh --create --topic test-orders --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
	@go test -p 1 ./... -coverpkg=./... -coverprofile=coverage.out -v
	@grep -v "/mocks/" coverage.out > coverage_filtered.out
	@grep -v "/cmd/producer/" coverage_filtered.out > coverage_res.out
	@go tool cover -func=coverage_res.out
	@go tool cover -html=coverage_res.out -o cover.html
	@rm -f coverage.out
	@rm -f coverage_filtered.out
	@rm -f coverage_res.out
	@rm -f ./internal/configs/.env
	@rm -f ./internal/configs/config.yaml	
	-@docker-compose down

coverage:
	go test ./... -coverprofile=coverage.out -v
	go tool cover -html=coverage.out -o cover.html
	