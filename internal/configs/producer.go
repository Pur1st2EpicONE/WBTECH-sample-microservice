package configs

import (
	"time"

	"github.com/spf13/viper"
)

type Producer struct {
	Brokers           []string
	Topic             string
	ClientID          string
	MsgsToSend        int
	FlushTimeOut      int
	ProduceRetryDelay time.Duration
	EventTimeout      time.Duration
	Kafka             *KafkaProducer
}

type KafkaProducer struct {
	Acks              string
	EnableIdempotence bool
	Retries           int
	LingerMs          int
	BatchSize         int
	CompressionType   string
}

type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Headers   map[string]string
	Timestamp time.Time
	Metadata  map[string]any
	DLQ       bool
	WorkerID  int
}

func ProdConfig() (Producer, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return Producer{}, err
	}

	return Producer{
		Brokers:           viper.GetStringSlice("kafka.producer.brokers"),
		Topic:             viper.GetString("kafka.producer.topic"),
		ClientID:          viper.GetString("kafka.producer.client_id"),
		MsgsToSend:        viper.GetInt("kafka.producer.messages_to_send"),
		FlushTimeOut:      viper.GetInt("kafka.producer.flush_time_out_ms"),
		ProduceRetryDelay: viper.GetDuration("kafka.producer.produce_retry_delay"),
		EventTimeout:      viper.GetDuration("kafka.producer.event_timeout"),
		Kafka:             kafkaProdConfig(),
	}, nil
}

func kafkaProdConfig() *KafkaProducer {
	return &KafkaProducer{
		Acks:              viper.GetString("kafka.producer.acks"),
		Retries:           viper.GetInt("kafka.producer.retry_max"),
		LingerMs:          viper.GetInt("kafka.producer.linger_ms"),
		BatchSize:         viper.GetInt("kafka.producer.batch_size"),
		CompressionType:   viper.GetString("kafka.producer.compression_type"),
		EnableIdempotence: viper.GetBool("kafka.producer.enable_idempotence"),
	}
}
