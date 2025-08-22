package broker

import (
	"context"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/broker/kafka"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
)

type EventConsumer interface {
	Run(ctx context.Context, storage *repository.Storage, logger logger.Logger)
	Close(logger logger.Logger)
}

type Consumer struct {
	EventConsumer
}

func NewConsumer(config configs.Consumer) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	return &Consumer{EventConsumer: consumer}, nil
}
