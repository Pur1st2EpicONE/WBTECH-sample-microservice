package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/kafka"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/repository/postgres"
)

func main() {
	consumer, err := kafka.NewConsumer([]string{"localhost:9092"}, "orders", "orders")
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}
	config := postgres.PgConfig{
		Host:     "localhost",
		Port:     "5433",
		Username: "postgres",
		Password: "qwerty",
		DBName:   "postgres",
		SSLMode:  "disable",
	}
	db, err := postgres.ConnectPostgres(config)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	storage := repository.NewStorage(db)
	go func() {
		if err := consumer.Run(*storage); err != nil {
			log.Fatalf("Consumer failed: %v", err)
		}
	}()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
	close(sigCh)
}
