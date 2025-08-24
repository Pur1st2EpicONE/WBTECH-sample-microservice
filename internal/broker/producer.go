package broker

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/broker/kafka"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
)

type Producer interface {
	Produce(message configs.Message) error
	Close()
}

func NewProducer(config configs.Producer, logger logger.Logger) (Producer, error) {
	producer, err := kafka.NewProducer(config, logger)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
