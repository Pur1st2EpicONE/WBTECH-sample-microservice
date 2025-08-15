package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/config"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/kafka"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/server"
)

type App struct {
	srv          *server.Server
	consumer     *kafka.Consumer
	storage      *repository.Storage
	closeLogFile func() error
}

func New() *App {

	db, err := postgres.ConnectPostgres(config.Pg)
	if err != nil {
		logger.LogFatal("failed to connect db: %v", err)
	}
	storage := repository.NewStorage(db)

	consumer, err := kafka.NewConsumer([]string{"localhost:9092"}, "orders", "orders")
	if err != nil {
		logger.LogFatal("failed to create consumer: %v", err)
	}

	srv := server.NewServer(config.HTTPPort, storage)

	return &App{srv: srv, consumer: consumer, storage: storage}
}

func (a *App) RunServer() {
	logger.LogInfo("app — starting the server")
	a.srv.Run()
}

func (a *App) RunConsumer() {
	if err := a.consumer.Run(a.storage); err != nil {
		logger.LogFatal("consumer failed: %v", err)
	}
}

func (a *App) WaitForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh

	a.srv.Shutdown(context.Background())
	a.storage.Close()
}
