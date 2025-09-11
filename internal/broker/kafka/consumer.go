/*
Package kafka provides Kafka-based implementations of broker interfaces.

It includes:
  - KafkaConsumer: a consumer instance that processes messages from Kafka topics.
  - KafkaProducer: a producer instance used for sending messages (e.g., to a DLQ).

KafkaConsumer handles message consumption, retries, DLQ routing, and critical error
notifications. It is designed to be supervised by the orchestration layer
(App) and can trigger self-termination on unrecoverable errors.
*/
package kafka

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/notifier"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

/*
KafkaConsumer represents a single Kafka consumer instance.

It is responsible for:
  - Polling messages from a Kafka topic.
  - Processing messages and saving them to storage.
  - Handling retries for message processing and offset commits.
  - Sending failed messages to a dead-letter queue (DLQ).
  - Logging critical errors and notifying via a notifier.
  - Self-termination if unrecoverable errors occur.

KafkaConsumer is typically managed and monitored by the App orchestration layer.
*/
type KafkaConsumer struct {
	consumer                 *kafka.Consumer   // underlying Kafka consumer
	handler                  *Handler          // message handler for processing orders
	dlq                      *KafkaProducer    // producer for dead-letter queue
	dlqTopic                 string            // DLQ topic name
	saveOrderRetryDelay      time.Duration     // delay between retries when saving order fails
	saveOrderRetryMax        int               // maximum retries for saving an order
	commitRetryDelay         time.Duration     // delay between retries when committing offset
	commitRetryMax           int               // maximum retries for committing offset
	eventTypeErrorsMax       int               // max consecutive broker errors before panic
	eventTypeErrorRetryDelay time.Duration     // delay after broker error before retry
	dbConnectionCheckDelay   time.Duration     // delay between database connection checks when connection errors occur
	notifier                 notifier.Notifier // notifier for critical errors
}

/*
NewConsumer creates a new KafkaConsumer instance with the provided configuration.

It initializes:
  - A Kafka consumer connected to the specified topic.
  - A DLQ producer for handling failed messages.
  - A handler for processing messages.
  - A notifier for critical errors.

Returns the fully initialized KafkaConsumer or an error if setup fails.
*/
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
		dbConnectionCheckDelay:   config.DbConnectionCheckDelay,
		notifier:                 notifier.NewNotifier(config.Notifier)}, nil
}

/*
toMap converts a consumer or producer configuration into a Kafka ConfigMap.

This helper function maps internal configuration structs to
the format expected by the Confluent Kafka Go client.
*/
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

