package broker

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/broker/kafka"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
)

type EventProducer interface {
	Produce(data []byte, topic string) error
	Close()
}

type Producer struct {
	EventProducer
}

func NewProducer(config configs.Producer) (*Producer, error) {
	kafka, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}
	return &Producer{EventProducer: kafka}, nil
}
