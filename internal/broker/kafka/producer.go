package kafka

import (
	"fmt"
	"strings"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	flushTimeOutMs = 5000
)

type KafkaProducer struct {
	producer *kafka.Producer
	logger   logger.Logger
}

func NewProducer(config configs.Producer, logger logger.Logger) (*KafkaProducer, error) {
	kafkaProducer, err := kafka.NewProducer(toMap(config))
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: kafkaProducer, logger: logger}, nil
}

func (p *KafkaProducer) Produce(message configs.Message) error {
	order := NewKafkaMessage(message.Key, message.Value, message.Topic)
	eventChan := make(chan kafka.Event)
	var err error
	for range 5 {
		if err = p.producer.Produce(order, eventChan); err != nil {
			time.Sleep(5 * time.Second)
			p.logger.LogError("failed to send order", err, "orderUID", strings.Trim(string(message.Key), `"`), "layer", "broker.kafka")
			continue
		} else {
			if message.DLQ {
				p.logger.LogInfo(fmt.Sprintf("worker %d — order is sent to DLQ", message.WorkerID), "orderUID", strings.Trim(string(message.Key), `"`), "workerID", fmt.Sprintf("%d", message.WorkerID), "layer", "broker.kafka")
			}
			return nil
		}
	}
	if message.DLQ {
		p.logger.LogError(fmt.Sprintf("worker %d — failed to send order to DLQ after 5 attempts", message.WorkerID), err, "orderUID", strings.Trim(string(message.Key), `"`), "workerID", fmt.Sprintf("%d", message.WorkerID), "layer", "broker.kafka")
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
	case <-time.After(5 * time.Second):
		return fmt.Errorf("event time out, make sure that Kafka is running and the address is correct")
	}
}

func NewKafkaMessage(key []byte, value []byte, topic string) *kafka.Message {
	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}
}

func (p *KafkaProducer) Close() {
	p.producer.Flush(flushTimeOutMs)
	p.producer.Close()
}
