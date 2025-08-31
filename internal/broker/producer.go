package broker

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/broker/kafka"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
)

/*
Producer defines a message broker producer instance.

A Producer is responsible for:
  - Sending messages to a broker (e.g., to a dead-letter queue).
  - Closing the producer and releasing any resources.
*/
type Producer interface {
	// Produce sends a message to the broker.
	Produce(message configs.Message) error

	// Close terminates the producer and cleans up resources.
	Close()
}

/*
NewProducer creates a new Producer instance based on the provided configuration.

Currently, it constructs a Kafka-based producer, but the function
can be extended to support other broker types in the future.
Returns the fully initialized Producer or an error if setup fails.
*/
func NewProducer(config configs.Producer, logger logger.Logger) (Producer, error) {
	producer, err := kafka.NewProducer(config, logger)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