/*
Run starts the KafkaConsumer loop for a single worker.

Behavior:
  - Polls messages from Kafka continuously.
  - Processes each message with retries.
  - Commits offsets with retries.
  - Sends messages to DLQ if processing fails.
  - Logs errors and triggers notifier notifications for critical errors.
  - Pauses order processing during database outages with periodic connection checks.
  - Panics for unrecoverable errors, which may trigger worker self-termination.
*/
func (c *KafkaConsumer) Run(ctx context.Context, storage repository.Storage, logger logger.Logger, workerID int, lastWorker *atomic.Int32) {
	logger.LogInfo(fmt.Sprintf("worker %d — receiving orders", workerID), "layer", "broker.kafka")
	eventTypeErrors := 0
	for {
		select {
		case <-ctx.Done():
			if lastWorker.Load() == int32(1) { // I do realize how utterly retarded this is
				c.dlq.Close() // should've delegated DLQ management to the Consumer, not embedded it in each worker
				logger.LogInfo(fmt.Sprintf("worker %d — DLQ closed", workerID), "workerID", fmt.Sprintf("%d", workerID), "layer", "broker.kafka")
			} else {
				lastWorker.Add(-1)
			}
			return
		default:
			event := c.consumer.Poll(100)
			if event == nil {
				continue
			}
			switch eventType := event.(type) {
			case *kafka.Message:
				logger.Debug(fmt.Sprintf("worker %d — received a new order from Kafka, will try saving it", workerID), "workerID", fmt.Sprintf("%d", workerID), "layer", "broker.kafka")
				eventTypeErrors = 0
				var lastErr error
				var notified bool
				retryCnt := 0
				for retryCnt < c.saveOrderRetryMax {
					if err := c.handler.SaveOrder(eventType.Value, storage, logger, workerID); err != nil {
						if strings.Contains(err.Error(), "connection refused") {
							if !notified {
								logger.LogInfo(fmt.Sprintf("worker %d — lost connection to database, order processing paused", workerID), "layer", "broker.kafka")
								_ = c.notifier.Notify(fmt.Sprintf("CRITICAL ERROR — database connection lost, consumer worker %d paused", workerID))
								notified = true
							}
							time.Sleep(c.dbConnectionCheckDelay)
							continue
						}
						notified = false
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
						_ = c.notifier.Notify(fmt.Sprintf("CRITICAL ERROR — Kafka commit failed\nworkerID=%d\norderUID=%s", workerID, ToStr(eventType.Key)))
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
					_ = c.notifier.Notify(fmt.Sprintf("CRITICAL ERROR — Kafka broker is unreachable\nworkerID=%d", workerID))
					panic(fmt.Sprintf("worker self-termination: kafka is down (workerID=%d)", workerID))
				}
				time.Sleep(c.eventTypeErrorRetryDelay)
			}
		}
	}
}

/*
commitWithRetry attempts to commit a Kafka message offset multiple times.

It retries up to commitRetryMax times with a configured delay between attempts.
Returns an error if the commit fails after all retries.
*/
func (c *KafkaConsumer) commitWithRetry(msg *kafka.Message) error {
	var err error
	for range c.commitRetryMax {
		if _, err = c.consumer.CommitMessage(msg); err != nil {
			time.Sleep(c.commitRetryDelay)
		} else {
			return nil
		}
	}
	return fmt.Errorf("failed to commit offset after %d attempts: %w", c.commitRetryMax, err)
}

/*
sendToDLQ sends a failed message to the dead-letter queue (DLQ).

It attempts to produce the message to the DLQ and commit its offset.
If either action fails, it logs the error, notifies via the notifier,
and panics to trigger worker self-termination.

This self-termination ensures that the worker does not keep consuming CPU
in a tight loop when Kafka is down or offset commits repeatedly fail,
allowing the orchestration layer to handle restart or shutdown.
*/
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
		_ = c.notifier.Notify(fmt.Sprintf("CRITICAL ERROR — failed to send order to DLQ\nworkerID=%d\norderUID=%s", workerID, ToStr(eventType.Key)))
		panic(fmt.Sprintf("worker self-termination: failed to send order to DLQ (workerID=%d, orderUID=%s)", workerID, ToStr(eventType.Key)))
	}
	if err := c.commitWithRetry(eventType); err != nil {
		_ = c.notifier.Notify(fmt.Sprintf("CRITICAL ERROR — order sent to DLQ but offset commit failed\nworkerID=%d\norderUID=%s", workerID, ToStr(eventType.Key)))
		panic(fmt.Sprintf("worker self-termination: order sent to DLQ but offset commit failed (workerID=%d, orderUID=%s)", workerID, ToStr(eventType.Key)))
	}
}

/*
ToStr converts a Kafka message key from bytes to a trimmed string.

Kafka message keys may arrive with surrounding quotes. This function trims them
so that the key appears consistent in logs, matching the format of order UIDs
used elsewhere in the application.
*/
func ToStr(key []byte) string {
	return strings.Trim(string(key), `"`)
}

/*
Close terminates the Kafka consumer instance.

It releases all resources, stops receiving messages,
and logs any errors encountered during shutdown.
*/
func (c *KafkaConsumer) Close(logger logger.Logger) {
	if err := c.consumer.Close(); err != nil {
		logger.LogError("consumer — failed to stop properly", err, "layer", "broker.kafka")
	}
	logger.LogInfo("consumer — stopped receiving orders", "layer", "broker.kafka")
}
