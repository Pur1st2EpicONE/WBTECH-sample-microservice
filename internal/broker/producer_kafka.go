package broker

import (
	"fmt"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	flushTimeOutMs = 5000
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(clusterHosts []string) (*Producer, error) {
	kafkaProducer, err := newKafkaProducer(clusterHosts)
	if err != nil {
		return nil, err
	}
	return &Producer{producer: kafkaProducer}, nil
}

func newKafkaProducer(clusterHosts []string) (*kafka.Producer, error) {
	config := newProducerConfig(clusterHosts)
	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

func newProducerConfig(clusterHosts []string) *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(clusterHosts, ","),
		"acks":               "all",
		"retries":            3,
		"linger.ms":          5,
		"batch.size":         65536,
		"compression.type":   "snappy",
		"enable.idempotence": false,
		"client.id":          "order-producer",
	}
}

func (p *Producer) Produce(data string, topic string) error {
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

func newKafkaMessage(data string, topic string) (*kafka.Message, chan kafka.Event) {
	eventChan := make(chan kafka.Event)
	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            nil,
		Value:          []byte(data),
	}, eventChan
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeOutMs)
	p.producer.Close()
}
