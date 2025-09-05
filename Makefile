.PHONY: all up down orders bad-order local local-compose local-down db-load create-topic app migrate-up migrate-down test-unit test

all: up 

up:
	@cat .env.example > .env	
	@cp ./configs/config.full.yaml ./config.yaml
	@cp ./deployments/docker-compose.full.yaml ./docker-compose.full.yaml
	@cp ./deployments/Dockerfile ./Dockerfile
	@docker-compose -f docker-compose.full.yaml build --no-cache
	docker-compose -f docker-compose.full.yaml up -d
	@echo "wb-service is up"

down:
	docker-compose -f docker-compose.full.yaml down -v --remove-orphans
	@rm -f docker-compose.full.yaml
	@rm -f Dockerfile
	@rm -f config.yaml
	@rm -f .env

orders:
	@docker exec wb-service ./producer

bad-order:
	@docker exec wb-service ./producer bad

local: local-compose db-load migrate-up create-topic app

local-compose:
	@cat .env.example > .env
	@cp ./deployments/docker-compose.dev.yaml ./docker-compose.dev.yaml
	@cp ./configs/config.dev.yaml ./config.yaml
	@docker-compose -f docker-compose.dev.yaml up -d postgres kafka
	@echo "Waiting for Kafka to start..."
	@sleep 7

local-down:
	@docker-compose -f docker-compose.dev.yaml down
	@rm -rf ./logs
	@rm -f config.yaml
	@rm -f docker-compose.dev.yaml
	@rm -f .env
	
db-load:
	@until docker exec postgres pg_isready -U Neo > /dev/null 2>&1; do sleep 0.5; done

create-topic:
	docker exec kafka /opt/kafka/bin/kafka-topics.sh --create --topic orders --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
	@echo "Topic created"

app:
	go run ./cmd/wb-service/main.go -o wb-service

migrate-up:
	@migrate -path ./schema -database 'postgres://Neo:0451@localhost:5433/wb-service-db?sslmode=disable' up

migrate-down:
	@migrate -path ./schema -database 'postgres://Neo:0451@localhost:5433/wb-service-db?sslmode=disable' down

test-unit:
	@cat .env.example > .env	
	@cat .env.example > ./internal/configs/.env
	@cp ./configs/config.dev.yaml .
	@cat ./config.dev.yaml > ./internal/configs/config.yaml	
	go test -v -cover -short ./...
	@rm -f coverage.out
	@rm -f coverage_filtered.out
	@rm -f coverage_res.out
	@rm -f ./internal/configs/.env
	@rm -f ./internal/configs/config.yaml	
	@rm -f ./config.dev.yaml 	

test:
	@cat .env.example > .env	
	@cat .env.example > ./internal/configs/.env
	@cp ./deployments/docker-compose.dev.yaml ./docker-compose.dev.yaml
	@cp ./configs/config.dev.yaml ./config.yaml
	@cat config.yaml > ./internal/configs/config.yaml		
	@docker-compose -f docker-compose.dev.yaml up -d postgres_test kafka_test
	@until docker exec postgres_test pg_isready -U Neo > /dev/null 2>&1; do sleep 0.5; done
	@sleep 10
	@migrate -path ./schema -database 'postgres://Neo:0451@localhost:5434/wb-service-db-test?sslmode=disable' up
	@docker exec kafka_test /opt/kafka/bin/kafka-topics.sh --create --topic test-orders --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
	@go test -p 1 ./... -coverpkg=./... -coverprofile=coverage.out -v
	@grep -v "/mocks/" coverage.out > coverage_filtered.out
	@grep -v "/cmd/producer/" coverage_filtered.out > coverage_res.out
	@go tool cover -func=coverage_res.out
	@go tool cover -html=coverage_res.out -o code_coverage.html
	@rm -f coverage.out
	@rm -f coverage_filtered.out
	@rm -f coverage_res.out
	@rm -f ./internal/configs/.env
	@rm -f ./internal/configs/config.yaml	
	@docker-compose -f docker-compose.dev.yaml down
	@rm -f ./docker-compose.dev.yaml up	
	@rm -f ./config.yaml

local-orders:
	go run ./cmd/producer

local-bad-order:
	go run ./cmd/producer bad
