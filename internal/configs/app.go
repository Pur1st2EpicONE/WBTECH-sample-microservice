package configs

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type App struct {
	Server   Server
	Database Database
	Consumer Consumer
}

type Server struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

type Database struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func Load() (App, error) {
	if err := godotenv.Load(); err != nil {
		return App{}, err
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return App{}, err
	}

	return App{
		Server:   srvConfig(),
		Database: dbConfig(),
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
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Username: viper.GetString("database.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("database.dbname"),
		SSLMode:  viper.GetString("database.sslmode"),
	}
}

// type Kafka struct {
// 	EnableAutoCommit  bool
// 	AutoOffsetReset   string
// 	SessionTimeoutMs  int
// 	MaxPollIntervalMs int

// 	Acks              string
// 	EnableIdempotence bool
// 	RetryMax          int
// 	Retries           int
// 	LingerMs          int
// 	BatchSize         int
// 	CompressionType   string
// }
