package configs

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type App struct {
	Server   Server
	Database Database
	Consumer Consumer
	Cache    Cache
}

type Server struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

type Database struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type Cache struct {
	SaveInCache   bool
	CacheSize     int
	BgCleanup     bool
	CleanupPeriod time.Duration
	OrderTTL      time.Duration
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
		Server:   srvConfig(),
		Database: dbConfig(),
		Cache:    cacheConfig(),
		Consumer: consConfig(),
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
		Driver:   viper.GetString("database.driver"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Username: viper.GetString("database.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("database.dbname"),
		SSLMode:  viper.GetString("database.sslmode"),
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
