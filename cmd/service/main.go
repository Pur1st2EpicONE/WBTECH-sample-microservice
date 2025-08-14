package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/kafka"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/server"
)

func main() {

	storage, err := initProject()
	if err != nil {
		logFatal("project init failed", err)
	}
	srv := server.InitServer("8080", storage)
	go func() {
		err := srv.Run()
		if err != nil && err != http.ErrServerClosed {
			logFatal("server run failed", err)
		}
	}()

	consumer, err := kafka.NewConsumer([]string{"localhost:9092"}, "orders", "orders")
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}
	go func() {
		if err := consumer.Run(storage); err != nil {
			log.Fatalf("consumer failed: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
	close(sigCh)

	if err := srv.Shutdown(context.Background()); err != nil {
		logFatal("server shutdown fail", err)
	}

	if err := storage.Close(); err != nil {
		logFatal("db connection failed to close properly", err)
	}
}

func initProject() (*repository.Storage, error) {

	logger := initLogger()
	slog.SetDefault(logger)

	db, err := postgres.ConnectPostgres(postgres.PgConfig{
		Host:     "localhost",
		Port:     "5433",
		Username: "postgres",
		Password: "qwerty",
		DBName:   "postgres",
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	storage := repository.NewStorage(db)
	return storage, nil
}

func initLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func logFatal(msg string, err error) {
	slog.Error(msg, slog.String("err", err.Error()))
	os.Exit(1)
}
