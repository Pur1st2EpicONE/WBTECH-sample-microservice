package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/broker"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/cache"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/config"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/server"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type App struct {
	srv      *server.Server
	consumer *broker.Consumer
	storage  *repository.Storage
	ctx      context.Context
	Stop     context.CancelFunc
	wg       *sync.WaitGroup
}

func Start() *App {

	if err := initConfigs(); err != nil {
		logger.LogFatal("app — failed to load config", err)
	}

	db, err := postgres.ConnectPostgres(postgres.GetConfig())
	if err != nil {
		logger.LogFatal("app — failed to connect to database", err)
	}
	storage := repository.NewStorage(db)

	consumer, err := broker.NewConsumer([]string{"localhost:9092"}, "orders", "orders")
	if err != nil {
		logger.LogFatal("app — failed to create consumer", err)
	}
	var wg sync.WaitGroup
	ctx, stop := newContext()
	cache := cache.LoadCache(storage, 20*time.Second)
	go cache.CacheCleaner(ctx)
	srv := server.NewServer(config.HTTPPort, cache, storage)

	return &App{srv: srv, consumer: consumer, storage: storage, ctx: ctx, Stop: stop, wg: &wg}
}

func initConfigs() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func newContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background()) // can't get signal info with signal.NotifyContext
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
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		<-a.ctx.Done()
		a.srv.Shutdown(a.ctx)
		a.storage.Close()
	}()
	a.srv.Run(a.ctx)
}

func (a *App) RunConsumer() {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		<-a.ctx.Done()
		a.consumer.Close()
	}()
	a.consumer.Run(a.ctx, a.storage)
}

func (a *App) Wait() {
	<-a.ctx.Done()
	a.wg.Wait()
}
