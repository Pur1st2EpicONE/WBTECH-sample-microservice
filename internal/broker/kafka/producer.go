package kafka

import (
	"fmt"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// KafkaProducer wraps a Confluent Kafka producer and provides
// retry and logging mechanisms for sending messages.
//
// It handles sending messages to Kafka topics, including
// dead-letter queue (DLQ) messages, and ensures that failed
// deliveries are retried with configured delays.
type KafkaProducer struct {
	producer          *kafka.Producer // underlying Kafka producer
	logger            logger.Logger   // logger for producer events and errors
	flushTimeOut      int             // timeout for flushing pending messages on Close
	RetryAttempts     int             // number of attempts to retry sending a message
	produceRetryDelay time.Duration   // delay between retries when sending fails
	eventTimeout      time.Duration   // maximum wait time for delivery events
}

// NewProducer creates a new KafkaProducer instance based on the provided configuration
// and logger. It initializes the underlying Confluent Kafka producer and sets
// retry and timeout parameters.
func NewProducer(config configs.Producer, logger logger.Logger) (*KafkaProducer, error) {
	kafkaProducer, err := kafka.NewProducer(toMap(config))
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{
		producer:          kafkaProducer,
		logger:            logger,
		flushTimeOut:      config.FlushTimeOut,
		RetryAttempts:     config.RetryAttempts,
		produceRetryDelay: config.ProduceRetryDelay,
		eventTimeout:      config.EventTimeout,
	}, nil
}

// Produce sends a message to a Kafka topic.
//
// Steps:
//  1. Constructs a Kafka message from key, value, and topic.
//  2. Attempts to produce the message, retrying up to `RetryAttempts` times on failure.
//  3. Logs errors for each failed attempt.
//  4. If the message is for the DLQ, logs additional info on success or failure.
//  5. Waits for the delivery event or times out based on `eventTimeout`.
//  6. Returns any delivery or event errors to the caller.
//
// This ensures reliable message delivery with configurable retries, logging, and
// proper handling for dead-letter queue messages. The retry mechanism prevents
// transient Kafka issues from immediately failing message processing.
func (p *KafkaProducer) Produce(message configs.Message) error {
	order := NewKafkaMessage(message.Key, message.Value, message.Topic)
	eventChan := make(chan kafka.Event)
	var err error
	for range p.RetryAttempts {
		if err = p.producer.Produce(order, eventChan); err != nil {
			time.Sleep(p.produceRetryDelay)
			p.logger.LogError("producer — failed to send order", err, "orderUID", ToStr(message.Key), "layer", "broker.kafka")
			continue
		} else {
			if message.DLQ {
				p.logger.LogInfo(fmt.Sprintf("worker %d — order is sent to DLQ", message.WorkerID), "orderUID", ToStr(message.Key), "workerID", fmt.Sprintf("%d", message.WorkerID), "layer", "broker.kafka")
			}
			return nil
		}
	}
	if message.DLQ {
		p.logger.LogError(fmt.Sprintf("worker %d — failed to send order to DLQ after %d attempts", message.WorkerID, p.RetryAttempts), err, "orderUID", ToStr(message.Key), "workerID", fmt.Sprintf("%d", message.WorkerID), "layer", "broker.kafka")
	}
	select {
	case event := <-eventChan:
		switch eventType := event.(type) {
		case *kafka.Message:
			if eventType.TopicPartition.Error != nil {
				return eventType.TopicPartition.Error
			}
			return nil
		case *kafka.Error:
			return eventType
		default:
			return fmt.Errorf("unknown type of event: %T", event)
		}
	case <-time.After(p.eventTimeout):
		return fmt.Errorf("event time out, make sure that Kafka is running and the address is correct")
	}
}

// NewKafkaMessage constructs a Kafka message with the given key, value, and topic.
//
// It sets the partition to kafka.PartitionAny to let Kafka decide which partition
// the message should go to.
func NewKafkaMessage(key []byte, value []byte, topic string) *kafka.Message {
	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}
}

// Close flushes all pending messages and closes the producer.
//
// Ensures that no messages remain unsent and releases resources.
func (p *KafkaProducer) Close() {
	p.producer.Flush(p.flushTimeOut)
	p.producer.Close()
}
