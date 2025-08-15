package kafka

import (
	"fmt"
	"strings"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const maxRetries = 3

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

func (c *Consumer) Run(storage *repository.Storage) {
	logger.LogInfo("consumer — receiving orders")
	for {
		kafkaMsg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logger.LogFatal("consumer — failed to read message: %v", err)
		}
		var lastErr error
		retryCnt := 0
		for retryCnt < maxRetries {
			if err := c.handler.SaveOrder(kafkaMsg.Value, *storage); err != nil {
				lastErr = err
				retryCnt++
				if retryCnt < maxRetries {
					time.Sleep(5 * time.Second)
					continue
				}
				break
			}
			if _, err := c.consumer.CommitMessage(kafkaMsg); err != nil {
				logger.LogError("consumer — failed to commit offset:", err)
			}
			break
		}
		if retryCnt >= maxRetries {
			logger.LogError(fmt.Sprintf("consumer — failed to get message after %d retries: %v", maxRetries, lastErr), lastErr)
			continue
		}
	}
}

func (c *Consumer) Close() {
	logger.LogInfo("consumer — stopping")
	err := c.consumer.Close()
	if err != nil {
		logger.LogFatal("consumer — failed to stop properly: %v", err)
	}
}
