package configs

import (
	"github.com/spf13/viper"
)

type Producer struct {
	Brokers  []string
	ClientID string
	Kafka    *KafkaProducer
	NATS     *NATSProducer
	RabbitMQ *RabbitMQProducer
}

type KafkaProducer struct {
	Acks              string
	EnableIdempotence bool
	Retries           int
	LingerMs          int
	BatchSize         int
	CompressionType   string
}

type NATSProducer struct {
	Subject   string
	Queue     string
	Durable   string
	AckWaitMs int
}

type RabbitMQProducer struct {
	Exchange   string
	RoutingKey string
	Mandatory  bool
	Immediate  bool
}

func ProdConfig() (Producer, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return Producer{}, err
	}

	return Producer{
		Brokers:  viper.GetStringSlice("kafka.producer.brokers"),
		ClientID: viper.GetString("kafka.producer.client_id"),
		Kafka:    kafkaProdConfig(),
		NATS:     nil,
		RabbitMQ: nil,
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
