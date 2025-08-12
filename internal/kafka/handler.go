package kafka

import (
	"encoding/json"
	"fmt"

	model "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/repository"
)

type OrderHandler struct {
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) SaveOrder(jsonMsg []byte, db repository.Storage) error {
	order := new(model.Order)
	if err := json.Unmarshal(jsonMsg, order); err != nil {
		return fmt.Errorf("failed to unmarshal order: %w", err)
	}
	db.SaveOrder(order)

	return nil
}
