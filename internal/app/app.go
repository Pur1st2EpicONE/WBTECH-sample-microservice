package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/config"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/kafka"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/server"
)

type App struct {
	srv      *server.Server
	consumer *kafka.Consumer
	storage  *repository.Storage
}

func New() *App {

	db, err := postgres.ConnectPostgres(config.Pg)
	if err != nil {
		logger.LogFatal("failed to connect to database: %v", err)
	}
	logger.LogInfo("postgres — connected to database")
	storage := repository.NewStorage(db)

	consumer, err := kafka.NewConsumer([]string{"localhost:9092"}, "orders", "orders")
	if err != nil {
		logger.LogFatal("failed to create consumer: %v", err)
	}

	srv := server.NewServer(config.HTTPPort, storage)

	return &App{srv: srv, consumer: consumer, storage: storage}
}

func (a *App) NewContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		logger.LogInfo("app — received signal: " + sig.String() + ", cancelling context")
		cancel()
	}()
	return ctx, cancel
}

func (a *App) RunServer(ctx context.Context) {
	go func() {
		<-ctx.Done()
		a.srv.Shutdown(context.Background())
		a.storage.Close()
	}()
	a.srv.Run(ctx)
}

func (a *App) RunConsumer(ctx context.Context) {
	go func() {
		<-ctx.Done()
		a.consumer.Close()
		logger.LogInfo("consumer stopped")
	}()
	a.consumer.Run(a.storage)
}

func (a *App) Wait(ctx context.Context) {
	<-ctx.Done()
	time.Sleep(1 * time.Second) // didn't want to clutter main with a wait group
}
