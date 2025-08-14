package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/repository"
)

type OrderHandler struct {
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) SaveOrder(jsonMsg []byte, db repository.Storage) error {
	order := new(models.Order)
	if err := json.Unmarshal(jsonMsg, order); err != nil {
		return fmt.Errorf("failed to unmarshal the order: %v", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("lost connection to database: %v", err)
	}
	if err := db.SaveOrder(order); err != nil {
		return fmt.Errorf("failed to save the order: %v", err)
	}
	return nil
}
