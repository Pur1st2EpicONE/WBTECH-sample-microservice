package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/notifier"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConsumer struct {
	consumer                 *kafka.Consumer
	handler                  *Handler
	dlq                      *KafkaProducer
	dlqTopic                 string
	saveOrderRetryDelay      time.Duration
	saveOrderRetryMax        int
	commitRetryDelay         time.Duration
	commitRetryMax           int
	eventTypeErrorsMax       int
	eventTypeErrorRetryDelay time.Duration
	notifier                 notifier.Notifier
}

func NewConsumer(config configs.Consumer, logger logger.Logger) (*KafkaConsumer, error) {
	kafkaConsumer, err := kafka.NewConsumer(toMap(config))
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	if err := kafkaConsumer.Subscribe(config.Topic, nil); err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}
	dlq, err := NewProducer(config.DLQ, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create DLQ: %w", err)
	}
	handler := newHandler()
	return &KafkaConsumer{
		consumer:                 kafkaConsumer,
		handler:                  handler,
		dlq:                      dlq,
		dlqTopic:                 config.DLQ.Topic,
		saveOrderRetryDelay:      config.SaveOrderRetryDelay,
		saveOrderRetryMax:        config.SaveOrderRetryMax,
		commitRetryDelay:         config.CommitRetryDelay,
		commitRetryMax:           config.CommitRetryMax,
		eventTypeErrorsMax:       config.EventTypeErrorsMax,
		eventTypeErrorRetryDelay: config.EventTypeErrorRetryDelay,
		notifier:                 notifier.NewNotifier(config.Notifier)}, nil
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

func (c *KafkaConsumer) Run(ctx context.Context, storage repository.Storage, logger logger.Logger, workerID int) {
	logger.LogInfo(fmt.Sprintf("worker %d — receiving orders", workerID), "layer", "broker.kafka")
	eventTypeErrors := 0
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
				eventTypeErrors = 0
				var lastErr error
				retryCnt := 0
				for retryCnt < c.saveOrderRetryMax {
					if err := c.handler.SaveOrder(eventType.Value, storage, logger, workerID); err != nil {
						lastErr = err
						retryCnt++
						if retryCnt < c.saveOrderRetryMax {
							time.Sleep(c.saveOrderRetryDelay)
							continue
						}
						break
					}
					if err := c.commitWithRetry(eventType); err != nil {
						logger.LogError(fmt.Sprintf("worker %d — critical error", workerID), err, "orderUID", ToStr(eventType.Key), "workerID", fmt.Sprintf("%d", workerID), "layer", "broker.kafka")
						c.notifier.Notify(fmt.Sprintf("CRITICAL ERROR — Kafka commit failed\nworkerID=%d\norderUID=%s", workerID, ToStr(eventType.Key)))
						panic(fmt.Sprintf("worker self-termination: offset commit failed (workerID=%d, orderUID=%s)", workerID, ToStr(eventType.Key)))
					}
					break
				}
				if retryCnt >= c.saveOrderRetryMax {
					logger.LogError(fmt.Sprintf("worker %d — failed to process order after %d retries", workerID, c.saveOrderRetryMax), lastErr, "orderUID", ToStr(eventType.Key), "workerID", fmt.Sprintf("%d", workerID), "layer", "broker.kafka")
					c.sendToDLQ(eventType, retryCnt, workerID)
				}
			case kafka.Error:
				eventTypeErrors++
				logger.LogError("consumer — event type error", eventType, "layer", "broker.kafka")
				if eventTypeErrors > c.eventTypeErrorsMax {
					c.notifier.Notify(fmt.Sprintf("CRITICAL ERROR — Kafka broker is unreachable\nworkerID=%d", workerID))
					panic(fmt.Sprintf("worker self-termination: kafka is down (workerID=%d)", workerID))
				}
				time.Sleep(c.eventTypeErrorRetryDelay)
			}
		}
	}
}

func (c *KafkaConsumer) commitWithRetry(msg *kafka.Message) error {
	var err error
	for range c.commitRetryMax {
		if _, err = c.consumer.CommitMessage(msg); err != nil {
			time.Sleep(c.commitRetryDelay)
		} else {
			return nil
		}
	}
	return fmt.Errorf("critical error — failed to commit offset after %d attempts: %w", c.commitRetryMax, err)
}

func (c *KafkaConsumer) sendToDLQ(eventType *kafka.Message, retryCnt int, workerID int) {
	headers := make(map[string]string)
	msg := configs.Message{
		Topic:     c.dlqTopic,
		Key:       eventType.Key,
		Value:     eventType.Value,
		Headers:   headers,
		Timestamp: eventType.Timestamp,
		Metadata:  map[string]any{"retryCount": retryCnt},
		DLQ:       true,
		WorkerID:  workerID,
	}
	if err := c.dlq.Produce(msg); err != nil {
		c.notifier.Notify(fmt.Sprintf("CRITICAL ERROR — failed to send order to DLQ\nworkerID=%d\norderUID=%s", workerID, ToStr(eventType.Key)))
		panic(fmt.Sprintf("worker self-termination: failed to send order to DLQ (workerID=%d, orderUID=%s)", workerID, ToStr(eventType.Key)))
	}
	if err := c.commitWithRetry(eventType); err != nil {
		c.notifier.Notify(fmt.Sprintf("CRITICAL ERROR — order sent to DLQ but offset commit failed\nworkerID=%d\norderUID=%s", workerID, ToStr(eventType.Key)))
		panic(fmt.Sprintf("worker self-termination: order sent to DLQ but offset commit failed (workerID=%d, orderUID=%s)", workerID, ToStr(eventType.Key)))
	}
}

func ToStr(key []byte) string {
	return strings.Trim(string(key), `"`)
}

func (c *KafkaConsumer) Close(logger logger.Logger) {
	if err := c.consumer.Close(); err != nil {
		logger.LogError("consumer — failed to stop properly", err, "layer", "broker.kafka")
	}
	logger.LogInfo("consumer — stopped receiving orders", "layer", "broker.kafka")
}
