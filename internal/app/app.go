package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/broker"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/cache"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/handler"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/server"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/service"
	"github.com/jmoiron/sqlx"
)

type App struct { // orchestration layer
	logger       logger.Logger
	logFile      *os.File
	server       *server.Server
	consumer     broker.Consumer
	workers      int
	restartDelay time.Duration
	cache        cache.Cache
	storage      repository.Storage
	ctx          context.Context
	Stop         context.CancelFunc
	wg           *sync.WaitGroup
}

func Start() *App {
	var wg sync.WaitGroup

	config, err := configs.Load()
	if err != nil {
		log.Fatalf("app — failed to load configs: %v", err)
	}

	logger, logFile := logger.NewLogger(config.Logger)

	ctx, stop := newContext(logger)

	db, err := repository.ConnectDB(config.Database)
	if err != nil {
		logger.LogFatal("app — failed to connect to database", err, "layer", "app")
	}
	logger.LogInfo("app — connected to database", "layer", "app")

	consumer, err := broker.NewConsumer(config.Consumer, logger)
	if err != nil {
		logger.LogFatal("app — failed to create consumer", err, "layer", "app")
	}

	server, cache, storage := wireApp(db, config, logger)

	return &App{
		logger:       logger,
		logFile:      logFile,
		server:       server,
		consumer:     consumer,
		workers:      config.Workers,
		restartDelay: config.RestartDelay,
		cache:        cache,
		storage:      storage,
		ctx:          ctx,
		Stop:         stop,
		wg:           &wg}
}

func wireApp(db *sqlx.DB, config configs.App, logger logger.Logger) (*server.Server, cache.Cache, repository.Storage) {
	storage := repository.NewStorage(db, logger)
	cache := cache.NewCache(storage, config.Cache, logger)
	service := service.NewService(storage, cache)
	handler := (handler.NewHandler(service, logger)).InitRoutes()
	server := server.NewServer(config.Server, handler)
	return server, cache, storage
}

func newContext(logger logger.Logger) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background()) // can't get signal info with signal.NotifyContext
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		logger.LogInfo("app — received signal "+sig.String()+", cancelling context", "layer", "app")
		cancel()
	}()
	return ctx, cancel
}

func (a *App) RunCacheCleaner() {
	a.cache.CacheCleaner(a.ctx, a.logger)
}

func (a *App) RunServer() {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		<-a.ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		a.server.Shutdown(ctx, a.logger)
	}()
	err := a.server.Run(a.ctx, a.logger)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.logger.LogFatal("app — server run failed", err, "layer", "app")
	}
}

func (a *App) RunConsumer() {
	go func() {
		<-a.ctx.Done()
		a.consumer.Close(a.logger)
	}()
	for i := 1; i < a.workers+1; i++ {
		a.wg.Add(1)
		go a.runWorker(i)
	}
}

func (a *App) runWorker(workerID int) {
	defer a.wg.Done()
	for {
		func() {
			defer func() {
				if panicErr := recover(); panicErr != nil {
					a.logger.LogError(fmt.Sprintf("worker %d panicked, restarting", workerID), fmt.Errorf("%v", panicErr))
				}
			}()
			a.consumer.Run(a.ctx, a.storage, a.logger, workerID)
		}()
		if a.ctx.Err() != nil {
			return
		}
		time.Sleep(a.restartDelay)
	}
}

func (a *App) Wait() {
	<-a.ctx.Done()
	a.wg.Wait()
	a.storage.Close()
	if a.logFile != nil {
		a.logFile.Close()
	}
}
