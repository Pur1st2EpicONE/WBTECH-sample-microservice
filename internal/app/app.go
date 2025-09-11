/*
Package app provides the core layer of the service.

It defines the main application instance that orchestrates
all critical components, including the HTTP server, message broker consumer,
critical error notifier, cache, storage, and logging.
*/
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
	"sync/atomic"
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

/*
App structure is the core instance of the project,
coordinating initialization, runtime operations, and graceful shutdown.
*/
type App struct {
	logger          logger.Logger      // main logger for the service
	logFile         *os.File           // output for logs, can be a file or stdout
	server          *server.Server     // HTTP server instance
	consumer        broker.Consumer    // consumer for processing orders
	notifier        notifier.Notifier  // notifies about critical errors
	workers         int                // number of worker goroutines for message processing
	restartOnPanic  bool               // whether workers should restart on panic
	restartDelay    time.Duration      // delay before restarting a worker after panic
	cache           cache.Cache        // application cache for storing orders
	storage         repository.Storage // database storage interface
	dbCheckInterval time.Duration      // interval between DB connectivity checks
	dbMaxChecks     int                // max number of failed DB checks before action
	ctx             context.Context    // root context for graceful shutdown
	Stop            context.CancelFunc // cancels the root context
	wg              *sync.WaitGroup    // tracks background goroutines
}

/*
Start initializes the application and wires together all the main components.

It performs the following steps:
 1. Loads application configuration (database, server, cache, consumer, etc.).
 2. Sets up logging (file/stdout).
 3. Creates a root context with cancellation for graceful shutdown.
 4. Connects to the database and checks connectivity.
 5. Initializes the message broker consumer.
 6. Sets up a notifier to report critical errors.
 7. Wires dependencies: repository, cache, service, HTTP handlers, and server.
 8. Returns a fully configured App instance ready to run.
*/
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
		wg:              wg,
	}
}

/*
newContext creates a root context that handles graceful shutdown.

Instead of using signal.NotifyContext, it sets up a custom signal handler
for SIGINT and SIGTERM. When such a signal is received, it:
  - Logs the shutdown signal.
  - Cancels the context to notify all running goroutines.

This ensures that shutdown signals are logged and all components terminate gracefully.
*/
func newContext(logger logger.Logger) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		logger.LogInfo("app — received signal "+sig.String()+", initiating graceful shutdown", "layer", "app")
		cancel()
	}()
	return ctx, cancel
}

/*
wireApp performs dependency injection for the application.

It wires together the core layers — storage, cache, service, HTTP handler,
and server — ensuring all components are properly constructed and connected.

Returns the fully initialized server, cache, and storage instances.
*/
func wireApp(db *sqlx.DB, config configs.App, logger logger.Logger) (*server.Server, cache.Cache, repository.Storage) {
	storage := repository.NewStorage(db, logger)
	cache := cache.NewCache(storage, config.Cache, logger)
	service := service.NewService(storage, cache)
	handler := (handler.NewHandler(service, logger)).InitRoutes()
	server := server.NewServer(config.Server, handler)
	return server, cache, storage
}

/*
RunCacheCleaner launches a background cleanup process that links database health monitoring
with cache management.

  - Continuously checks DB availability in the background.
  - Switches to "cache-only mode" and sends an alert if the DB is unreachable.
  - Restores normal operation and re-enables cleanup when the DB recovers.
  - Provides DB status updates to the cache cleaner.

Notes:
  - The cleanup runs in the background and removes only cache entries
    that haven’t been accessed for a long time.
  - Cache overflow itself is prevented by a ring buffer,
    so cleanup is an additional mechanism to keep the cache fresh.
  - The cleanup mechanism itself can be enabled or disabled through the service configuration.
*/
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
					_ = a.notifier.Notify("CRITICAL ERROR — database connection lost\ncache cleaner suspended\nservice is now in cache-only mode")
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

/*
RunServer starts the HTTP server and listens for shutdown signals.

It works as follows:
  - Launches a goroutine that listens for context cancellation
    and gracefully shuts down the server with a timeout if a signal is received.
  - Runs the server (blocking).
  - If the server exits unexpectedly (not due to shutdown),
    logs a fatal error and terminates the application.
*/
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

/*
RunConsumer acts as a system monitor for message-processing workers.

  - Starts a supervisory goroutine that shuts down the consumer when the context is cancelled.
  - Spawns multiple worker goroutines to process messages concurrently.
  - Each worker is monitored: panics and failures are reported and handled according to the restart policy defined in the configuration.
*/
func (a *App) RunConsumer() {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		<-a.ctx.Done()
		a.consumer.Close(a.logger)
	}()
	var lastWorker atomic.Int32
	for workerID := 1; workerID < a.workers+1; workerID++ {
		a.wg.Add(1)
		lastWorker.Add(1)
		go a.runWorker(workerID, &lastWorker)
	}
}

/*
runWorker executes a single worker instance for message processing.

The worker continuously polls messages from the broker and handles any panics.
Errors are logged, and the worker is either restarted or terminated according to
the configured restart policy. If all workers terminate without restart, an emergency
shutdown is triggered.
*/
func (a *App) runWorker(workerID int, lastWorker *atomic.Int32) {
	defer a.wg.Done()
	mu := new(sync.Mutex)
	for {
		select {
		case <-a.ctx.Done():
			a.logger.LogInfo(fmt.Sprintf("consumer — worker %d stopped", workerID), "layer", "app")
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
				a.consumer.Run(a.ctx, a.storage, a.logger, workerID, lastWorker)
			}()
			if !a.restartOnPanic {
				return
			}
			time.Sleep(a.restartDelay)
		}
	}
}

/*
Wait blocks until all application components shut down.

Steps:
 1. Waits for the root context cancellation (ctx acts as a blocking point to prevent premature main exit).
 2. Waits for all goroutines (server, consumer, workers, cache cleaner) to finish.
 3. Closes the storage (DB connection).
 4. Closes the log file if one was used.

This ensures a clean and deterministic application exit.
*/
func (a *App) Wait() {
	<-a.ctx.Done()
	a.wg.Wait()
	a.storage.Close()
	if a.logFile != nil {
		_ = a.logFile.Close()
	}
}
