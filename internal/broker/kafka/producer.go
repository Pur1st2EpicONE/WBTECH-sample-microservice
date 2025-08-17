package kafka

import (
	"fmt"
	"strings"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	flushTimeOutMs = 5000
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(config configs.Producer) (*Producer, error) {
	kafkaProducer, err := kafka.NewProducer(confToMap(config))
	if err != nil {
		return nil, err
	}
	return &Producer{producer: kafkaProducer}, nil
}

func confToMap(config configs.Producer) *kafka.ConfigMap {
	acks := config.Kafka.Acks
	var acksValue int
	switch strings.ToLower(acks) {
	case "all", "-1":
		acksValue = -1
	case "0":
		acksValue = 0
	case "1":
		acksValue = 1
	default:
		acksValue = -1
	}

	return &kafka.ConfigMap{
		"bootstrap.servers":     strings.Join(config.Brokers, ","),
		"request.required.acks": acksValue,
		"retries":               config.Kafka.Retries,
		"linger.ms":             config.Kafka.LingerMs,
		"batch.size":            config.Kafka.BatchSize,
		"compression.codec":     config.Kafka.CompressionType,
		"enable.idempotence":    config.Kafka.EnableIdempotence,
		"client.id":             config.ClientID,
	}
}

func (p *Producer) Produce(data []byte, topic string) error {
	kafkaMessage, eventChan := newKafkaMessage(data, topic)
	if err := p.producer.Produce(kafkaMessage, eventChan); err != nil {
		return err
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
		fmt.Println("event time out, make sure that Kafka is running and the address is correct")
		return nil
	}
}

func newKafkaMessage(data []byte, topic string) (*kafka.Message, chan kafka.Event) {
	eventChan := make(chan kafka.Event)
	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            nil,
		Value:          data,
	}, eventChan
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeOutMs)
	p.producer.Close()
}
