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
	dlq      *KafkaProducer
	dlqTopic string
}

func NewConsumer(config configs.Consumer, logger logger.Logger) (*KafkaConsumer, error) {
	kafkaConsumer, err := kafka.NewConsumer(toMap(config))
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer -> %w", err)
	}
	if err := kafkaConsumer.Subscribe(config.Topic, nil); err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic -> %w", err)
	}
	dlq, err := NewProducer(config.DLQ, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create DLQ -> %w", err)
	}
	handler := newHandler()
	return &KafkaConsumer{
		consumer: kafkaConsumer,
		handler:  handler,
		dlq:      dlq,
		dlqTopic: config.DLQ.Topic}, nil
}

func toMap(config any) *kafka.ConfigMap {
	switch c := config.(type) {
	case configs.Consumer:
		return &kafka.ConfigMap{
			"bootstrap.servers":  strings.Join(c.Brokers, ","),
			"group.id":           c.GroupID,
			"auto.offset.reset":  c.Kafka.AutoOffsetReset,
			"enable.auto.commit": c.Kafka.EnableAutoCommit,
			"client.id":          c.ClientID,
		}
	case configs.Producer:
		var acksValue int
		switch c.Kafka.Acks {
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
			"bootstrap.servers":     strings.Join(c.Brokers, ","),
			"request.required.acks": acksValue,
			"retries":               c.Kafka.Retries,
			"linger.ms":             c.Kafka.LingerMs,
			"batch.size":            c.Kafka.BatchSize,
			"compression.codec":     c.Kafka.CompressionType,
			"enable.idempotence":    c.Kafka.EnableIdempotence,
			"client.id":             c.ClientID,
		}
	default:
		return nil
	}
}

func (c *KafkaConsumer) Run(ctx context.Context, storage repository.Storage, logger logger.Logger) {
	logger.LogInfo("receiving orders", "layer", "consumer")
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
					if err := c.handler.SaveOrder(eventType.Value, storage, logger); err != nil {
						lastErr = err
						retryCnt++
						if retryCnt < maxRetries {
							time.Sleep(5 * time.Second)
							continue
						}
						break
					}
					if err := c.commitWithRetry(eventType, logger); err != nil {
						panic(err.Error())
					}
					break
				}
				if retryCnt >= maxRetries {
					logger.LogError(fmt.Sprintf("failed to process order after %d retries -> %v", maxRetries, lastErr), lastErr, "orderUID", string(eventType.Key), "layer", "consumer")
					headers := make(map[string]string, len(eventType.Headers))
					for _, h := range eventType.Headers {
						headers[h.Key] = string(h.Value)
					}
					msg := configs.Message{
						Topic:     c.dlqTopic,
						Key:       eventType.Key,
						Value:     eventType.Value,
						Headers:   headers,
						Timestamp: eventType.Timestamp,
						Metadata:  map[string]any{"retryCount": retryCnt},
						DLQ:       true,
					}
					c.dlq.Produce(msg)
					if err := c.commitWithRetry(eventType, logger); err != nil {
						panic("kafka — critical error")
					}
				}
			case kafka.Error:
				logger.LogError("event type error -> %v", eventType, "layer", "consumer")
			}
		}
	}
}

func (c *KafkaConsumer) commitWithRetry(msg *kafka.Message, logger logger.Logger) error {
	for range 3 {
		logger.LogInfo("attempting commit",
			"topic", *msg.TopicPartition.Topic,
			"partition", msg.TopicPartition.Partition,
			"offset", msg.TopicPartition.Offset,
		)
		if _, err := c.consumer.CommitMessage(msg); err != nil {
			time.Sleep(1 * time.Second)
		} else {
			return nil
		}
	}
	return fmt.Errorf("failed to commit offset after 3 attempts")
}

func (c *KafkaConsumer) Close(logger logger.Logger) {
	if err := c.consumer.Close(); err != nil {
		logger.LogError("failed to stop properly -> %v", err, "layer", "consumer")
	}
	logger.LogInfo("stopped receiving orders", "layer", "consumer.kafka")
}
