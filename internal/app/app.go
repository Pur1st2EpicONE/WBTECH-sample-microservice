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
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/server"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/service"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/notifier"
	"github.com/jmoiron/sqlx"
)

type App struct { // orchestration layer
	logger          logger.Logger
	logFile         *os.File
	server          *server.Server
	consumer        broker.Consumer
	notifier        notifier.Notifier
	workers         int
	restartOnPanic  bool
	restartDelay    time.Duration
	cache           cache.Cache
	storage         repository.Storage
	dbCheckInterval time.Duration
	dbMaxChecks     int
	ctx             context.Context
	Stop            context.CancelFunc
	wg              *sync.WaitGroup
}

func Start() *App {

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

	notifier := notifier.NewNotifier(config.Notifier)
	server, cache, storage := wireApp(db, config, logger)
	wg := new(sync.WaitGroup)

	return &App{
		logger:          logger,
		logFile:         logFile,
		server:          server,
		consumer:        consumer,
		notifier:        notifier,
		workers:         config.Workers,
		restartOnPanic:  config.RestartOnPanic,
		restartDelay:    config.RestartDelay,
		cache:           cache,
		storage:         storage,
		dbCheckInterval: config.DbCheckInterval,
		dbMaxChecks:     config.DbMaxChecks,
		ctx:             ctx,
		Stop:            stop,
		wg:              wg}
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
		logger.LogInfo("app — received signal "+sig.String()+", initiating graceful shutdown", "layer", "app")
		cancel()
	}()
	return ctx, cancel
}

func (a *App) RunCacheCleaner() {
	dbStatus := make(chan bool, 1)
	go func() {
		var notified bool
		for {
			time.Sleep(a.dbCheckInterval)
			if err := a.storage.Ping(); err != nil {
				for range a.dbMaxChecks {
					if err = a.storage.Ping(); err != nil {
						time.Sleep(a.dbCheckInterval)
						continue
					}
				}
				if !notified {
					a.notifier.Notify("CRITICAL ERROR — database connection lost\ncache cleaner suspended\nservice is now in cache-only mode")
					notified = true
				}
				dbStatus <- false
			} else {
				notified = false
				dbStatus <- true
			}
		}
	}()
	a.cache.CacheCleaner(a.ctx, a.logger, dbStatus)
}

func (a *App) RunServer() {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		<-a.ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), a.server.ShutdownTimeout)
		defer cancel()
		a.server.Shutdown(ctx, a.logger)
	}()
	err := a.server.Run(a.ctx, a.logger)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.logger.LogFatal("app — server run failed", err, "layer", "app")
	}
}

func (a *App) RunConsumer() {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		<-a.ctx.Done()
		a.consumer.Close(a.logger)
	}()
	for workerID := 1; workerID < a.workers+1; workerID++ {
		a.wg.Add(1)
		go a.runWorker(workerID)
	}
}

func (a *App) runWorker(workerID int) {
	defer a.wg.Done()
	mu := new(sync.Mutex)
	for {
		select {
		case <-a.ctx.Done():
			a.logger.LogInfo(fmt.Sprintf("consumer — worker %d shutting down", workerID), "layer", "app")
			return
		default:
			func() {
				defer func() {
					if panicErr := recover(); panicErr != nil {
						a.logger.LogError(fmt.Sprintf("consumer — worker %d panicked", workerID), fmt.Errorf("%v", panicErr))
						if !a.restartOnPanic {
							a.logger.LogInfo(fmt.Sprintf("consumer — worker %d terminated", workerID), "layer", "app")
							mu.Lock()
							a.workers--
							if a.workers == 0 {
								a.logger.LogInfo("consumer — all workers have panicked, initiating emergency shutdown", "layer", "app")
								a.Stop()
							}
							mu.Unlock()
						}
					}
				}()
				a.consumer.Run(a.ctx, a.storage, a.logger, workerID)
			}()
			if !a.restartOnPanic {
				return
			}
			time.Sleep(a.restartDelay)
		}
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
