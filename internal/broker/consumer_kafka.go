package broker

import (
	"context"
	"strings"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer *kafka.Consumer
	handler  *Handler
}

func NewConsumer(clusterHosts []string, consGroupID string, topic string) (*Consumer, error) {
	kafkaConsumer, err := newKafkaConsumer(clusterHosts, consGroupID)
	if err != nil {
		return nil, err
	}
	if err := kafkaConsumer.Subscribe(topic, nil); err != nil {
		return nil, err
	}
	orderHandler := NewHandler()
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

func (c *Consumer) Run(ctx context.Context, storage *repository.Storage) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		event := c.consumer.Poll(100)
		if event == nil {
			continue
		}
		switch eventType := event.(type) {
		case *kafka.Message:
			if err := c.processMessage(ctx, eventType, storage); err != nil {
				logger.LogError("consumer — failed to process message", err)
			}
		case kafka.Error:
			logger.LogError("consumer — Kafka error:", eventType)
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg *kafka.Message, storage *repository.Storage) error {
	const maxRetries = 5
	var lastErr error
	for try := range maxRetries {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		if err := c.handler.SaveOrder(msg.Value, *storage); err != nil {
			lastErr = err
			time.Sleep(time.Second << try)
			continue
		}
		if _, err := c.consumer.CommitMessage(msg); err != nil {
			logger.LogError("consumer — failed to commit offset:", err)
		}
		return nil
	}
	return lastErr
}

func (c *Consumer) Close() {
	if err := c.consumer.Close(); err != nil {
		logger.LogError("consumer — failed to stop properly: %v", err)
	}
	logger.LogInfo("consumer — stopped")
}
