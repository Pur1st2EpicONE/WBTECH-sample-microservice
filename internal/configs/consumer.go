package configs

import (
	"time"

	"github.com/spf13/viper"
)

type Consumer struct {
	Brokers                  []string
	Topic                    string
	ClientID                 string
	GroupID                  string
	AutoAck                  bool
	SessionTimeoutMs         int
	MaxPollIntervalMs        int
	SaveOrderRetryDelay      time.Duration
	SaveOrderRetryMax        int
	CommitRetryDelay         time.Duration
	CommitRetryMax           int
	EventTypeErrorsMax       int
	EventTypeErrorRetryDelay time.Duration
	DLQ                      Producer
	Notifier                 Notifier
	Kafka                    *Kafka // interchangeable
}

type Kafka struct {
	EnableAutoCommit    bool
	AutoOffsetReset     string
	SessionTimeoutMs    int
	MaxPollIntervalMs   int
	Acks                string
	EnableIdempotence   bool
	RetryMax            int
	SaveOrderRetryDelay time.Duration
	SaveOrderRetryMax   int
	CommitRetryDelay    time.Duration
	CommitRetryMax      int
}

func consConfig() Consumer {
	return Consumer{
		Brokers:                  viper.GetStringSlice("kafka.consumer.brokers"),
		Topic:                    viper.GetString("kafka.consumer.topic"),
		ClientID:                 viper.GetString("kafka.consumer.client_id"),
		GroupID:                  viper.GetString("kafka.consumer.group_id"),
		AutoAck:                  viper.GetBool("kafka.consumer.auto_ack"),
		SessionTimeoutMs:         viper.GetInt("kafka.consumer.session_timeout_ms"),
		MaxPollIntervalMs:        viper.GetInt("kafka.consumer.max_poll_interval_ms"),
		SaveOrderRetryDelay:      viper.GetDuration("kafka.consumer.save_order_retry_delay"),
		SaveOrderRetryMax:        viper.GetInt("kafka.consumer.save_order_retry_max"),
		CommitRetryDelay:         viper.GetDuration("kafka.consumer.commit_retry_delay"),
		CommitRetryMax:           viper.GetInt("kafka.consumer.commit_retry_max"),
		EventTypeErrorsMax:       viper.GetInt("kafka.consumer.event_type_errors_max"),
		EventTypeErrorRetryDelay: viper.GetDuration("kafka.consumer.event_type_error_retry_delay"),
		DLQ:                      dlqConfig(),
		Notifier:                 notifierConfig(),
		Kafka:                    kafkaConfig(),
	}
}

func kafkaConfig() *Kafka {
	return &Kafka{
		EnableAutoCommit: viper.GetBool("kafka.consumer.enable_auto_commit"),
		AutoOffsetReset:  viper.GetString("kafka.consumer.auto_offset_reset"),
	}
}

func dlqConfig() Producer {
	return Producer{
		Brokers:  viper.GetStringSlice("kafka.dlq.brokers"),
		Topic:    viper.GetString("kafka.dlq.topic"),
		ClientID: viper.GetString("kafka.dlq.client_id"),
		Kafka:    kafkaDlqConfig(),
	}
}

func kafkaDlqConfig() *KafkaProducer {
	return &KafkaProducer{
		Acks:              viper.GetString("kafka.dlq.acks"),
		Retries:           viper.GetInt("kafka.dlq.retry_max"),
		LingerMs:          viper.GetInt("kafka.dlq.linger_ms"),
		BatchSize:         viper.GetInt("kafka.dlq.batch_size"),
		CompressionType:   viper.GetString("kafka.dlq.compression_type"),
		EnableIdempotence: viper.GetBool("kafka.dlq.enable_idempotence"),
	}
}
