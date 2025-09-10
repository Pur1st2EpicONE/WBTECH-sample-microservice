/*
Package broker provides abstractions for message broker interactions.

It defines interfaces and factory functions for consumers and producers.
Consumers are responsible for receiving and processing messages from a broker,
while producers send messages (for example, to a dead-letter queue).
Both are managed and supervised by the application's orchestration layer.
*/
package broker

import (
	"context"
	"sync/atomic"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/broker/kafka"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
)

/*
Consumer defines a message broker consumer instance.

A Consumer is responsible for:
  - Running a worker loop to consume messages.
  - Gracefully shutting down when requested.
  - Being supervised by the orchestration layer (App) for panics or errors.
*/
type Consumer interface {
	// Run starts the consumer loop for a single worker.
	// It processes messages until the context is cancelled.
	Run(ctx context.Context, storage repository.Storage, logger logger.Logger, workerID int, lastWorker *atomic.Int32)

	// Close terminates the consumer and releases any underlying resources.
	Close(logger logger.Logger)
}

/*
NewConsumer creates a new Consumer instance based on the provided configuration.

Currently, it constructs a Kafka-based consumer, but the function
can be extended to support other broker types in the future.
Returns the fully initialized Consumer or an error if setup fails.
*/
func NewConsumer(config configs.Consumer, logger logger.Logger) (Consumer, error) {
	consumer, err := kafka.NewConsumer(config, logger)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}
