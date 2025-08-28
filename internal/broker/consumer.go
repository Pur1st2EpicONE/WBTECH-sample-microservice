package broker

import (
	"context"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/broker/kafka"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
)

type Consumer interface {
	Run(ctx context.Context, storage repository.Storage, logger logger.Logger, workerID int)
	Close(logger logger.Logger)
}

func NewConsumer(config configs.Consumer, logger logger.Logger) (Consumer, error) {
	consumer, err := kafka.NewConsumer(config, logger)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}
