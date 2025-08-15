package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/cache"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/config"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/kafka"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/server"
)

type App struct {
	srv           *server.Server
	consumer      *kafka.Consumer
	storage       *repository.Storage
	ctx           context.Context
	CancelContext context.CancelFunc
}

func New() *App {

	db, err := postgres.ConnectPostgres(config.Pg)
	if err != nil {
		logger.LogFatal("app — failed to connect to database: %v", err)
	}
	storage := repository.NewStorage(db)

	consumer, err := kafka.NewConsumer([]string{"localhost:9092"}, "orders", "orders")
	if err != nil {
		logger.LogFatal("app — failed to create consumer: %v", err)
	}
	ctx, cancel := newContext()
	cache := cache.LoadCache(storage, 20*time.Second)
	go cache.CacheCleaner(ctx)
	srv := server.NewServer(config.HTTPPort, cache, storage)

	return &App{srv: srv, consumer: consumer, storage: storage, ctx: ctx, CancelContext: cancel}
}

func newContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		logger.LogInfo("app — received signal " + sig.String() + ", cancelling context")
		cancel()
	}()
	return ctx, cancel
}

func (a *App) RunServer() {
	go func() {
		<-a.ctx.Done()
		a.srv.Shutdown(context.Background())
		a.storage.Close()
	}()
	a.srv.Run(a.ctx)
}

func (a *App) RunConsumer() {
	go func() {
		<-a.ctx.Done()
		a.consumer.Close()
		logger.LogInfo("app — consumer stopped")
	}()
	a.consumer.Run(a.storage)
}

func (a *App) Wait() {
	<-a.ctx.Done()
	time.Sleep(1 * time.Second) // didn't want to clutter main with a wait group
}
