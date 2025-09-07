package configs

import (
	"time"

	"github.com/spf13/viper"
)

// Producer holds configuration for a message producer.
//
// Includes connection info, topic, retry policy, batching, and optional Kafka-specific settings.
type Producer struct {
	Brokers           []string       // list of Kafka brokers for the producer
	Topic             string         // topic to produce messages to
	ClientID          string         // producer client ID
	MsgsToSend        int            // number of messages (orders) to send (order-producer-specific)
	FlushTimeOut      int            // maximum time to wait for message flush
	RetryAttempts     int            // number of application-level retry attempts
	ProduceRetryDelay time.Duration  // delay between application-level retry attempts
	EventTimeout      time.Duration  // overall timeout for event processing
	Kafka             *KafkaProducer // Kafka-specific producer parameters
}

// KafkaProducer stores Kafka-specific producer parameters.
//
// Configures acknowledgements, idempotence, retries, batching, and compression.
type KafkaProducer struct {
	Acks              string // wait for all in-sync replicas to acknowledge
	EnableIdempotence bool   // number of replicas that must acknowledge writes
	Retries           int    // maximum number of automatic retry attempts by Kafka library
	LingerMs          int    // time to wait before sending a batch
	BatchSize         int    // maximum batch size in bytes
	CompressionType   string // compression algorithm for messages
}

// Message represents a Kafka message payload.
//
// Contains topic, key/value data, headers, timestamps, metadata, DLQ flag, and originating worker ID.
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
		RetryAttempts:     viper.GetInt("kafka.producer.produce_retry_attempts"),
		ProduceRetryDelay: viper.GetDuration("kafka.producer.produce_retry_delay"),
		EventTimeout:      viper.GetDuration("kafka.producer.event_timeout"),
		Kafka:             kafkaProdConfig(),
	}, nil
}

func kafkaProdConfig() *KafkaProducer {
	return &KafkaProducer{
		Acks:              viper.GetString("kafka.producer.acks"),
		EnableIdempotence: viper.GetBool("kafka.producer.enable_idempotence"),
		Retries:           viper.GetInt("kafka.producer.retry_max"),
		LingerMs:          viper.GetInt("kafka.producer.linger_ms"),
		BatchSize:         viper.GetInt("kafka.producer.batch_size"),
		CompressionType:   viper.GetString("kafka.producer.compression_type"),
	}
}
