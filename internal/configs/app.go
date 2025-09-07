/*
Package configs provides structured configuration for the application.

It loads settings from environment variables and configuration files (.env and config.yaml)
using godotenv and viper. The package centralizes configuration for:

  - HTTP server parameters
  - Database connections
  - In-memory caching
  - Kafka consumer/producer settings
  - Logging
  - Notifications (e.g., Telegram bot)
  - Worker behavior and shutdown policies

Sensitive values like passwords and API tokens are read from environment variables
to avoid hardcoding secrets. Explicit viper.Get* calls are used to preserve polymorphism,
allowing optional components (Kafka, NATS, RabbitMQ) without tying them directly to the App struct.
*/
package configs

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// App holds all top-level configuration for the application,
// including server, database, cache, consumer, logger, notifier, and worker settings.
type App struct {
	Server          Server
	Database        Database
	Consumer        Consumer
	Cache           Cache
	Logger          Logger
	Notifier        Notifier
	Workers         int
	RestartOnPanic  bool
	RestartDelay    time.Duration
	DbCheckInterval time.Duration
	DbMaxChecks     int
}

// Server contains HTTP server configuration.
type Server struct {
	Port            string        // port to listen on
	ReadTimeout     time.Duration // maximum duration for reading the request
	WriteTimeout    time.Duration // maximum duration before timing out writes
	MaxHeaderBytes  int           // maximum size of request headers
	ShutdownTimeout time.Duration // graceful shutdown timeout
}

// Database holds database connection parameters.
type Database struct {
	Driver          string
	Host            string
	Port            string
	Username        string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// Cache contains in-memory caching configuration.
type Cache struct {
	SaveInCache     bool          // whether to store orders in memory
	CacheSize       int           // maximum number of orders to cache
	BgCleanup       bool          // whether background cleaner is enabled
	CleanupInterval time.Duration // period between cleanup cycles
	OrderTTL        time.Duration // time-to-live for cached orders
	PauseDuration   time.Duration // pause duration when DB is unreachable
}

// Logger defines logging configuration.
type Logger struct {
	LogDir string // directory to store logs
	Debug  bool   // whether debug mode is enabled
}

// Notifier holds configuration for external notifications.
type Notifier struct {
	Token    string // authentication token (e.g., Telegram bot)
	Receiver string // recipient identifier (e.g., Telegram chat ID)
}

// Load reads environment variables and configuration files to
// populate the App struct with all settings.
//
// It uses godotenv for .env files and viper for config.yaml.
// Sensitive fields like database password and bot token are read from the environment.
//
// Explicit viper.Get* calls are used instead of viper.Unmarshal
// to preserve polymorphism in interchangeable sub-structures like Kafka, NATS and RabbitMQ.
func Load() (App, error) {
	if err := godotenv.Load(); err != nil {
		return App{}, fmt.Errorf("godotenv — failed to %v", err) // phrasing is odd here, but gives a clean error message in logs

	}

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return App{}, fmt.Errorf("viper — %v", err)
	}

	return App{
		Server:          srvConfig(),
		Database:        dbConfig(),
		Cache:           cacheConfig(),
		Consumer:        consConfig(),
		Logger:          loggerConfig(),
		Notifier:        notifierConfig(),
		Workers:         viper.GetInt("app.workers.active_consumer_workers"),
		RestartOnPanic:  viper.GetBool("app.workers.restart_on_panic"),
		RestartDelay:    viper.GetDuration("app.workers.restart_delay"),
		DbCheckInterval: viper.GetDuration("app.db.connection_check_interval"),
		DbMaxChecks:     viper.GetInt("app.db.max_rtrs_bfr_cache_only_mode"),
	}, nil
}

// srvConfig reads server-related configuration from viper.
func srvConfig() Server {
	return Server{
		Port:            viper.GetString("server.port"),
		ReadTimeout:     viper.GetDuration("server.read_timeout"),
		WriteTimeout:    viper.GetDuration("server.write_timeout"),
		MaxHeaderBytes:  viper.GetInt("server.max_header_bytes"),
		ShutdownTimeout: viper.GetDuration("server.shutdown_timeout"),
	}
}

// dbConfig reads database-related configuration from viper and environment variables.
func dbConfig() Database {
	return Database{
		Driver:          viper.GetString("database.driver"),
		Host:            viper.GetString("database.host"),
		Port:            viper.GetString("database.port"),
		Username:        viper.GetString("database.username"),
		Password:        os.Getenv("DB_PASSWORD"),
		DBName:          viper.GetString("database.dbname"),
		SSLMode:         viper.GetString("database.sslmode"),
		MaxOpenConns:    viper.GetInt("database.max_open_conns"),
		MaxIdleConns:    viper.GetInt("database.max_idle_conns"),
		ConnMaxLifetime: viper.GetDuration("database.conn_max_lifetime"),
		ConnMaxIdleTime: viper.GetDuration("database.conn_max_idle_time"),
	}
}

// cacheConfig reads cache settings from viper.
func cacheConfig() Cache {
	return Cache{
		SaveInCache:     viper.GetBool("cache.save_in_cache"),
		CacheSize:       viper.GetInt("cache.cache_size"),
		BgCleanup:       viper.GetBool("cache.background_cleanup"),
		CleanupInterval: viper.GetDuration("cache.cleanup_interval"),
		OrderTTL:        viper.GetDuration("cache.order_ttl"),
		PauseDuration:   viper.GetDuration("cache.clnr_pause_on_db_conn_check"),
	}
}

// loggerConfig reads logger settings from viper.
func loggerConfig() Logger {
	return Logger{
		LogDir: viper.GetString("app.logger.log_directory"),
		Debug:  viper.GetBool("app.logger.debug_mode"),
	}
}

// notifierConfig reads notifier settings from viper and environment variables.
func notifierConfig() Notifier {
	return Notifier{
		Token:    os.Getenv("TG_BOT_TOKEN"),
		Receiver: viper.GetString("notifier.telegram.chat_id"),
	}
}
