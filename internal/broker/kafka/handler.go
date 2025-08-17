package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
)

type MessageHandler interface {
	SaveOrder(jsonMsg []byte, storage repository.Storage) error
}

type Handler struct{}

func NewHandler() *Handler {
	return new(Handler)
}

func (h *Handler) SaveOrder(jsonMsg []byte, storage repository.Storage) error {
	order := new(models.Order)
	if err := json.Unmarshal(jsonMsg, order); err != nil {
		return fmt.Errorf("consumer-handler — failed to unmarshal the order: %v", err)
	}
	if err := storage.Ping(); err != nil {
		return fmt.Errorf("consumer-handler — lost connection to database: %v", err)
	}
	if err := storage.SaveOrder(order); err != nil {
		return fmt.Errorf("consumer-handler — failed to save order %s to database: %v", order.OrderUID, err)
	}
	logger.LogInfo("consumer-handler — saved order", "orderUID", order.OrderUID)
	return nil
}
