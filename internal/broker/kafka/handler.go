package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
	"github.com/go-playground/validator/v10"
)

type MessageHandler interface {
	SaveOrder(jsonMsg []byte, storage repository.Storage, logger logger.Logger, consumerID int) error
}

type Handler struct{}

func newHandler() *Handler {
	return new(Handler)
}

func (h *Handler) SaveOrder(jsonMsg []byte, storage repository.Storage, logger logger.Logger, workerID int) error {
	validate := validator.New()
	order := new(models.Order)
	if err := json.Unmarshal(jsonMsg, order); err != nil {
		return fmt.Errorf("failed to unmarshal the order: %v", err)
	}
	if err := validate.Struct(order); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}
	if err := storage.Ping(); err != nil {
		return fmt.Errorf("lost connection to database: %v", err)
	}
	if err := storage.SaveOrder(order); err != nil {
		return fmt.Errorf("failed to save order %s to database: %v", order.OrderUID, err)
	}
	logger.Debug(fmt.Sprintf("worker %d — saved order to DB", workerID), "orderUID", order.OrderUID, "workerID", fmt.Sprintf("%d", workerID), "layer", "broker.kafka")
	return nil
}
