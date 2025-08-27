package configs

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type App struct {
	Server         Server
	Database       Database
	Consumer       Consumer
	Cache          Cache
	Logger         Logger
	Notifier       Notifier
	Workers        int
	RestartOnPanic bool
	RestartDelay   time.Duration
}

type Server struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

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

type Cache struct {
	SaveInCache   bool
	CacheSize     int
	BgCleanup     bool
	CleanupPeriod time.Duration
	OrderTTL      time.Duration
}

type Logger struct {
	LogDir string
	Debug  bool
}

type Notifier struct {
	Token    string
	Receiver string
}

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
		Server:         srvConfig(),
		Database:       dbConfig(),
		Cache:          cacheConfig(),
		Consumer:       consConfig(),
		Logger:         loggerConfig(),
		Notifier:       notifierConfig(),
		Workers:        viper.GetInt("app.workers.active_consumer_workers"),
		RestartOnPanic: viper.GetBool("app.workers.restart_on_panic"),
		RestartDelay:   viper.GetDuration("app.workers.restart_delay"),
	}, nil
}

/*
Using explicit viper.Get* calls instead of viper.Unmarshal because
unmarshaling would force embedding KafkaConsumer fields directly into App,
breaking polymorphism. This way, we keep Kafka, NATS, RabbitMQ as optional
interchangeable sub-structures without tying them to App directly.
*/

func srvConfig() Server {
	return Server{
		Port:           viper.GetString("server.port"),
		ReadTimeout:    viper.GetDuration("server.read_timeout"),
		WriteTimeout:   viper.GetDuration("server.write_timeout"),
		MaxHeaderBytes: viper.GetInt("server.max_header_bytes"),
	}
}

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

func cacheConfig() Cache {
	return Cache{
		SaveInCache:   viper.GetBool("cache.save_in_cache"),
		CacheSize:     viper.GetInt("cache.cache_size"),
		BgCleanup:     viper.GetBool("cache.background_cleanup"),
		CleanupPeriod: viper.GetDuration("cache.cleanup_period"),
		OrderTTL:      viper.GetDuration("cache.order_ttl"),
	}
}

func loggerConfig() Logger {
	return Logger{
		LogDir: viper.GetString("app.logger.log_directory"),
		Debug:  viper.GetBool("app.logger.debug_mode"),
	}
}

func notifierConfig() Notifier {
	return Notifier{
		Token:    os.Getenv("TG_BOT_TOKEN"),
		Receiver: viper.GetString("notifier.telegram.chat_id"),
	}
}
