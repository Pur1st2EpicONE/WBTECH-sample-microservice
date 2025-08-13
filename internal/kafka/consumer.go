package kafka

import (
	"log"
	"strings"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/repository"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer *kafka.Consumer
	handler  *OrderHandler
}

func NewConsumer(clusterHosts []string, consGroupID string, topic string) (*Consumer, error) {
	kafkaConsumer, err := newKafkaConsumer(clusterHosts, consGroupID)
	if err != nil {
		return nil, err
	}
	if err := kafkaConsumer.Subscribe(topic, nil); err != nil {
		return nil, err
	}
	orderHandler := NewOrderHandler()
	return &Consumer{consumer: kafkaConsumer, handler: orderHandler}, nil
}

func newKafkaConsumer(clusterHosts []string, consGroupID string) (*kafka.Consumer, error) {
	config := newConsumerConfig(clusterHosts, consGroupID)
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func newConsumerConfig(clusterHosts []string, consGroupID string) *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(clusterHosts, ","),
		"group.id":           consGroupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
		"client.id":          "order-consumer",
	}
}

func (c *Consumer) Run(db repository.Storage) error {
	for {
		kafkaMsg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			return err
		}
		if err := c.handler.SaveOrder(kafkaMsg.Value, db); err != nil {
			log.Printf("Failed to save order: %v", err)
			continue
		}
		if _, err := c.consumer.CommitMessage(kafkaMsg); err != nil {
			log.Printf("Failed to commit offset: %v", err)
		}

	}
}
