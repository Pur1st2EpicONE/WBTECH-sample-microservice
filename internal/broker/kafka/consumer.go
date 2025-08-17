package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const maxRetries = 3

type KafkaConsumer struct {
	consumer *kafka.Consumer
	handler  *Handler
}

func NewConsumer(config configs.Consumer) (*KafkaConsumer, error) {
	kafkaConsumer, err := kafka.NewConsumer(toMap(config))
	if err != nil {
		return nil, err
	}
	if err := kafkaConsumer.Subscribe(config.Topic, nil); err != nil {
		return nil, err
	}
	orderHandler := NewHandler()
	return &KafkaConsumer{consumer: kafkaConsumer, handler: orderHandler}, nil
}

func toMap(config configs.Consumer) *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(config.Brokers, ","),
		"group.id":           config.GroupID,
		"auto.offset.reset":  config.Kafka.AutoOffsetReset,
		"enable.auto.commit": config.Kafka.EnableAutoCommit,
		"client.id":          config.ClientID,
	}
}

func (c *KafkaConsumer) Run(ctx context.Context, storage *repository.Storage) {
	logger.LogInfo("consumer — receiving orders")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			event := c.consumer.Poll(100)
			if event == nil {
				continue
			}
			switch eventType := event.(type) {
			case *kafka.Message:
				var lastErr error
				retryCnt := 0
				for retryCnt < maxRetries {
					if err := c.handler.SaveOrder(eventType.Value, *storage); err != nil {
						lastErr = err
						retryCnt++
						if retryCnt < maxRetries {
							time.Sleep(5 * time.Second)
							continue
						}
						break
					}
					if _, err := c.consumer.CommitMessage(eventType); err != nil {
						logger.LogError("consumer — failed to commit offset:", err)
					}
					break
				}
				if retryCnt >= maxRetries {
					logger.LogError(fmt.Sprintf("consumer — failed to process message after %d retries: %v", maxRetries, lastErr), lastErr)
				}
			case kafka.Error:
				logger.LogError("consumer — Kafka error:", eventType)
			}
		}
	}
}

func (c *KafkaConsumer) Close() {
	if err := c.consumer.Close(); err != nil {
		logger.LogError("consumer — failed to stop properly: %v", err)
	}
	logger.LogInfo("consumer — stopped")
}
