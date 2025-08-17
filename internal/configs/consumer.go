package configs

import "github.com/spf13/viper"

type Consumer struct {
	Brokers           []string
	Topic             string
	ClientID          string
	GroupID           string
	AutoAck           bool
	SessionTimeoutMs  int
	MaxPollIntervalMs int
	Kafka             *Kafka
	NATS              *NATS
	RabbitMQ          *RabbitMQ
}

type Kafka struct {
	EnableAutoCommit  bool
	AutoOffsetReset   string
	SessionTimeoutMs  int
	MaxPollIntervalMs int
	Acks              string
	EnableIdempotence bool
	RetryMax          int
}

type NATS struct {
	DurableName string
	AckWaitMs   int
	Subject     string
}

type RabbitMQ struct {
	QueueName     string
	ConsumerTag   string
	PrefetchCount int
	Exchange      string
	RoutingKey    string
	Persistent    bool
}

func consConfig() Consumer {
	return Consumer{
		Brokers:           viper.GetStringSlice("kafka.consumer.brokers"),
		Topic:             viper.GetString("kafka.consumer.topic"),
		ClientID:          viper.GetString("kafka.consumer.client_id"),
		GroupID:           viper.GetString("kafka.consumer.group_id"),
		AutoAck:           viper.GetBool("kafka.consumer.auto_ack"),
		SessionTimeoutMs:  viper.GetInt("kafka.consumer.session_timeout_ms"),
		MaxPollIntervalMs: viper.GetInt("kafka.consumer.max_poll_interval_ms"),
		Kafka:             kafkaConfig(),
		NATS:              nil,
		RabbitMQ:          nil,
	}
}

func kafkaConfig() *Kafka {
	return &Kafka{
		EnableAutoCommit:  viper.GetBool("kafka.consumer.enable_auto_commit"),
		AutoOffsetReset:   viper.GetString("kafka.consumer.auto_offset_reset"),
		SessionTimeoutMs:  viper.GetInt("kafka.consumer.session_timeout_ms"),
		MaxPollIntervalMs: viper.GetInt("kafka.consumer.max_poll_interval_ms"),
	}
}
